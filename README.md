# go-api-starter

Production-ready Go API starter built with Gin + GORM, including JWT auth, RBAC, and a three-service background job architecture.

## Architecture

The template is split into three independent binaries that share the same
`internals` and `pkg` packages:

```
cmd/
  core/       HTTP API (Gin) + asynq client + asynqmon dashboard
  minion/     asynq worker, consumes tasks from Redis
  overseer/   asynq cron scheduler, enqueues recurring tasks
  seed/       idempotent database seeder

internals/
  config/         environment loading + validation
  handlers/       HTTP handlers
  middlewares/    auth, CORS, logging, RBAC
  minions/        async task types + task handlers
  models/         GORM models
  repositories/   data access
  services/       business logic

pkg/
  database/   GORM connection + migrations
  jwt/        token signing/verification
  logger/     logrus loggers + asynq logger adapter
  exceptions/ shared error types
  styles/     UI styles
```

Request flow:

```
HTTP  ->  handler  ->  service  ->  repository  ->  DB
core  ->  asynq client  ->  Redis  ->  minion  ->  handler
overseer  ->  cron schedule  ->  Redis  ->  minion
```

`core` never processes background jobs. `minion` has no HTTP server and can be
scaled horizontally. `overseer` only enqueues tasks, so run exactly one replica
per environment.

## Features

- Gin HTTP API with clean repository/service/handler layering
- JWT authentication (`login`, `register`, `logout`)
- Role-based access control (users, roles, permissions)
- SQLite/Postgres/MySQL support via GORM
- Asynq task queue with dedicated worker and scheduler services
- Asynqmon monitoring dashboard at `/monitoring`
- Structured logging with Logrus
- Health and readiness endpoints
- Production and development Docker Compose stacks

## Quick Start

### 1) Prerequisites

- Go 1.25+
- Redis (for the minion/overseer)
- Make (optional but recommended)

### 2) Configure environment

```bash
cp .env.example .env
```

Adjust values in `.env` (especially `JWT_SECRET` and `REDIS_URL`).

### 3) Run locally

In separate terminals:

```bash
make dev           # core API
make dev-minion    # asynq worker
make dev-overseer  # cron scheduler
```

Use `make dev-watch` instead of `make dev` to run core with air hot reload.

The API runs on `http://localhost:8080`. Asynqmon is available at
`http://localhost:8080/monitoring` (using `MONITORING_USER` / `MONITORING_PASSWORD`).

### 4) Seed sample data (optional)

```bash
make seed
```

### 5) Run everything with Docker

```bash
make docker-up      # production-style stack
make dev-up         # development stack
```

## Useful Commands

```bash
make help
make dev
make dev-watch
make dev-minion
make dev-overseer
make build-all
make run-core
make run-minion
make run-overseer
make seed
make test
make lint
make docker-up
make dev-up
```

## Health Endpoints

- `GET /health` - process liveness
- `GET /health/ready` - readiness + database connectivity

## API Route Overview

Base path: `/api/v1`

- Public: `/login`, `/register`
- Authenticated: `/logout`, `/users/*`
- Admin: `/admin/users/*`, `/admin/roles/*`, `/admin/permissions/*`

## Adding a background task

1. Add a task type constant in `internals/minions/tasks.go`.
2. Write the handler in `internals/minions` (or delegate to a service).
3. Register the handler in `cmd/minion/main.go` with `mux.HandleFunc`.
4. Enqueue it from `core` via `services.NotificationService` or a new service.
5. Register a cron schedule in `cmd/overseer/main.go` if it is recurring.

The included `maintenance:heartbeat` task is only a template placeholder.

## Environment Variables

| Variable               | Default                  | Description                               |
| ---------------------- | ------------------------ | ----------------------------------------- |
| `SERVER_PORT`          | `8080`                   | HTTP server port                          |
| `GIN_MODE`             | `release`                | Gin mode (`debug`, `release`, `test`)     |
| `TIMEZONE`             | `UTC`                    | Timezone for logs and cron schedules      |
| `JWT_SECRET`           | `secret`                 | JWT signing secret (change in production) |
| `JWT_EXPIRES_IN`       | `24`                     | Token expiry in hours                     |
| `CORS_ALLOWED_ORIGINS` | `*`                      | Comma-separated origins or `*`            |
| `LOG_FILE_PATH`        | `logs/app.log`           | Log file output                           |
| `DB_TYPE`              | `sqlite`                 | `sqlite`, `postgres`, or `mysql`          |
| `DB_PATH`              | `core.db`                | SQLite database file path                 |
| `DATABASE_URL`         | empty                    | Full DSN, takes precedence for postgres   |
| `DB_HOST`              | `localhost`              | DB host for Postgres/MySQL                |
| `DB_PORT`              | `5432`                   | DB port for Postgres/MySQL                |
| `DB_USER`              | `user`                   | DB user for Postgres/MySQL                |
| `DB_PASSWORD`          | `password`               | DB password for Postgres/MySQL            |
| `DB_NAME`              | `dbname`                 | DB name for Postgres/MySQL                |
| `DB_SSL_MODE`          | `disable`                | Postgres SSL mode                         |
| `REDIS_URL`            | `redis://localhost:6379` | Redis URL used by core/minion/overseer    |
| `MONITORING_USER`      | `admin`                  | Asynqmon basic auth user                  |
| `MONITORING_PASSWORD`  | `changeme`               | Asynqmon basic auth password              |

## Docker

See [README.Docker.md](README.Docker.md) for details.

## Notes

- Default startup auto-runs migrations in `core`.
- Admin routes require JWT + admin role middleware.
- `minion` connects to the database for repository-backed tasks but does not
  run migrations. `core` owns schema changes.
