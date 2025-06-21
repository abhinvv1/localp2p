package security

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

type Authenticator struct {
    nodeID string
    secret []byte
}

func NewAuthenticator(nodeID string) *Authenticator {
    // Generate a simple secret for this session
    secret := make([]byte, 32)
    rand.Read(secret)
    
    return &Authenticator{
        nodeID: nodeID,
        secret: secret,
    }
}

func (a *Authenticator) GenerateChallenge() string {
    // Simple challenge-response for Phase 1
    challenge := make([]byte, 16)
    rand.Read(challenge)
    return hex.EncodeToString(challenge)
}

func (a *Authenticator) VerifyResponse(challenge, response string) bool {
    // Simple verification - in production, use proper crypto
    expected := a.computeResponse(challenge)
    return response == expected
}

func (a *Authenticator) ComputeResponse(challenge string) string {
    return a.computeResponse(challenge)
}

func (a *Authenticator) computeResponse(challenge string) string {
    h := sha256.New()
    h.Write([]byte(a.nodeID))
    h.Write([]byte(challenge))
    h.Write(a.secret)
    return hex.EncodeToString(h.Sum(nil))
}

func (a *Authenticator) GetNodeID() string {
    return a.nodeID
}