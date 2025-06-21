# Check all components
echo "🔍 Verifying complete setup..."
echo ""

echo "📋 System Info:"
echo "OS: $(uname -a)"
echo "Go: $(go version)"
echo "Node: $(node --version)"
echo "NPM: $(npm --version)"
echo "Docker: $(docker --version 2>/dev/null || echo 'Not available')"
echo ""

echo "🔧 Go Environment:"
echo "GOPATH: $GOPATH"
echo "GOROOT: $(go env GOROOT)"
echo "GOBIN: $GOBIN"
echo ""

echo "📁 Project Structure:"
ls -la ~/localp2p/
echo ""

echo "📦 Go Dependencies (core):"
cd ~/localp2p/core && go list -m all && cd ~
echo ""

echo "�� Node Dependencies (cli):"
cd ~/localp2p/cli && npm list --depth=0 && cd ~
echo ""

echo "✅ Environment setup complete!"
echo ""
echo "🚀 Next steps:"
echo "  1. Copy the LocalP2P code files from the previous artifact"
echo "  2. Run: cd ~/localp2p && ./scripts/build.sh"
echo "  3. Test with: ./scripts/test-docker.sh"
