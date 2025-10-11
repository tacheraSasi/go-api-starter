APP_NAME=go-api-starter
BINARY_NAME=bin/api
DOCKER_IMAGE_NAME=go-api-starter
DOCKER_CONTAINER_NAME=go-api-starter-container

.PHONY: help dev build run test lint fmt clean docker-build docker-run docker-up docker-down deps security-check

# Show help
help:
	@echo "Available commands:"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  seed          - Seed database with initial data"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  security-check- Run security checks"
	@echo "  install-tools - Install development tools"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-up     - Start with docker-compose"
	@echo "  docker-down   - Stop docker-compose"

# Development mode
dev:
	@echo "Starting development server..."
	go run ./cmd/api/main.go

# Install/update dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p bin
	go build -ldflags="-w -s" -o $(BINARY_NAME) ./cmd/api/main.go

# Build and run
run: build
	@echo "Running $(APP_NAME)..."
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Install it with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Security check
security-check:
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

# Docker compose up
docker-up:
	@echo "Starting services with docker-compose..."
	docker-compose up --build

# Docker compose down
docker-down:
	@echo "Stopping services..."
	docker-compose down

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Seed database
seed:
	@echo "Seeding database..."
	go run ./cmd/seed/main.go