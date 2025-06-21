package discovery

import (
    "context"
    "fmt"
    "log"
    "net"
    "strconv"
    "sync"
    "time"

    "github.com/grandcat/zeroconf"
    "localp2p/config"
)

type Peer struct {
    ID       string
    Name     string
    Address  string
    Port     int
    LastSeen time.Time
}

type Discovery struct {
    config    *config.Config
    peers     map[string]*Peer
    peerMutex sync.RWMutex
    server    *zeroconf.Server
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewDiscovery(cfg *config.Config) *Discovery {
    ctx, cancel := context.WithCancel(context.Background())
    return &Discovery{
        config: cfg,
        peers:  make(map[string]*Peer),
        ctx:    ctx,
        cancel: cancel,
    }
}

func (d *Discovery) Start() error {
    // Start mDNS server to advertise this node
    if err := d.startServer(); err != nil {
        return fmt.Errorf("failed to start mDNS server: %v", err)
    }
    
    // Start discovering other peers
    go d.startDiscovery()
    
    log.Printf("Discovery started for %s on port %d", d.config.NodeID, d.config.Port)
    return nil
}

func (d *Discovery) Stop() {
    if d.server != nil {
        d.server.Shutdown()
    }
    d.cancel()
}

func (d *Discovery) startServer() error {
    // Get local IP address
    ip, err := getLocalIP()
    if err != nil {
        return err
    }
    
    // Create mDNS server
    server, err := zeroconf.Register(
        d.config.NodeID,                    // Instance name
        d.config.ServiceName,               // Service type
        d.config.Domain,                    // Domain
        d.config.Port,                      // Port
        []string{"id=" + d.config.NodeID, "name=" + d.config.DisplayName}, // TXT records
        []net.IP{ip},                       // IP addresses
    )
    
    if err != nil {
        return err
    }
    
    d.server = server
    return nil
}

func (d *Discovery) startDiscovery() {
    resolver, err := zeroconf.NewResolver(nil)
    if err != nil {
        log.Printf("Failed to create resolver: %v", err)
        return
    }
    
    entries := make(chan *zeroconf.ServiceEntry)
    
    go func() {
        for entry := range entries {
            d.handleServiceEntry(entry)
        }
    }()
    
    // Continuous discovery
    for {
        select {
        case <-d.ctx.Done():
            return
        default:
            ctx, cancel := context.WithTimeout(d.ctx, 5*time.Second)
            err := resolver.Browse(ctx, d.config.ServiceName, d.config.Domain, entries)
            cancel()
            
            if err != nil {
                log.Printf("Discovery error: %v", err)
            }
            
            time.Sleep(10 * time.Second) // Discovery interval
        }
    }
}

func (d *Discovery) handleServiceEntry(entry *zeroconf.ServiceEntry) {
    // Skip our own service
    if entry.Instance == d.config.NodeID {
        return
    }
    
    d.peerMutex.Lock()
    defer d.peerMutex.Unlock()
    
    // Extract peer information
    peer := &Peer{
        ID:       entry.Instance,
        Name:     entry.Instance,
        Address:  entry.AddrIPv4[0].String(),
        Port:     entry.Port,
        LastSeen: time.Now(),
    }
    
    // Parse TXT records for additional info
    for _, txt := range entry.Text {
        if len(txt) > 5 && txt[:5] == "name=" {
            peer.Name = txt[5:]
        }
    }
    
    d.peers[peer.ID] = peer
    log.Printf("Discovered peer: %s (%s:%d)", peer.Name, peer.Address, peer.Port)
}

func (d *Discovery) GetPeers() []*Peer {
    d.peerMutex.RLock()
    defer d.peerMutex.RUnlock()
    
    var peers []*Peer
    now := time.Now()
    
    for _, peer := range d.peers {
        // Only return peers seen in the last 60 seconds
        if now.Sub(peer.LastSeen) < 60*time.Second {
            peers = append(peers, peer)
        }
    }
    
    return peers
}

func getLocalIP() (net.IP, error) {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return nil, err
    }
    defer conn.Close()
    
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP, nil
}