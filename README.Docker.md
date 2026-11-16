# Docker

The stack has three application services plus PostgreSQL and Redis.

| Service    | Binary            | Purpose                                      |
| ---------- | ----------------- | -------------------------------------------- |
| `core`     | `cmd/core`        | HTTP API, asynq client, asynqmon dashboard   |
| `minion`   | `cmd/minion`      | asynq worker, consumes tasks from Redis      |
| `overseer` | `cmd/overseer`    | asynq cron scheduler, enqueues recurring jobs|
| `db`       | postgres:16       | application database                         |
| `redis`    | redis:7           | asynq broker                                 |

## Production-style stack

```bash
cp .env.example .env
docker compose up -d --build
```

The API is available at `http://localhost:8080` and Asynqmon at
`http://localhost:8080/monitoring`.

Scale workers independently:

```bash
docker compose up -d --scale minion=3
```

Do not scale `overseer`: every replica registers the same cron schedules and
tasks would be enqueued more than once.

## Development stack

```bash
docker compose -f compose.dev.yaml up --build
```

`core` runs with air for hot reload. `minion` and `overseer` run with `go run`.

## Building images directly

```bash
docker build -f docker/Dockerfile -t go-api-starter-core .
docker build -f docker/Dockerfile.minion -t go-api-starter-minion .
docker build -f docker/Dockerfile.overseer -t go-api-starter-overseer .
```

`core` and `minion` are built with CGO enabled so the SQLite driver works.
`overseer` does not touch the database and is built with CGO disabled as a
small non-root image.
