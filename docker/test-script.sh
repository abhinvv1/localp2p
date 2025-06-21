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