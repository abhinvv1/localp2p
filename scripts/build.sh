#!/bin/bash

set -e

echo "🏗️  Building LocalP2P..."

# Build Go core
echo "📦 Building Go core..."
cd core
go mod tidy
go build -o localp2p
cd ..

# Install CLI dependencies
echo "📦 Installing CLI dependencies..."
cd cli
npm install
cd ..

# Make scripts executable
chmod +x scripts/*.sh
chmod +x cli/bin/localp2p
chmod +x docker/test-script.sh

echo "✅ Build completed!"
echo ""
echo "🚀 Quick start:"
echo "  1. ./scripts/test-local.sh  # Test locally"
echo "  2. ./scripts/test-docker.sh # Test with Docker"