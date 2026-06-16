# Druna Frontend API Guide

Contract for frontend / mobile / Telegram Mini App clients integrating with **DrunaServer**.

Backend repo: [DrunaServer](https://github.com/TG4-Dev/DrunaServer)  
Related docs: [README.md](../README.md) · [AGENTS.md](../AGENTS.md)

---

## Base URLs

| Environment | Base URL |
|-------------|----------|
| Local dev | `http://localhost:8000` |
| Docker (`make up`) | `http://localhost:22000` |
| Swagger UI | `{base}/swagger/index.html` |

Use **`/api/v1/...`** for all protected routes. Legacy **`/api/...`** mirrors the same handlers but v1 is preferred.

---

## Response envelope

Every JSON response follows this shape.

**Success:**

```json
{
  "data": { },
  "error": null
}
```

**Error:**

```json
{
  "data": null,
  "error": {
    "message": "invalid credentials",
    "code": 401
  }
}
```

On the client, read the payload from **`response.data.data`**, not from the HTTP body root.

Handle errors via `response.data.error.message` and `response.data.error.code`.

---

## Headers

| Header | When |
|--------|------|
| `Content-Type: application/json` | All POST/PATCH bodies |
| `Authorization: Bearer <accessToken>` | All `/api/v1/*` routes |
| `X-Request-ID` | Optional; server echoes generated ID |

---

## Authentication

### Token types

| Token | Use |
|-------|-----|
| `accessToken` | All `/api/v1/*` requests |
| `refreshToken` | Only `POST /auth/renew-token` |

The auth middleware **rejects refresh tokens** on API routes (401).

Refresh tokens are **rotated**: after renew, the old refresh token is revoked. Always store the new pair.

### Sign up

`POST /auth/sign-up`

```json
{
  "name": "Alice",
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret123"
}
```

Legacy field `passwordHash` is also accepted (same value as plaintext password).

**Response `data`:**

```json
{ "id": 1 }
```

### Sign in

`POST /auth/sign-in`

```json
{
  "username": "alice",
  "password": "secret123"
}
```

**Response `data`:**

```json
{
  "accessToken": "eyJ...",
  "refreshToken": "eyJ..."
}
```

### Renew tokens

`POST /auth/renew-token`

```json
{ "refreshToken": "eyJ..." }
```

Alternative: `Authorization: Bearer <refreshToken>` header (body field can be omitted).

**Response `data`:** new `accessToken` + `refreshToken`.

### Telegram WebApp login

`POST /auth/telegram`

```json
{ "initData": "<Telegram.WebApp.initData string>" }
```

**Response `data`:** same as sign-in (`accessToken`, `refreshToken`).

Requires matching `BOT_TOKEN` on the server.

### Recommended client flow

```
1. Sign in → store accessToken + refreshToken
2. API call with Authorization: Bearer <accessToken>
3. On 401 → POST /auth/renew-token → retry once with new accessToken
4. If renew fails → clear tokens, redirect to login
```

---

## Events

All routes require auth. Prefix: `/api/v1/events`.

### List events

`GET /api/v1/events/`

Query params (all optional):

| Param | Type | Description |
|-------|------|-------------|
| `limit` | int | Page size (default 50) |
| `offset` | int | Pagination offset |
| `type` | string | Filter by event type |
| `dateFrom` | RFC3339 | Events ending on or after |
| `dateTo` | RFC3339 | Events starting on or before |

**Response `data`:**

```json
{
  "events": [
    {
      "eventID": 1,
      "userID": 1,
      "title": "Meeting",
      "startTime": "2026-06-17T10:00:00Z",
      "endTime": "2026-06-17T11:00:00Z",
      "type": "work"
    }
  ],
  "total": 1,
  "limit": 50,
  "offset": 0
}
```

### Create event

`POST /api/v1/events/`

```json
{
  "title": "Meeting",
  "startTime": "2026-06-17T10:00:00Z",
  "endTime": "2026-06-17T11:00:00Z",
  "type": "work"
}
```

**Response `data`:** `{ "eventId": 1 }`

Validation errors (400):

- `end time must be after start time`
- `event overlaps with an existing event`

### Update event

`PATCH /api/v1/events/:id`

Same body as create. `userID` is taken from the token.

**Response `data`:** `{ "message": "event updated" }`

### Delete event

`DELETE /api/v1/events/:id`

**Response `data`:** `{ "message": "event deleted" }`

### Personal free time

`POST /api/v1/events/free-time`

```json
{ "date": "2026-06-17" }
```

Date format: **`YYYY-MM-DD`**.

**Response `data`:**

```json
{
  "freeSlots": [
    {
      "start": "2026-06-17T00:00:00Z",
      "end": "2026-06-17T10:00:00Z"
    },
    {
      "start": "2026-06-17T11:00:00Z",
      "end": "2026-06-18T00:00:00Z"
    }
  ]
}
```

---

## Friends

Prefix: `/api/v1/friends`. All routes require auth.

| Method | Path | Body | Response `data` |
|--------|------|------|-----------------|
| GET | `/list` | — | `{ "friends": [FriendInfo] }` |
| GET | `/search?username=ali` | — | `{ "users": [FriendInfo] }` |
| GET | `/request-list` | — | `{ "friends": [FriendInfo] }` (all pending) |
| GET | `/requests/incoming` | — | `{ "friends": [FriendInfo] }` |
| GET | `/requests/outgoing` | — | `{ "friends": [FriendInfo] }` |
| POST | `/request` | `{ "username": "bob" }` | `{ "message": "friend request sent" }` |
| POST | `/accept` | `{ "username": "bob" }` | `{ "message": "friend request accepted" }` |
| POST | `/reject` | `{ "username": "bob" }` | `{ "message": "friend request rejected" }` |
| DELETE | `/` | `{ "username": "bob" }` | `{ "message": "friend deleted", "username": "bob" }` |

**FriendInfo:**

```json
{
  "id": 2,
  "name": "Bob",
  "username": "bob"
}
```

Business rules (400 errors):

- Cannot send request to yourself
- Cannot re-send if pending, already friends, or previously rejected

---

## Groups

Prefix: `/api/v1/groups`. All routes require auth.

| Method | Path | Body | Notes |
|--------|------|------|-------|
| POST | `/create` | `{ "name": "Weekend trip" }` | Owner = current user |
| GET | `/list` | — | Groups user owns or belongs to |
| GET | `/:id` | — | Details + members |
| POST | `/:id/members` | `{ "username": "bob" }` | Owner only |
| POST | `/:id/confirm` | `{ "confirmedTime": "2026-06-20T18:00:00Z" }` | Member confirms time |
| POST | `/:id/free-time` | `{ "date": "2026-06-20" }` | Intersection of member slots |
| POST | `/:id/leave` | — | Owner cannot leave (must delete) |
| DELETE | `/:id` | — | Owner only |

**Create response `data`:** `{ "message": "group created", "groupId": 1 }`

**List response `data`:** `{ "groups": [Group] }`

**Group object:**

```json
{
  "groupID": 1,
  "ownerID": 1,
  "name": "Weekend trip",
  "confirmedTime": "2026-06-20T18:00:00Z"
}
```

**Details response `data`:**

```json
{
  "groupID": 1,
  "ownerID": 1,
  "name": "Weekend trip",
  "confirmedTime": "2026-06-20T18:00:00Z",
  "members": [
    {
      "id": 1,
      "name": "Alice",
      "username": "alice",
      "confirmedTime": "2026-06-20T18:00:00Z"
    }
  ]
}
```

**Group free-time response `data`:** same shape as personal free time (`freeSlots`).

---

## Public utility routes

| Method | Path | Response `data` |
|--------|------|-----------------|
| GET | `/ping/` | `{ "status": "ok", "db": "ok" }` — `503` if DB down |
| GET | `/metrics` | Prometheus format (not JSON envelope) |
| GET | `/swagger/*` | Swagger UI |

---

## Dates and timezones

- Event timestamps: **ISO 8601 / RFC3339** (e.g. `2026-06-17T10:00:00Z`)
- Free-time date param: **`YYYY-MM-DD`** (day interpreted in server/local parsing; send UTC dates for consistency)
- Display times in the user's local timezone on the client

---

## CORS

Server reads `CORS_ORIGINS` env var (comma-separated). Default: `*`.

For cookie-based auth set explicit origins on the backend. Current API uses Bearer tokens in headers — `*` works for dev.

---

## Rate limits

`/auth/*`: **30 requests/minute** per IP → `429` with `{ "error": { "message": "rate limit exceeded" } }`.

---

## Suggested MVP screens

1. **Auth** — login, register, Telegram login (if Mini App)
2. **Calendar / Events** — list, create, edit, delete
3. **Free time** — day picker + slot list (personal)
4. **Friends** — search, incoming/outgoing tabs, accept/reject
5. **Groups** — list, create, detail (members), confirm time, group free-time

---

## TypeScript client sketch

```typescript
type ApiResponse<T> = {
  data: T | null;
  error: { message: string; code: number } | null;
};

async function api<T>(
  path: string,
  options: RequestInit = {},
  accessToken?: string
): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
      ...options.headers,
    },
  });
  const body: ApiResponse<T> = await res.json();
  if (body.error) throw new Error(body.error.message);
  return body.data as T;
}
```

---

## Testing against local backend

```bash
# terminal 1
cp configs/config.yaml.example configs/config.yaml
cp .env.example .env
go run cmd/main.go

# terminal 2
make smoke
# or BASE_URL=http://localhost:8000 make smoke
```

Smoke script covers: sign-up → sign-in → event → free-time → group via `/api/v1/...`.

---

## Changelog notes for frontend

- All responses wrapped in `{ data, error }` envelope
- Prefer `password` over legacy `passwordHash` on sign-in/sign-up
- Use `/api/v1/` prefix
- Access and refresh tokens are distinct; only access works on API routes
- Refresh token rotation: always save new tokens after renew
