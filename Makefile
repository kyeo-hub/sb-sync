.PHONY: build test lint clean install run help

VERSION := $(shell cat VERSION)
BINARY_NAME := sb-sync
BUILD_DIR := dist
INSTALL_DIR := /usr/local/bin

BUILD_LDFLAGS := -s -w -X sb-sync/cmd.Version=$(VERSION)

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@if [ "$$(uname)" = "Windows_NT" ]; then \
		mv $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/$(BINARY_NAME).exe; \
	fi
	@echo "Built successfully!"

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows

build-linux-amd64:
	@echo "Building Linux amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

build-linux-arm64:
	@echo "Building Linux arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

build-darwin-amd64:
	@echo "Building Darwin amd64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

build-darwin-arm64:
	@echo "Building Darwin arm64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-windows:
	@echo "Building Windows amd64..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME).exe .

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Coverage:"
	@go tool cover -func=coverage.out | tail -1

test-coverage:
	@echo "Generating coverage report..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@echo "Running linters..."
	go vet ./...
	@echo "Checking formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "Linting passed!"

fmt:
	@echo "Formatting code..."
	gofmt -w .

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Cleaned!"

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	@if [ -w $(INSTALL_DIR) ]; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/; \
		chmod +x $(INSTALL_DIR)/$(BINARY_NAME); \
	else \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/; \
		sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@echo "Installed to $(INSTALL_DIR)/$(BINARY_NAME)"

run:
	@echo "Running $(BINARY_NAME)..."
	go run .

deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod download

check: lint test

help:
	@echo "sb-sync Makefile targets:"
	@echo ""
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms"
	@echo "  test           - Run tests with coverage"
	@echo "  test-coverage  - Generate HTML coverage report"
	@echo "  lint           - Run linters and formatters"
	@echo "  fmt            - Format code"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install binary to system"
	@echo "  run            - Run without building"
	@echo "  deps           - Update dependencies"
	@echo "  check          - Run lint and tests"
	@echo "  help           - Show this help"
