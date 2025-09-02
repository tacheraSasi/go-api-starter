DEV_BINARY=bin/invoice-api

.PHONY: dev build run test lint

dev:
	go run ./cmd/api/main.go

build:
	go build -o $(DEV_BINARY) ./cmd/api/main.go

run: build
	./$(DEV_BINARY)

test:
	go test ./...

lint:
	golangci-lint run || echo "Install golangci-lint for linting support."