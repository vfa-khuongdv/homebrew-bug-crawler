#!/bin/bash

# setup.sh - Helper script để setup bug-crawler

set -e

echo "🐛 Setting up Bug Crawler..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Build
echo "📥 Building application..."
./build.sh

# Create config directory
CONFIG_DIR="$HOME/.config/bug-crawler"
if [ ! -d "$CONFIG_DIR" ]; then
    mkdir -p "$CONFIG_DIR"
    echo "✓ Created config directory: $CONFIG_DIR"
fi

# Check if token file exists
if [ -f "$CONFIG_DIR/token" ]; then
    echo "✓ Token file found"
else
    echo "ℹ️  No token file found. You'll be prompted to enter token on first run."
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Get a GitHub token from: https://github.com/settings/tokens"
echo "2. Run: ./bug-crawler"
echo "3. Follow the interactive prompts"
echo ""
echo "For more information, see README.md and USAGE.md"
