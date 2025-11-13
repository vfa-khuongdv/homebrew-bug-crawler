#!/bin/bash

# Build script for bug-crawler

set -e

echo "🐛 Building Bug Crawler..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi

# Download dependencies
echo "📥 Downloading dependencies..."
go mod download
go mod tidy

# Build the binary
echo "🔨 Building binary..."
go build -o bug-crawler ./cmd/main.go

echo "✅ Build successful!"
echo "📦 Binary location: ./bug-crawler"
echo "🚀 Run with: ./bug-crawler"
