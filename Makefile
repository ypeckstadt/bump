.PHONY: build clean test install dev-deps lint fmt vet security security-install security-json security-sarif security-ci security-clean setup-dev

# Build variables
BINARY_NAME=bump
VERSION=1.0.0
BUILD_DIR=bin
LDFLAGS=-ldflags "-X bump/pkg/version.Version=$(VERSION) \
                  -X bump/pkg/version.GitCommit=$(shell git rev-parse --short HEAD) \
                  -X bump/pkg/version.BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bump

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/bump
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/bump
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/bump
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/bump
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/bump

# Install binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./cmd/bump

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f bump-security.json results.sarif coverage.out coverage.html
	@go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Install development dependencies
dev-deps:
	@echo "Installing development dependencies..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

# Setup development environment
setup-dev: dev-deps
	@echo "Setting up development environment..."
	@mkdir -p .git/hooks
	@echo "Development environment setup complete."

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Vet code
vet:
	@echo "Running go vet..."
	go vet ./...

# Install gosec if not already installed
security-install:
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)

# Run security scan (default text output)
security: security-install
	@echo "Running security scan..."
	gosec ./...

# Run security scan with JSON output
security-json: security-install
	@echo "Running security scan (JSON output)..."
	gosec -fmt json -out bump-security.json ./...
	@echo "Security report saved to bump-security.json"

# Run security scan with SARIF output (same as CI)
security-sarif: security-install
	@echo "Running security scan (SARIF output)..."
	gosec -fmt sarif -out results.sarif ./...
	@echo "Security report saved to results.sarif"

# Run security scan exactly like CI pipeline
security-ci: security-sarif

# Clean security reports
security-clean:
	@echo "Cleaning security reports..."
	@rm -f bump-security.json results.sarif

# Run all checks including security
check: fmt vet lint test security

# Development build (with race detector)
dev:
	@echo "Building development version..."
	@mkdir -p $(BUILD_DIR)
	go build -race $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-dev ./cmd/bump

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  clean         - Clean build artifacts and reports"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  dev-deps      - Install development dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  security      - Run security scan"
	@echo "  security-json - Run security scan with JSON output"
	@echo "  security-sarif- Run security scan with SARIF output"
	@echo "  security-ci   - Run security scan like CI pipeline"
	@echo "  setup-dev     - Setup development environment"
	@echo "  check         - Run all checks (fmt, vet, lint, test, security)"
	@echo "  dev           - Build development version with race detector"
	@echo "  help          - Show this help"