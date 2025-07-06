# Druna Server

Druna is a backend service for managing users, events, friends, and groups. It provides a RESTful API using Go (Golang), Gin, and PostgreSQL, with support for JWT authentication and Swagger-based documentation.

## Features

- User registration and login
- JWT-based authentication
- Event creation, listing, and deletion
- Friendship management (request, list)
- Group planning support (WIP)
- Swagger API documentation

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

### 3. Setup .env variables

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=druna
```

### 4. Run the server.

```bash
go run cmd/main.go
```

## API Documentation
```
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger docs
```bash
swag init -g cmd/main.go
```
- All Swagger documentation models are defined in model/structs_doc.go.

## Authentication
Use the JWT access token returned from /auth/sign-in in the Authorization header for all /api/* requests:

```
Authorization: Bearer <token>
```

Use the JWT refresh token and /auth/renew-token endpoint to get a new access token:

```
Authorization: Bearer <refresh token>
```

## Endpoints

| Method | Path                  | Description           |
| ------ | --------------------- | --------------------- |
| POST   | /auth/sign-up         | Register a new user   |
| POST   | /auth/sign-in         | Login and get JWT     |
| POST   | /auth/renew-token     | Renew JWT access token|
| POST   | /api/events/add-event | Create a new event    |
| POST   | /api/events/list      | List user events      |
| DELETE | /api/events/\:id      | Delete an event       |
| GET    | /api/friends/list     | List user's friends   |
| POST   | /api/friends/request  | Send a friend request |
