.PHONY: build test lint clean help

# Build the kubertino binary
build:
	@echo "Building kubertino..."
	go build -o kubertino cmd/kubertino/main.go
	@echo "Build complete: ./kubertino"

# Run all tests with coverage
test:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Coverage report:"
	go tool cover -func=coverage.out

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f kubertino
	rm -f coverage.out
	@echo "Clean complete"

# Display available targets
help:
	@echo "Available targets:"
	@echo "  build  - Compile the kubertino binary"
	@echo "  test   - Run all tests with coverage report"
	@echo "  lint   - Run golangci-lint"
	@echo "  clean  - Remove build artifacts (binary, coverage files)"
	@echo "  help   - Display this help message"