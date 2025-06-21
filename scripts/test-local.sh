#!/bin/bash

set -e

echo "🧪 Running local tests..."

# Build first
./scripts/build.sh

# Start core in background
echo "🚀 Starting LocalP2P core..."
cd core
./localp2p --rpc-port=9090 &
CORE_PID=$!
cd ..

# Wait for core to start
sleep 3

# Test CLI commands
echo "🔍 Testing discovery..."
cd cli
npm start discover
echo ""

echo "📊 Testing status..."
npm start status
echo ""

# Cleanup
echo "🧹 Cleaning up..."
kill $CORE_PID 2>/dev/null || true

echo "✅ Local tests completed!"
