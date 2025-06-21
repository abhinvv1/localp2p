package transport

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "sync"
    "time"

    "localp2p/security"
)

type Message struct {
    Type      string    `json:"type"`
    From      string    `json:"from"`
    To        string    `json:"to"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}

type AuthMessage struct {
    Type      string `json:"type"`
    NodeID    string `json:"node_id"`
    Challenge string `json:"challenge,omitempty"`
    Response  string `json:"response,omitempty"`
}

type Connection struct {
    conn        net.Conn
    peer        string
    authenticated bool
    lastSeen    time.Time
}

type Transport struct {
    nodeID       string
    port         int
    listener     net.Listener
    connections  map[string]*Connection
    connMutex    sync.RWMutex
    auth         *security.Authenticator
    messageQueue chan Message
}

func NewTransport(nodeID string, port int) *Transport {
    return &Transport{
        nodeID:       nodeID,
        port:         port,
        connections:  make(map[string]*Connection),
        auth:         security.NewAuthenticator(nodeID),
        messageQueue: make(chan Message, 100),
    }
}

func (t *Transport) Start() error {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", t.port))
    if err != nil {
        return err
    }
    
    t.listener = listener
    log.Printf("Transport listening on port %d", t.port)
    
    // Start accepting connections
    go t.acceptConnections()
    
    // Start connection cleanup
    go t.cleanupConnections()
    
    return nil
}

func (t *Transport) Stop() {
    if t.listener != nil {
        t.listener.Close()
    }
    
    t.connMutex.Lock()
    for _, conn := range t.connections {
        conn.conn.Close()
    }
    t.connMutex.Unlock()
}

func (t *Transport) acceptConnections() {
    for {
        conn, err := t.listener.Accept()
        if err != nil {
            log.Printf("Accept error: %v", err)
            continue
        }
        
        go t.handleConnection(conn)
    }
}

func (t *Transport) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // Perform authentication
    peerID, err := t.authenticate(conn)
    if err != nil {
        log.Printf("Authentication failed: %v", err)
        return
    }
    
    // Store connection
    t.connMutex.Lock()
    t.connections[peerID] = &Connection{
        conn:          conn,
        peer:          peerID,
        authenticated: true,
        lastSeen:      time.Now(),
    }
    t.connMutex.Unlock()
    
    log.Printf("Authenticated connection from %s", peerID)
    
    // Handle messages
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        var msg Message
        if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
            log.Printf("Invalid message from %s: %v", peerID, err)
            continue
        }
        
        msg.From = peerID
        t.messageQueue <- msg
        
        // Update last seen
        t.connMutex.Lock()
        if connection, exists := t.connections[peerID]; exists {
            connection.lastSeen = time.Now()
        }
        t.connMutex.Unlock()
    }
    
    // Remove connection
    t.connMutex.Lock()
    delete(t.connections, peerID)
    t.connMutex.Unlock()
    
    log.Printf("Connection closed with %s", peerID)
}

func (t *Transport) authenticate(conn net.Conn) (string, error) {
    // Set timeout for authentication
    conn.SetDeadline(time.Now().Add(10 * time.Second))
    defer conn.SetDeadline(time.Time{})
    
    // Send challenge
    challenge := t.auth.GenerateChallenge()
    authMsg := AuthMessage{
        Type:      "challenge",
        NodeID:    t.nodeID,
        Challenge: challenge,
    }
    
    if err := t.sendAuthMessage(conn, authMsg); err != nil {
        return "", err
    }
    
    // Receive response
    var response AuthMessage
    if err := t.receiveAuthMessage(conn, &response); err != nil {
        return "", err
    }
    
    // Verify response
    if !t.auth.VerifyResponse(challenge, response.Response) {
        return "", fmt.Errorf("authentication failed")
    }
    
    return response.NodeID, nil
}

func (t *Transport) sendAuthMessage(conn net.Conn, msg AuthMessage) error {
    data, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    
    _, err = conn.Write(append(data, '\n'))
    return err
}

func (t *Transport) receiveAuthMessage(conn net.Conn, msg *AuthMessage) error {
    scanner := bufio.NewScanner(conn)
    if !scanner.Scan() {
        return fmt.Errorf("failed to read message")
    }
    
    return json.Unmarshal(scanner.Bytes(), msg)
}

func (t *Transport) ConnectToPeer(address string, port int) error {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
    if err != nil {
        return err
    }
    
    // Perform client-side authentication
    peerID, err := t.authenticateClient(conn)
    if err != nil {
        conn.Close()
        return err
    }
    
    // Store connection
    t.connMutex.Lock()
    t.connections[peerID] = &Connection{
        conn:          conn,
        peer:          peerID,
        authenticated: true,
        lastSeen:      time.Now(),
    }
    t.connMutex.Unlock()
    
    log.Printf("Connected to peer %s", peerID)
    
    // Handle incoming messages
    go t.handleConnection(conn)
    
    return nil
}

func (t *Transport) authenticateClient(conn net.Conn) (string, error) {
    // Set timeout for authentication
    conn.SetDeadline(time.Now().Add(10 * time.Second))
    defer conn.SetDeadline(time.Time{})
    
    // Receive challenge
    var challenge AuthMessage
    if err := t.receiveAuthMessage(conn, &challenge); err != nil {
        return "", err
    }
    
    // Compute and send response
    response := t.auth.ComputeResponse(challenge.Challenge)
    authMsg := AuthMessage{
        Type:     "response",
        NodeID:   t.nodeID,
        Response: response,
    }
    
    if err := t.sendAuthMessage(conn, authMsg); err != nil {
        return "", err
    }
    
    return challenge.NodeID, nil
}

func (t *Transport) SendMessage(to, content string) error {
    t.connMutex.RLock()
    connection, exists := t.connections[to]
    t.connMutex.RUnlock()
    
    if !exists {
        return fmt.Errorf("no connection to peer %s", to)
    }
    
    msg := Message{
        Type:      "message",
        From:      t.nodeID,
        To:        to,
        Content:   content,
        Timestamp: time.Now(),
    }
    
    data, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    
    _, err = connection.conn.Write(append(data, '\n'))
    return err
}

func (t *Transport) GetMessages()  60*time.Second {
                conn.conn.Close()
                delete(t.connections, peer)
                log.Printf("Cleaned up stale connection to %s", peer)
            }
        }
        t.connMutex.Unlock()
    }
}