.PHONY: test build clean

# Default target
all: test build

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Build the binary
build:
	@echo "Building..."
	@go build -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f openraster-go
	@go clean

# Run tests with coverage
cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Help target
help:
	@echo "Available targets:"
	@echo "  test   - Run tests"
	@echo "  build  - Build the binary"
	@echo "  clean  - Clean build artifacts"
	@echo "  cover  - Run tests with coverage"
	@echo "  fmt    - Format code"
	@echo "  lint   - Run linter"
	@echo "  help   - Show this help message"
