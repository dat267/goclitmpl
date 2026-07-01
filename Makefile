.PHONY: all build run test test-coverage lint fmt clean help cross-compile

# Build parameters
BINARY_NAME=goclitmpl
BUILD_DIR=bin
MAIN_PATH=./cmd/goclitmpl

# Versioning (derived from Git or fallback defaults)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags to inject variables and optimize size (strip debugging symbols)
LDFLAGS=-ldflags="-s -w \
	-X github.com/dat267/goclitmpl/internal/cli.Version=$(VERSION) \
	-X github.com/dat267/goclitmpl/internal/cli.Commit=$(COMMIT) \
	-X github.com/dat267/goclitmpl/internal/cli.Date=$(DATE)"

all: fmt lint test build

build:
	@echo "==> Building binary to $(BUILD_DIR)/$(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH) $(filter-out $@,$(MAKECMDGOALS))

test:
	@echo "==> Running unit tests..."
	go test -v -race ./...

test-coverage:
	@echo "==> Running tests with coverage..."
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"

lint:
	@echo "==> Running golangci-lint..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run; \
	else \
		echo "Warning: golangci-lint is not installed. Skipping. Install via: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	}

fmt:
	@echo "==> Formatting code..."
	go fmt ./...

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf $(BUILD_DIR) coverage.out coverage.html

cross-compile:
	@echo "==> Cross-compiling for Linux, macOS, and Windows..."
	@mkdir -p $(BUILD_DIR)
	# Linux 64-bit
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	# macOS Apple Silicon
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	# macOS Intel
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	# Windows 64-bit
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

help:
	@echo "Available make targets:"
	@echo "  build           - Build optimized executable for host platform"
	@echo "  run             - Build and run program"
	@echo "  test            - Run all unit tests with race detection"
	@echo "  test-coverage   - Run tests with coverage and generate HTML report"
	@echo "  lint            - Run static analysis via golangci-lint"
	@echo "  fmt             - Format all Go files"
	@echo "  clean           - Remove build directories and reports"
	@echo "  cross-compile   - Compile binaries for multiple platforms"
	@echo "  all             - Format, lint, test, and build binary"
