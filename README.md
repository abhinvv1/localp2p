# LocalP2P - Phase 1 Implementation

A secure peer-to-peer communication system for local networks without internet dependency.

## Quick Start

### 1. Build the Project
```bash
./scripts/build.sh
```

### 2. Test Locally
```bash
# Terminal 1: Start first peer
cd core && ./localp2p --rpc-port=9090

# Terminal 2: Start second peer  
cd core && ./localp2p --rpc-port=9091

# Terminal 3: Test CLI
cd cli
npm start discover
npm start connect
npm start send "Hello World!"
```

### 3. Docker Testing (Recommended)
```bash
./scripts/test-docker.sh
```

## Manual Docker Testing

### Start Services
```bash
cd docker
docker-compose up -d --build
```

### Test Commands
```bash
# Discover peers
docker exec -it localp2p-peer1 ./cli/bin/localp2p discover
docker exec -it localp2p-peer2 ./cli/bin/localp2p discover

# Connect peers
docker exec -it localp2p-peer1 ./cli/bin/localp2p connect --address 172.20.0.11 --port 8080

# Check connections
docker exec -it localp2p-peer1 ./cli/bin/localp2p status

# Send messages
docker exec -it localp2p-peer1 ./cli/bin/localp2p send "Hello from peer1!"
```

### Cleanup
```bash
docker-compose down
docker-compose down -v  # Also remove volumes
```

## Project Structure

- `core/` - Go-based P2P engine with mDNS discovery and TCP transport
- `cli/` - Node.js CLI interface for user interaction
- `docker/` - Docker setup for testing
- `scripts/` - Build and test automation
- `test/` - Integration tests

## Phase 1 Features

✅ **mDNS-based peer discovery** - Zero-configuration peer finding
✅ **TCP direct connections** - Reliable peer-to-peer communication  
✅ **Basic authentication** - Challenge-response authentication
✅ **CLI interface** - Easy-to-use command line tools
✅ **Cross-platform support** - Windows, Linux, macOS
✅ **Docker testing** - Containerized testing environment

## Testing Strategy

1. **Local Testing** - Run multiple instances on different ports
2. **Docker Testing** - Isolated containers with custom network
3. **Integration Testing** - Automated test suite
4. **Manual Testing** - Interactive CLI testing

## Troubleshooting

### Core Won't Start
- Check if ports are available: `netstat -an | grep 8080`
- Verify Go binary is built: `cd core && go build -o localp2p`

### Discovery Not Working
- Ensure both peers are on same network
- Check firewall settings for mDNS (port 5353)
- Verify mDNS is enabled on system

### Connection Failures
- Check if target peer is reachable: `ping <peer-ip>`
- Verify ports are open: `telnet <peer-ip> 8080`
- Check logs for authentication errors

## Next Steps (Phase 2)

- Implement Signal Protocol encryption
- Add UDP hole punching for NAT traversal
- Enhanced security with proper key exchange
- Connection health monitoring and failover

## Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Node.js CLI   │    │   Node.js CLI   │
├─────────────────┤    ├─────────────────┤
│   Go P2P Core   │◄──►│   Go P2P Core   │
│                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │ mDNS Disc.  │ │    │ │ mDNS Disc.  │ │
│ ├─────────────┤ │    │ ├─────────────┤ │
│ │ TCP Transport│◄┼────┼►│ TCP Transport│ │
│ ├─────────────┤ │    │ ├─────────────┤ │
│ │ Basic Auth  │ │    │ │ Basic Auth  │ │
│ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘
        Peer 1                 Peer 2
```