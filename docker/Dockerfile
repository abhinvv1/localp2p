FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy Go modules
COPY core/go.mod core/go.sum ./
RUN go mod download

# Copy source code
COPY core/ ./

# Build the application
RUN go build -o localp2p

FROM node:18-alpine

# Install system dependencies
RUN apk add --no-cache curl

# Set working directory
WORKDIR /app

# Copy built Go binary
COPY --from=builder /app/localp2p ./core/localp2p

# Copy Node.js CLI
COPY cli/ ./cli/

# Install CLI dependencies
WORKDIR /app/cli
RUN npm install

# Make CLI executable
RUN chmod +x bin/localp2p

# Expose ports
EXPOSE 8080 9090

# Set working directory back to app root
WORKDIR /app

# Default command
CMD ["./cli/bin/localp2p", "start"]
```

### docker/docker-compose.yml
```yaml
version: '3.8'

services:
  peer1:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    container_name: localp2p-peer1
    hostname: peer1
    ports:
      - "8081:8080"
      - "9091:9090"
    environment:
      - NODE_ENV=development
    networks:
      localp2p-net:
        ipv4_address: 172.20.0.10
    volumes:
      - peer1-data:/app/data

  peer2:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    container_name: localp2p-peer2
    hostname: peer2
    ports:
      - "8082:8080"
      - "9092:9090"
    environment:
      - NODE_ENV=development
    networks:
      localp2p-net:
        ipv4_address: 172.20.0.11
    volumes:
      - peer2-data:/app/data

networks:
  localp2p-net:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16

volumes:
  peer1-data:
  peer2-data:
```

### docker/test-script.sh
```bash
#!/bin/bash

set -e

echo "🚀 Starting LocalP2P Docker Test..."

# Build and start containers
echo "📦 Building and starting containers..."
docker-compose up -d --build

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 10

# Function to run CLI command in container
run_cli() {
    local container=$1
    shift
    docker exec -it "localp2p-${container}" ./cli/bin/localp2p "$@"
}

# Function to run CLI command and capture output
run_cli_output() {
    local container=$1
    shift
    docker exec "localp2p-${container}" ./cli/bin/localp2p "$@"
}

echo "🔍 Testing peer discovery..."
echo "Peer1 discovering peers:"
run_cli peer1 discover

echo ""
echo "Peer2 discovering peers:"
run_cli peer2 discover

echo ""
echo "🔗 Testing connection..."
echo "Connecting peer1 to peer2..."

# Get peer2's IP and connect
PEER2_IP="172.20.0.11"
run_cli peer1 connect --address "$PEER2_IP" --port 8080

echo ""
echo "📊 Checking connection status..."
echo "Peer1 connections:"
run_cli peer1 status

echo ""
echo "Peer2 connections:"
run_cli peer2 status

echo ""
echo "💬 Testing message sending..."
echo "Sending message from peer1 to peer2..."

# Send test message
run_cli peer1 send "Hello from peer1!" --to "peer2"

echo ""
echo "✅ Test completed!"
echo ""
echo "🛠️  Manual testing commands:"
echo "docker exec -it localp2p-peer1 ./cli/bin/localp2p discover"
echo "docker exec -it localp2p-peer1 ./cli/bin/localp2p connect --address 172.20.0.11 --port 8080"
echo "docker exec -it localp2p-peer1 ./cli/bin/localp2p send 'Hello World!'"
echo ""
echo "To stop: docker-compose down"