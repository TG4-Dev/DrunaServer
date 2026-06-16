# AGENTS.md — DrunaServer

Backend REST API for the Druna app: users, events, friends, and groups.

## Tech stack

- Go 1.24+, Gin, PostgreSQL (sqlx)
- JWT auth, Swagger (swaggo), Viper + godotenv for config
- GitHub Actions CI, golangci-lint

## Architecture

```
HTTP request → pkg/handler → pkg/service → pkg/repository → PostgreSQL
```

| Layer        | Path               | Responsibility                          |
|--------------|--------------------|-----------------------------------------|
| Entry        | `cmd/main.go`      | Config, DB, graceful shutdown           |
| HTTP         | `pkg/handler/`     | Routes, middleware, JSON responses      |
| Business     | `pkg/service/`     | JWT, password hashing, domain rules     |
| Data         | `pkg/repository/`  | SQL via sqlx                            |
| Models       | `pkg/model/`       | Domain structs; Swagger types in `structs_doc.go` |

## Quick start (local)

```bash
go mod tidy
cp configs/config.yaml.example configs/config.yaml
cp .env.example .env
go run cmd/main.go
```

Server listens on the port from `configs/config.yaml` (default `8000`).

## Environment variables

| Variable       | Required | Description                          |
|----------------|----------|--------------------------------------|
| `DB_PASSWORD`  | Yes      | PostgreSQL password                  |
| `JWT_SECRET`   | Yes      | HMAC key for signing JWT tokens      |
| `BOT_TOKEN`    | Telegram | Telegram bot token for WebApp auth   |
| `CORS_ORIGINS` | No       | Comma-separated allowed origins (default `*`) |
| `DATABASE_URL` | Docker   | Used by entrypoint for auto-migrations |

## Adding a new endpoint

1. Define method on repository interface in `pkg/repository/repository.go`
2. Implement in `pkg/repository/*_postgres.go`
3. Add service method in `pkg/service/` and interface in `service.go`
4. Add handler in `pkg/handler/` and register route in `pkg/handler/handler.go`
5. Add Swagger annotations; run `swag init -g cmd/main.go`
6. Update `README.md` endpoint table if the route is public

## Conventions

- Register routes only in `pkg/handler/handler.go`
- Return API errors via `NewErrorResponse()` from `pkg/handler/response.go`
- Authenticated user ID comes from JWT middleware (`userCtx`); never trust query/body for identity
- Passwords are hashed with bcrypt on sign-up; clients send plaintext in the `passwordHash` JSON field (legacy name)
- Token response keys: `accessToken`, `refreshToken`
- `/auth/*` routes are rate-limited (30 req/min per IP)

## Swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go
```

## Docker

```bash
make build
make up
```

Migrations run automatically on container start via `docker-entrypoint.sh`.

App: `http://localhost:22000`. Postgres: `localhost:5432`, database `druna_db`.

## Testing

```bash
make test          # unit tests
make lint          # golangci-lint
make smoke         # end-to-end curl script (server must be running)
bash scripts/smoke_test.sh
```

## Known limitations

- Password hashing uses bcrypt without DB migration — existing users must re-register
- Group scheduling (confirmed time negotiation) is not implemented yet

## Pre-PR checklist

- [ ] `make test`
- [ ] `go vet ./...`
- [ ] `make lint` (or CI green)
- [ ] Swagger regenerated if API changed
- [ ] No secrets logged to stdout
