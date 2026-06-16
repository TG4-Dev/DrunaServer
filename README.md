# Druna Server

Druna is a backend service for managing users, events, friends, and groups. It provides a RESTful API using Go (Golang), Gin, and PostgreSQL, with support for JWT authentication and Swagger-based documentation.

## Features

- User registration and login
- JWT-based authentication (bcrypt password hashing)
- Telegram WebApp authentication
- Event creation, listing, deletion, and free-time calculation
- Friendship management (request, accept, reject, list, incoming/outgoing)
- Groups: create, list, details, add members
- CORS and auth rate limiting
- Swagger API documentation
- CI (GitHub Actions) and unit tests

## Tech Stack

- Go (Golang)
- Gin Web Framework
- PostgreSQL
- JWT (Authorization)
- Swaggo (Swagger documentation)

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/DrunaServer.git
cd DrunaServer
```

### 2. Install the dependencies

```bash
go mod tidy
```

### 3. Setup configuration

```bash
cp configs/config.yaml.example configs/config.yaml
cp .env.example .env
```

Edit `.env`:

```
DB_PASSWORD=yourpassword
JWT_SECRET=your-secret-key-at-least-32-chars
BOT_TOKEN=          # required for Telegram auth
CORS_ORIGINS=*      # optional, comma-separated origins
```

Database settings are in `configs/config.yaml` (host, port, dbname).

### 4. Run the server

```bash
go run cmd/main.go
```

## Docker

```bash
make build
make up
```

Migrations run automatically when the app container starts. App: `http://localhost:22000`

Manual migration (optional):

```bash
make migrate-up
```

## Testing

```bash
make test    # unit tests
make lint    # golangci-lint
make smoke   # curl smoke test (server must be running on BASE_URL)
```

## API Documentation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go
```

Swagger UI: `http://localhost:8000/swagger/index.html` (or port from config)

Documentation models are defined in `pkg/model/structs_doc.go`.

## Authentication

Use the JWT access token from `/auth/sign-in` in the `Authorization` header for all `/api/*` requests:

```
Authorization: Bearer <accessToken>
```

Renew tokens via `POST /auth/renew-token` with either:

```json
{ "refreshToken": "<refresh token>" }
```

or `Authorization: Bearer <refresh token>` header.

Telegram WebApp auth: `POST /auth/telegram` with `{ "initData": "<telegram init data>" }`.

**Note:** Password hashing uses bcrypt. Existing users created before this change must re-register.

## Endpoints

| Method | Path                          | Auth | Description                   |
| ------ | ----------------------------- | ---- | ----------------------------- |
| GET    | /ping/                        | No   | Health check                  |
| POST   | /auth/sign-up                 | No   | Register a new user           |
| POST   | /auth/sign-in                 | No   | Login and get JWT             |
| POST   | /auth/renew-token             | No   | Renew JWT tokens              |
| POST   | /auth/telegram                | No   | Telegram WebApp login         |
| GET    | /api/events/                  | Yes  | List user events              |
| POST   | /api/events/                  | Yes  | Create an event               |
| DELETE | /api/events/:id               | Yes  | Delete an event               |
| POST   | /api/events/free-time         | Yes  | Free time slots for a day     |
| GET    | /api/friends/list             | Yes  | List friends                  |
| GET    | /api/friends/request-list     | Yes  | All pending requests          |
| GET    | /api/friends/requests/incoming| Yes  | Incoming friend requests      |
| GET    | /api/friends/requests/outgoing| Yes  | Outgoing friend requests      |
| POST   | /api/friends/request          | Yes  | Send a friend request         |
| POST   | /api/friends/accept           | Yes  | Accept a friend request       |
| POST   | /api/friends/reject           | Yes  | Reject a friend request       |
| DELETE | /api/friends/                 | Yes  | Remove a friend               |
| POST   | /api/groups/create            | Yes  | Create a group                |
| GET    | /api/groups/list              | Yes  | List user groups              |
| GET    | /api/groups/:id               | Yes  | Group details with members    |
| POST   | /api/groups/:id/members       | Yes  | Add member (owner only)       |

See [AGENTS.md](AGENTS.md) for development conventions.
