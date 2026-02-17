# go-api-starter

Production-ready Go API starter built with Gin + GORM, including JWT auth, RBAC, customer/invoice modules, and Docker support.

## Features

- Gin-based HTTP API with clean repository/service/handler layering
- JWT authentication (`login`, `register`, `logout`)
- Role-based access control (users, roles, permissions)
- Customer and invoice modules
- SQLite/Postgres/MySQL support via GORM
- Structured logging with Logrus
- Health and readiness endpoints

## Quick Start

### 1) Prerequisites

- Go 1.25+
- Make (optional but recommended)

### 2) Configure environment

```bash
cp .env.example .env
```

Adjust values in `.env` (especially `JWT_SECRET`).

### 3) Run locally

```bash
make dev
```

API runs on `http://localhost:8080` by default.

### 4) Seed sample data (optional)

```bash
make seed
```

## Useful Commands

```bash
make help
make dev
make build
make run
make test
make lint
make docker-up
```

## Health Endpoints

- `GET /health` - process liveness
- `GET /health/ready` - readiness + database connectivity

## API Route Overview

Base path: `/api/v1`

- Public: `/login`, `/register`
- Authenticated: `/logout`, `/users/*`, `/customers/*`, `/invoices/*`
- Admin: `/admin/users/*`, `/admin/roles/*`, `/admin/permissions/*`

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8080` | HTTP server port |
| `GIN_MODE` | `release` | Gin mode (`debug`, `release`, `test`) |
| `JWT_SECRET` | `secret` | JWT signing secret (change in production) |
| `JWT_EXPIRES_IN` | `24` | Token expiry in hours |
| `CORS_ALLOWED_ORIGINS` | `*` | Comma-separated origins or `*` |
| `LOG_FILE_PATH` | `logs/app.log` | Log file output |
| `DB_TYPE` | `sqlite` | `sqlite`, `postgres`, or `mysql` |
| `DB_PATH` | `core.db` | SQLite database file path |
| `DB_HOST` | `localhost` | DB host for Postgres/MySQL |
| `DB_PORT` | `5432` | DB port for Postgres/MySQL |
| `DB_USER` | `user` | DB user for Postgres/MySQL |
| `DB_PASSWORD` | `password` | DB password for Postgres/MySQL |
| `DB_NAME` | `dbname` | DB name for Postgres/MySQL |

## Docker

Use:

```bash
make docker-up
```

Or see [README.Docker.md](README.Docker.md) for details.

## Project Structure

- `cmd/api` - API application entrypoint
- `cmd/seed` - seed command
- `internals/handlers` - HTTP handlers
- `internals/services` - business logic
- `internals/repositories` - data access
- `pkg/database` - DB setup and migration

## Notes

- Default startup auto-runs migrations.
- Admin routes require JWT + admin role middleware.
