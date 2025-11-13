#!/bin/bash

# Install script for bug-crawler

set -e

echo "🐛 Installing Bug Crawler..."

# Build the binary
./build.sh

# Create installation directory
INSTALL_DIR="${HOME}/.local/bin"
if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
fi

# Copy binary
cp bug-crawler "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/bug-crawler"

echo "✅ Installation successful!"
echo "📍 Binary installed at: $INSTALL_DIR/bug-crawler"
echo ""
echo "To use the command globally, add to your PATH:"
echo "export PATH=\"\$PATH:$INSTALL_DIR\""
echo ""
echo "Then you can run: bug-crawler"
