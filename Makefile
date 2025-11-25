.PHONY: help build test clean install run lint fmt vendor

# Variables
BINARY_NAME=bug-crawler
MAIN_PACKAGE=./cmd/main.go
BINARY_PATH=./$(BINARY_NAME)
GO=go
GOFLAGS=-v

# Default target
.DEFAULT_GOAL := help

# Help target
help:
	@echo "🐛 Bug Crawler - Makefile Targets"
	@echo "=================================="
	@echo ""
	@echo "Build Targets:"
	@echo "  make build           - Build the bug-crawler binary"
	@echo "  make build-debug     - Build with debug flags"
	@echo "  make install         - Install the binary to GOPATH/bin"
	@echo ""
	@echo "Test Targets:"
	@echo "  make test            - Run all tests"
	@echo "  make test-verbose    - Run tests with verbose output"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make test-short      - Run tests in short mode"
	@echo ""
	@echo "Utility Targets:"
	@echo "  make clean           - Remove built binaries and temp files"
	@echo "  make fmt             - Format Go code"
	@echo "  make lint            - Run linter (requires golangci-lint)"
	@echo "  make vendor          - Download and vendor dependencies"
	@echo "  make run             - Run the application"
	@echo "  make help            - Display this help message"
	@echo ""

# ============================================================================
# BUILD TARGETS
# ============================================================================

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@$(GO) mod download
	@$(GO) mod tidy
	@$(GO) build -o $(BINARY_PATH) $(MAIN_PACKAGE)
	@echo "✅ Build successful!"
	@echo "📦 Binary: $(BINARY_PATH)"

build-debug:
	@echo "🔨 Building $(BINARY_NAME) with debug flags..."
	@$(GO) build -v -x -o $(BINARY_PATH) $(MAIN_PACKAGE)
	@echo "✅ Debug build successful!"

install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@$(GO) install $(MAIN_PACKAGE)
	@echo "✅ Installation successful!"

# ============================================================================
# TEST TARGETS
# ============================================================================

test:
	@echo "🧪 Running tests..."
	@$(GO) test -v ./...
	@echo "✅ Tests completed!"

test-verbose:
	@echo "🧪 Running tests (verbose)..."
	@$(GO) test -v -race ./...
	@echo "✅ Tests completed!"

test-coverage:
	@echo "📊 Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-short:
	@echo "🧪 Running short tests..."
	@$(GO) test -short ./...
	@echo "✅ Short tests completed!"

# ============================================================================
# UTILITY TARGETS
# ============================================================================

clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(BINARY_PATH)
	@$(GO) clean
	@rm -f coverage.out coverage.html
	@echo "✅ Cleanup completed!"

fmt:
	@echo "📝 Formatting code..."
	@$(GO) fmt ./...
	@echo "✅ Code formatted!"

lint:
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null || (echo "❌ golangci-lint not installed. Install with: brew install golangci-lint"; exit 1)
	@golangci-lint run ./...
	@echo "✅ Linting completed!"

vendor:
	@echo "📥 Downloading and vendoring dependencies..."
	@$(GO) mod download
	@$(GO) mod vendor
	@$(GO) mod tidy
	@echo "✅ Dependencies vendored!"

run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@$(BINARY_PATH)

# ============================================================================
# COMBINED TARGETS
# ============================================================================

all: clean fmt build test
	@echo "✅ All tasks completed successfully!"

check: fmt lint test
	@echo "✅ All checks passed!"
