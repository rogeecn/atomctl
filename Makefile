# Makefile for pkg/ast/provider development

.PHONY: help test benchmark lint clean setup install-tools

# Default target
help:
	@echo "Available targets:"
	@echo "  setup         - Install required tools"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  benchmark     - Run benchmark tests"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  install-tools - Install development tools"

# Install required tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/AlekSi/gocov-xml@latest

# Setup development environment
setup: install-tools
	@echo "Setting up development environment..."
	go mod download
	go mod tidy

# Run all tests
test:
	@echo "Running tests..."
	go test -v -race -parallel=4 ./tests/... ./pkg/ast/provider/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./tests/... ./pkg/ast/provider/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmark tests
benchmark:
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem -count=3 ./pkg/ast/provider/...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./pkg/ast/provider/... ./tests/...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f coverage.out coverage.html
	rm -rf vendor/
	go clean -cache -testcache

# Run performance validation
validate-performance:
	@echo "Validating performance requirements..."
	@echo "Running single file parsing benchmark..."
	@go test -bench=ParseFile -benchtime=100ms -count=5 ./pkg/ast/provider/... | grep -E "(ParseFile|MB/s|allocs/op)"
	@echo "Running project parsing benchmark..."
	@go test -bench=ParseProject -benchtime=1s -count=3 ./pkg/ast/provider/... | grep -E "(ParseProject|MB/s|allocs/op)"

# Check test coverage meets requirements
check-coverage: test-coverage
	@echo "Checking coverage meets 90% requirement..."
	@go tool cover -func=coverage.out | grep total | awk '{if ($$3 < 90.0) {print "Coverage " $$3 " is below 90% requirement"; exit 1} else {print "Coverage " $$3 " meets requirement"}}'