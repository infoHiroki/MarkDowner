# MarkDowner Build Makefile

# Variables
BINARY_NAME=markdowner
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=${VERSION}"

# Default target
.PHONY: all
all: clean build

# Clean build artifacts
.PHONY: clean
clean:
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}.exe
	rm -f ${BINARY_NAME}-linux
	rm -f ${BINARY_NAME}-darwin
	rm -rf dist/

# Build for current platform
.PHONY: build
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for all platforms
.PHONY: build-all
build-all: clean
	mkdir -p dist
	# Windows
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .
	# Linux
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	# macOS Intel
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	# macOS Apple Silicon
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .

# Run tests
.PHONY: test
test:
	go test ./... -v

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	go vet ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Run the application
.PHONY: run
run: build
	./${BINARY_NAME}

# Development setup
.PHONY: dev
dev: tidy fmt lint test build

# Release preparation
.PHONY: release
release: clean test lint build-all

# Install locally
.PHONY: install
install:
	go install ${LDFLAGS} .

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         - Clean and build for current platform"
	@echo "  build       - Build for current platform"
	@echo "  build-all   - Build for all platforms"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"
	@echo "  tidy        - Tidy dependencies"
	@echo "  run         - Build and run the application"
	@echo "  dev         - Full development cycle"
	@echo "  release     - Prepare release builds"
	@echo "  install     - Install to GOPATH/bin"
	@echo "  help        - Show this help"