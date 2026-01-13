# Makefile for Leaf project

.PHONY: help test test-verbose test-watch test-coverage test-app test-storage clean build run

# Default target
help:
	@echo "🌱 Leaf - Available commands:"
	@echo ""
	@echo "  make test          - Run all tests with gotestsum"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-watch    - Run tests in watch mode"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-app      - Run only app tests"
	@echo "  make test-storage  - Run only storage tests"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make clean         - Clean build artifacts"
	@echo ""

# Run all tests with gotestsum
test:
	@echo "🧪 Running tests..."
	@gotestsum --format pkgname-and-test-fails -- -cover ./...

# Run tests with verbose output
test-verbose:
	@echo "🧪 Running tests (verbose)..."
	@gotestsum --format standard-verbose -- -cover -v ./...

# Run tests in watch mode
test-watch:
	@echo "👀 Running tests in watch mode..."
	@gotestsum --watch --format testname

# Run tests with coverage report
test-coverage:
	@echo "📊 Running tests with coverage..."
	@gotestsum --format pkgname-and-test-fails -- -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Run only app tests
test-app:
	@echo "🧪 Running app tests..."
	@gotestsum --format testname -- -v ./tests/app/...

# Run only storage tests
test-storage:
	@echo "🧪 Running storage tests..."
	@gotestsum --format testname -- -v ./tests/storage/...

# Build the application
build:
	@echo "🔨 Building leaf..."
	@go build -o bin/leaf ./cmd/leaf
	@echo "✓ Built: bin/leaf"

# Run the application
run:
	@echo "🚀 Running leaf..."
	@go run ./cmd/leaf

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/ coverage.out coverage.html
	@echo "✓ Cleaned"

# Install gotestsum if not present
install-tools:
	@echo "📦 Installing development tools..."
	@go install gotest.tools/gotestsum@latest
	@echo "✓ Tools installed"
