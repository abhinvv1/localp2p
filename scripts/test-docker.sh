#!/bin/bash

set -e

echo "🐳 Running Docker integration tests..."

# Navigate to docker directory
cd docker

# Run the test script
./test-script.sh

echo "✅ Docker tests completed!"
