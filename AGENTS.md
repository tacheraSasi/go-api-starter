# AGENTS.md - go-api-starter

## Repo structure

```
cmd/
  core/       HTTP API (Gin) + asynq client + asynqmon dashboard
  minion/     asynq worker (no HTTP port, scale horizontally)
  overseer/   asynq cron scheduler (exactly one replica)
  seed/       idempotent data seeder

internals/
  config/         env loading + validation
  handlers/       HTTP handlers
  middlewares/    auth, CORS, logging, RBAC
  minions/        task type constants + asynq handlers
  models/         GORM models
  repositories/   data access
  services/       business logic
  utils/          response helpers

pkg/
  database/       GORM connection + migrations (core owns migrations)
  jwt/            token signing/verification
  logger/         logrus service logger + asynq logger adapter
  exceptions/
  styles/

docker/
  Dockerfile          core image
  Dockerfile.minion   minion image
  Dockerfile.overseer overseer image
  Dockerfile.dev      shared dev image
```

## Commands

- `make dev` - core with `go run`
- `make dev-watch` - core with air hot reload
- `make dev-minion` - minion with `go run`
- `make dev-overseer` - overseer with `go run`
- `make build-all` - build all three binaries into `bin/`
- `make run-core` / `make run-minion` / `make run-overseer`
- `make seed` - seed roles, permissions, and sample users
- `make test` - `go test -v ./...`
- `make lint` - requires `golangci-lint`
- `make docker-up` - compose stack (core, minion, overseer, db, redis)
- `make dev-up` - compose.dev stack with hot reload

## Architecture rules

- `core` handles HTTP only and enqueues background work. It never processes
  asynq tasks.
- `minion` consumes asynq tasks. It has no HTTP server and connects to the
  database for repository-backed handlers, but never runs migrations.
- `overseer` only registers cron schedules and enqueues tasks. Run one replica.
- Shared task type constants live in `internals/minions/tasks.go`. A task type
  must be registered in `cmd/minion/main.go` and, if recurring, scheduled in
  `cmd/overseer/main.go`.
- Keep `cmd/*` binaries thin: wire dependencies, register handlers, start and
  gracefully stop the service. Business logic belongs in `internals/services`
  or dedicated `internals/minions` handlers.
- Redis is the only communication path between the three services. There are no
  direct service-to-service HTTP calls.
- Cron schedules use the `TIMEZONE` config value (default `UTC`).
