# syntax=docker/dockerfile:1

################################################################################
# core: HTTP API server (Gin) + asynq client + asynqmon dashboard
################################################################################

ARG GO_VERSION=1.25.0
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

# CGO is needed by the SQLite driver (mattn/go-sqlite3).
RUN apk add --no-cache gcc musl-dev

# Retry module downloads: proxy.golang.org occasionally aborts a stream
# mid-transfer and a transient blip should not fail the image build.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    for i in $(seq 1 5); do \
      go mod download && { echo "[go] modules downloaded via proxy"; exit 0; }; \
      echo "[go] download attempt $i/5 failed, retrying in 5s..."; \
      sleep 5; \
    done; \
    for i in $(seq 1 3); do \
      GOPROXY=direct go mod download && { echo "[go] modules downloaded via direct"; exit 0; }; \
      echo "[go] direct attempt $i/3 failed, retrying in 5s..."; \
      sleep 5; \
    done; \
    echo "[go] go mod download failed after all retries"; \
    exit 1

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=1 GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o /bin/core ./cmd/core

################################################################################
# Runtime
################################################################################
FROM alpine:3.21 AS final

RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    libc6-compat \
    && update-ca-certificates

RUN mkdir -p /app/logs /app/data

WORKDIR /app

COPY --from=build /bin/core /bin/core

ENV GIN_MODE=release \
    SERVER_PORT=8080 \
    LOG_FILE_PATH=/app/logs/app.log \
    DB_TYPE=postgres

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/bin/core"]
