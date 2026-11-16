APP_NAME=go-api-starter
BIN_DIR=bin
CORE_BIN=$(BIN_DIR)/core
MINION_BIN=$(BIN_DIR)/minion
OVERSEER_BIN=$(BIN_DIR)/overseer
SEED_BIN=$(BIN_DIR)/seed

.PHONY: help dev dev-watch dev-minion dev-overseer build build-core build-minion build-overseer build-all \
	run run-core run-minion run-overseer seed test test-coverage lint fmt clean deps \
	security-check docker-build docker-build-core docker-build-minion docker-build-overseer \
	docker-build-all docker-up docker-down docker-logs dev-up dev-down

# Show help
help:
	@echo "Available commands:"
	@echo "  dev              - Run core (API) locally"
	@echo "  dev-watch        - Run core with live reload via air"
	@echo "  dev-minion       - Run minion (asynq worker) locally"
	@echo "  dev-overseer     - Run overseer (cron scheduler) locally"
	@echo "  build-core       - Build bin/core (HTTP API)"
	@echo "  build-minion     - Build bin/minion (asynq worker)"
	@echo "  build-overseer   - Build bin/overseer (cron scheduler)"
	@echo "  build-all        - Build all three services"
	@echo "  run-core         - Build and run core"
	@echo "  run-minion       - Build and run minion"
	@echo "  run-overseer     - Build and run overseer"
	@echo "  seed             - Seed database with initial data"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  clean            - Clean build artifacts"
	@echo "  deps             - Download dependencies"
	@echo "  security-check   - Run security checks"
	@echo "  docker-build-all - Build all three Docker images"
	@echo "  docker-up        - Start the production-style stack"
	@echo "  docker-down      - Stop the production-style stack"
	@echo "  dev-up           - Start the development stack (compose.dev.yaml)"
	@echo "  dev-down         - Stop the development stack"

# Development mode
dev:
	@echo "Starting core..."
	go run ./cmd/core

# Hot reload via air (installed by docker/Dockerfile.dev or `make install-tools`)
dev-watch:
	@echo "Starting core with air..."
	air -c .air.toml

dev-minion:
	@echo "Starting minion..."
	go run ./cmd/minion

dev-overseer:
	@echo "Starting overseer..."
	go run ./cmd/overseer

# Install/update dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Builds
build: build-all

build-core:
	@echo "Building core..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-w -s" -o $(CORE_BIN) ./cmd/core

build-minion:
	@echo "Building minion..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-w -s" -o $(MINION_BIN) ./cmd/minion

build-overseer:
	@echo "Building overseer..."
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-w -s" -o $(OVERSEER_BIN) ./cmd/overseer

build-all: build-core build-minion build-overseer

# Run built binaries
run: run-core

run-core: build-core
	@echo "Running core..."
	./$(CORE_BIN)

run-minion: build-minion
	@echo "Running minion..."
	./$(MINION_BIN)

run-overseer: build-overseer
	@echo "Running overseer..."
	./$(OVERSEER_BIN)

# Seed database
seed:
	@echo "Seeding database..."
	@mkdir -p $(BIN_DIR)
	go run ./cmd/seed

# Tests
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint / format
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed. Install it with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html
	go clean

security-check:
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install it with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

# Docker
docker-build: docker-build-all

docker-build-core:
	docker build -f docker/Dockerfile -t $(APP_NAME)-core .

docker-build-minion:
	docker build -f docker/Dockerfile.minion -t $(APP_NAME)-minion .

docker-build-overseer:
	docker build -f docker/Dockerfile.overseer -t $(APP_NAME)-overseer .

docker-build-all: docker-build-core docker-build-minion docker-build-overseer

docker-up:
	@echo "Starting stack (core + minion + overseer)..."
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

dev-up:
	docker compose -f compose.dev.yaml up --build

dev-down:
	docker compose -f compose.dev.yaml down

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
