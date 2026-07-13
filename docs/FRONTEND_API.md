# Druna Frontend API Guide

Contract for frontend / mobile / **Telegram Mini App** clients integrating with **DrunaServer**.

Backend repo: [DrunaServer](https://github.com/TG4-Dev/DrunaServer)  
Related docs: [README.md](../README.md) · [AGENTS.md](../AGENTS.md)

---

## Table of contents

1. [Telegram Mini App (start here)](#telegram-mini-app)
2. [Base URLs](#base-urls)
3. [Response envelope](#response-envelope)
4. [Authentication](#authentication)
5. [Profile](#profile)
6. [Events](#events)
7. [Friends](#friends)
8. [Groups](#groups)
9. [Public utility routes](#public-utility-routes)
10. [Shared code: Web ↔ TMA ↔ Mobile](#shared-code-web--tma--mobile)
11. [TypeScript client](#typescript-client)
12. [Testing](#testing)

---

## Telegram Mini App

This section contains everything needed to build the Druna UI as a **Telegram Mini App (TMA)**.

### Architecture

```
User in Telegram
    ↓ opens Mini App (WebView)
React/Vue app (your frontend)
    ↓ POST /auth/telegram { initData }
DrunaServer → validates HMAC → JWT tokens
    ↓ Authorization: Bearer <accessToken>
DrunaServer /api/v1/* (events, friends, groups)
```

The **bot** is optional infrastructure:

- Menu Button / Web App URL → opens your Mini App
- Push notifications, `/start`, deep links (future)

The **Mini App** is a normal web app (HTML/JS) using Telegram WebApp SDK + Druna REST API.

### What the backend does on Telegram login

1. Receives raw `initData` string from the client
2. Validates HMAC signature using server `BOT_TOKEN` (same token as your bot)
3. Rejects initData older than `TELEGRAM_AUTH_TTL_HOURS` (default 24h) via `auth_date`
4. Parses `user` JSON from initData (`id`, `first_name`, `username`, `photo_url`, …)
4. Finds user by `telegram_id` in DB, or **auto-registers** a new account:
   - `username`: `@username` from Telegram, or `tg_{telegram_id}` if no username
   - `email`: `{username}@telegram.local`
   - `name`: `first_name` + `last_name`
5. Returns `accessToken` + `refreshToken` (same as password login)

**No separate sign-up screen is required in TMA** — first open = register, next opens = login.

### Prerequisites

| Item | Who sets it |
|------|-------------|
| Telegram bot | Create via [@BotFather](https://t.me/BotFather) |
| `BOT_TOKEN` | Server `.env` — **must match the bot** that opens the Mini App |
| Mini App URL | BotFather → Bot Settings → Menu Button / Web App |
| HTTPS | **Required in production** — Telegram opens WebView only on HTTPS URLs |

Server without `BOT_TOKEN` returns:

```json
{
  "data": null,
  "error": { "message": "telegram auth failed: BOT_TOKEN is not configured", "code": 401 }
}
```

### BotFather setup

1. `/newbot` → get token → put in server `.env`:
   ```
   BOT_TOKEN=123456789:ABCdefGHI...
   ```
2. Bot Settings → **Menu Button** → Configure → enter your app URL:
   ```
   https://your-domain.com/tma/
   ```
3. (Optional) `/setdomain` if using Login Widget on external sites

For **local development**, expose HTTPS via tunnel:

```bash
# example with ngrok
ngrok http 5173
# use https://xxxx.ngrok.io as Menu Button URL in BotFather
```

Point API base URL to your backend (see [Base URLs](#base-urls)).

### Recommended npm packages

```bash
npm install @twa-dev/sdk
# optional
npm install @telegram-apps/sdk
```

Official docs: [Telegram Mini Apps](https://core.telegram.org/bots/webapps)

### TMA bootstrap (React example)

```typescript
import WebApp from "@twa-dev/sdk";

// Call once on app start
WebApp.ready();
WebApp.expand();

// Apply Telegram theme to CSS variables
document.documentElement.style.setProperty(
  "--tg-theme-bg-color",
  WebApp.themeParams.bg_color ?? "#ffffff"
);
document.documentElement.style.setProperty(
  "--tg-theme-text-color",
  WebApp.themeParams.text_color ?? "#000000"
);
```

Use `WebApp.themeParams` for native look: `bg_color`, `text_color`, `button_color`, `hint_color`, etc.

### Authentication flow for Mini App

```typescript
const API_BASE = import.meta.env.VITE_API_URL; // e.g. https://api.druna.app

type Tokens = { accessToken: string; refreshToken: string };

async function loginWithTelegram(): Promise<Tokens> {
  const initData = WebApp.initData; // raw query string — DO NOT parse manually on client for auth

  if (!initData) {
    throw new Error("Open this app from Telegram, not in a regular browser");
  }

  const res = await fetch(`${API_BASE}/auth/telegram`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ initData }),
  });

  const body = await res.json();
  if (body.error) throw new Error(body.error.message);

  return body.data as Tokens;
}
```

**Important:**

- Send **`WebApp.initData`** as-is (query string format: `query_id=...&user=...&auth_date=...&hash=...`)
- Never trust `WebApp.initDataUnsafe.user.id` alone — always authenticate via server
- Do **not** implement HMAC validation on the client — server already does it

### Token storage in Mini App

| Storage | Pros | Cons |
|---------|------|------|
| `localStorage` | Simple | Cleared if WebView cache cleared |
| `Telegram.WebApp.CloudStorage` | Syncs in Telegram cloud | Async API, size limits |
| In-memory only | Safest from persistence | Re-auth on every cold start |

Recommended: **CloudStorage for refreshToken**, in-memory for accessToken, re-login via `initData` on cold start (initData is reissued each session anyway).

```typescript
// CloudStorage (promisify)
function cloudGet(key: string): Promise<string | null> {
  return new Promise((resolve) => {
    WebApp.CloudStorage.getItem(key, (err, value) => {
      resolve(err ? null : value ?? null);
    });
  });
}
```

### Session lifecycle (recommended)

```
App mount
  → WebApp.initData present?
      YES → POST /auth/telegram → store tokens → load app
      NO  → show "Open in Telegram" fallback (browser dev mode: use /auth/sign-in)

API call → 401
  → POST /auth/renew-token with refreshToken
  → retry once
  → still 401? → re-run /auth/telegram with fresh initData
```

`initData` is available on every Mini App open — **re-authenticating via `/auth/telegram` is cheap** and avoids stale refresh issues.

### API calls after login

Same as web — all protected routes use **access token only**:

```typescript
const res = await fetch(`${API_BASE}/api/v1/events/`, {
  headers: {
    Authorization: `Bearer ${accessToken}`,
    "Content-Type": "application/json",
  },
});
const body = await res.json();
if (body.error) throw new Error(body.error.message);
const events = body.data;
```

Use **`/api/v1/`** prefix for all protected routes.

### Telegram UI integration

| SDK API | Use in Druna |
|---------|--------------|
| `WebApp.MainButton` | "Save event", "Confirm time", "Send request" |
| `WebApp.BackButton` | Navigate back in multi-step forms |
| `WebApp.HapticFeedback` | Success/error feedback |
| `WebApp.showAlert()` | Error messages |
| `WebApp.close()` | Exit after action complete |
| `WebApp.openTelegramLink()` | Open friend's @username |

```typescript
WebApp.MainButton.setText("Create event");
WebApp.MainButton.onClick(() => submitEvent());
WebApp.MainButton.show();
```

### Screens for TMA MVP

| Screen | API used |
|--------|----------|
| Splash / auto-login | `POST /auth/telegram` |
| Profile | `GET /api/v1/users/me` |
| My events (list) | `GET /api/v1/events/` |
| Create / edit event | `POST` / `PATCH /api/v1/events/:id` |
| Day free time | `POST /api/v1/events/free-time` |
| Friends list | `GET /api/v1/friends/list` |
| Friend requests | `GET .../requests/incoming`, `.../outgoing` |
| Search & add friend | `GET .../search?username=`, `POST .../request` |
| Groups | `GET /api/v1/groups/list`, `POST .../create` |
| Group detail | `GET /api/v1/groups/:id` |
| Group scheduling | `POST .../confirm`, `POST .../free-time` |
| Group events | `GET` / `POST /api/v1/groups/:id/events`, `PATCH` / `DELETE .../:eventId` |

No username/password forms needed unless you support browser fallback.

### User identity in the app

After login, JWT contains `user_id` and `username`. Fetch full profile via `GET /api/v1/users/me`. Telegram account maps to:

| Field | Value |
|-------|-------|
| `telegram_id` | Stable Telegram user ID (DB) |
| `username` | Telegram `@username` or `tg_{id}` |
| `name` | Display name from Telegram profile |

For friend search, users can search by Telegram `@username` if the user has one; otherwise by `tg_{id}` pattern.

### CORS & networking

Mini App runs in Telegram WebView — requests are cross-origin from your app domain to API domain.

- Server default `CORS_ORIGINS=*` works for Bearer token auth
- Set explicit origins in production if needed:
  ```
  CORS_ORIGINS=https://your-tma-domain.com
  ```

### Local dev checklist

```bash
# Terminal 1 — backend
cp configs/config.yaml.example configs/config.yaml
cp .env.example .env
# set DB_PASSWORD, JWT_SECRET, BOT_TOKEN
go run cmd/main.go

# Terminal 2 — frontend (Vite example)
VITE_API_URL=http://localhost:8000 npm run dev

# Terminal 3 — HTTPS tunnel for Telegram
ngrok http 5173
# → set ngrok HTTPS URL in BotFather Menu Button
# → set VITE_API_URL to reachable backend (ngrok or LAN IP if testing)
```

**Browser-only dev** (without Telegram): use `POST /auth/sign-in` with test user as fallback when `WebApp.initData` is empty.

### Production checklist

- [ ] Mini App served over **HTTPS**
- [ ] `BOT_TOKEN` on server matches the bot opening the app
- [ ] `VITE_API_URL` points to production API (HTTPS)
- [ ] CORS configured if not using `*`
- [ ] Token refresh / re-auth via initData on 401
- [ ] Theme colors from `WebApp.themeParams`
- [ ] Test on iOS and Android Telegram clients (WebView differs slightly)

### Security notes

- Server validates `initData` HMAC — client must not skip `/auth/telegram`
- `initDataUnsafe` is for UI prefill only (name, photo), not for auth
- Access token in memory; avoid logging tokens
- Backend currently does **not** enforce `auth_date` TTL on initData — rely on fresh initData each session; do not cache initData long-term

### Optional companion bot (notifications)

A separate small bot process can send messages via Telegram Bot API (`sendMessage`) for:

- Incoming friend requests
- Group time confirmed
- Group event created (`group_event_created`)
- Event reminders

That bot is **not part of DrunaServer** today — Mini App handles UI; bot is for pushes only.

---

## Base URLs

| Environment | API base URL | Typical TMA frontend URL |
|-------------|--------------|--------------------------|
| Local dev | `http://localhost:8000` | `http://localhost:5173` (browser) |
| Docker API | `http://localhost:22000` | — |
| Production | `https://api.yourdomain.com` | `https://app.yourdomain.com` |

Configure in frontend:

```env
VITE_API_URL=http://localhost:8000
```

Swagger: `{API_BASE}/swagger/index.html`

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

**Parsing (fetch):**

```typescript
const body = await response.json();
if (body.error) throw new Error(body.error.message);
const payload = body.data;
```

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

**Token TTL:** access **12 hours**, refresh **7 days**.

### Telegram login (primary for Mini App)

`POST /auth/telegram`

```json
{ "initData": "<WebApp.initData raw string>" }
```

**Response `data`:**

```json
{
  "accessToken": "eyJ...",
  "refreshToken": "eyJ..."
}
```

Server env: `BOT_TOKEN` must match the bot linked to this Mini App.

### Sign up (web fallback)

`POST /auth/sign-up`

```json
{
  "name": "Alice",
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret123"
}
```

Password must be at least **8 characters**.

**Response `data`:** `{ "id": 1 }`

### Sign in (web fallback / browser dev)

`POST /auth/sign-in`

```json
{
  "username": "alice",
  "password": "secret123"
}
```

Legacy field `passwordHash` is also accepted.

### Renew tokens

`POST /auth/renew-token`

```json
{ "refreshToken": "eyJ..." }
```

Alternative: `Authorization: Bearer <refreshToken>` header.

**Response `data`:** new `accessToken` + `refreshToken`.

### Recommended client flow (all platforms)

```
1. Obtain tokens (Telegram: /auth/telegram, Web: /auth/sign-in)
2. API call with Authorization: Bearer <accessToken>
3. On 401 → POST /auth/renew-token → retry once
4. TMA: if renew fails → POST /auth/telegram again with fresh initData
5. Web: if renew fails → redirect to login
```

---

## Profile

Prefix: `/api/v1/users`. All routes require auth.

### Get current user

`GET /api/v1/users/me`

**Response `data`:**

```json
{
  "id": 1,
  "name": "Alice",
  "username": "alice",
  "email": "alice@example.com",
  "avatarURL": "https://example.com/avatar.png",
  "telegramID": 123456789
}
```

`telegramID` is omitted for password-only accounts.

### Update profile

`PATCH /api/v1/users/me`

```json
{
  "name": "Alice Updated",
  "avatarURL": "https://example.com/new.png"
}
```

Both fields are optional; omitted fields are unchanged.

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

Validation errors (400): end before start, overlapping event.

### Update event

`PATCH /api/v1/events/:id` — same body as create.

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
    { "start": "2026-06-17T08:00:00Z", "end": "2026-06-17T10:00:00Z" }
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
| GET | `/request-list` | — | all pending |
| GET | `/requests/incoming` | — | incoming pending |
| GET | `/requests/outgoing` | — | outgoing pending |
| POST | `/request` | `{ "username": "bob" }` | `{ "message": "friend request sent" }` |
| POST | `/accept` | `{ "username": "bob" }` | `{ "message": "friend request accepted" }` |
| POST | `/reject` | `{ "username": "bob" }` | `{ "message": "friend request rejected" }` |
| DELETE | `/` | `{ "username": "bob" }` | `{ "message": "friend deleted" }` |

**FriendInfo:** `{ "id": 2, "name": "Bob", "username": "bob" }`

Business rules (400): no self-request, no duplicate pending, no re-request after reject.

---

## Groups

Prefix: `/api/v1/groups`. All routes require auth.

| Method | Path | Body | Notes |
|--------|------|------|-------|
| POST | `/create` | `{ "name": "Weekend trip" }` | Owner = current user; response includes `groupId` |
| GET | `/list` | — | User's groups |
| GET | `/:id` | — | Details + members |
| POST | `/:id/members` | `{ "username": "bob" }` | Owner only; **member must be an accepted friend** |
| POST | `/:id/confirm` | `{ "confirmedTime": "2026-06-20T18:00:00Z" }` | Member confirms |
| POST | `/:id/free-time` | `{ "date": "2026-06-20" }` | Shared free slots |
| POST | `/:id/leave` | — | Not for owner |
| DELETE | `/:id` | — | Owner only |

**Group free-time** returns intersection of all members' free slots for the day. Busy time for each member now combines their **personal events** and the **group events** of every group they belong to, so a group event blocks the corresponding slot.

### Group events

Prefix: `/api/v1/groups/:id/events`. All routes require the caller to be a **member** of the group (non-members get **403**). A group event is a shared event scoped to the group; it does **not** appear in the personal `GET /api/v1/events/` list of its creator.

| Method | Path | Body | Notes |
|--------|------|------|-------|
| GET | `/:id/events` | — (query: `limit`, `offset`, `type`, `dateFrom`, `dateTo`) | List group events. `data` is the same envelope as personal events (`events`, `total`, `limit`, `offset`) |
| POST | `/:id/events` | `{ "title", "startTime", "endTime", "type" }` | Any member creates; response `{ "eventId": N }` |
| PATCH | `/:id/events/:eventId` | same body as create | Only the **creator** or the **group owner** |
| DELETE | `/:id/events/:eventId` | — | Only the **creator** or the **group owner** |

Event object (in list responses):

```json
{
  "eventID": 12,
  "userID": 3,
  "groupID": 7,
  "title": "Team sync",
  "startTime": "2026-06-17T14:00:00Z",
  "endTime": "2026-06-17T16:00:00Z",
  "type": "meeting"
}
```

Rules:

- `startTime`/`endTime`/`title` are required; `endTime` must be after `startTime` (**400**).
- A new/updated group event must not overlap another event **of the same group** (**400**).
- Update/delete by a member who is neither creator nor owner returns **403**; unknown event returns **404**.
- On create, the other group members receive a `group_event_created` notification in `notification_outbox`.

---

## Public utility routes

| Method | Path | Response |
|--------|------|----------|
| GET | `/ping/` | `{ "data": { "status": "ok", "db": "ok" } }` — on DB failure: HTTP **503**, `"status": "degraded"`, `"db": "error"` |
| GET | `/metrics` | Prometheus (not JSON envelope); disable with `METRICS_ENABLED=false` |
| GET | `/swagger/*` | Swagger UI — protect in production via reverse proxy |

Friend request, group confirm, and group event created events are enqueued in `notification_outbox` for a future companion Telegram bot.

---

## Shared code: Web ↔ TMA ↔ Mobile

Extract into shared packages (monorepo):

```
packages/
  api/       # fetch client, endpoints, 401 refresh interceptor
  types/     # ApiResponse, Event, Friend, Group, Tokens
  core/      # hooks: useEvents, useFriends, useGroups
  auth/      # AuthProvider interface + TelegramAuthProvider + WebAuthProvider
apps/
  web/       # Vite + React
  tma/       # same React, @twa-dev/sdk entry, Telegram theme
  mobile/    # Expo/RN, reuses api + types + core
```

| Layer | Web | TMA | Mobile |
|-------|-----|-----|--------|
| `packages/api` | ✅ | ✅ | ✅ |
| `packages/types` | ✅ | ✅ | ✅ |
| `packages/core` | ✅ | ✅ | ✅ |
| UI components | web | TMA-themed web | native |
| Login | `/auth/sign-in` | `/auth/telegram` | sign-in or OAuth later |

TMA and web can be **the same React codebase** with two entry points:

```typescript
// main-tma.tsx
if (WebApp.initData) await telegramAuthProvider.login();
else showBrowserFallback();

// main-web.tsx
await webAuthProvider.login(username, password);
```

---

## TypeScript client

```typescript
type ApiResponse<T> = {
  data: T | null;
  error: { message: string; code: number } | null;
};

type Tokens = { accessToken: string; refreshToken: string };

class DrunaClient {
  constructor(
    private baseUrl: string,
    private getAccessToken: () => string | null,
    private onTokensUpdated?: (t: Tokens) => void
  ) {}

  private async request<T>(path: string, init: RequestInit = {}): Promise<T> {
    const token = this.getAccessToken();
    const res = await fetch(`${this.baseUrl}${path}`, {
      ...init,
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...init.headers,
      },
    });
    const body: ApiResponse<T> = await res.json();
    if (body.error) throw new Error(body.error.message);
    return body.data as T;
  }

  authTelegram(initData: string) {
    return this.request<Tokens>("/auth/telegram", {
      method: "POST",
      body: JSON.stringify({ initData }),
    });
  }

  getEvents(params?: Record<string, string>) {
    const q = params ? "?" + new URLSearchParams(params) : "";
    return this.request(`/api/v1/events/${q}`);
  }

  createEvent(event: object) {
    return this.request("/api/v1/events/", {
      method: "POST",
      body: JSON.stringify(event),
    });
  }

  getFreeTime(date: string) {
    return this.request("/api/v1/events/free-time", {
      method: "POST",
      body: JSON.stringify({ date }),
    });
  }
}
```

---

## Testing

### Backend

```bash
cp configs/config.yaml.example configs/config.yaml
cp .env.example .env   # DB_PASSWORD, JWT_SECRET, BOT_TOKEN
go run cmd/main.go
```

### Smoke test (web auth path)

```bash
make smoke
# BASE_URL=http://localhost:8000 make smoke
```

### TMA manual test

1. Start backend + frontend with HTTPS tunnel
2. Set Menu Button URL in BotFather
3. Open bot → tap Menu → Mini App loads
4. Verify auto-login and `GET /api/v1/events/`

### Browser fallback (no Telegram)

Open app in Chrome → use sign-in form hitting `POST /auth/sign-in` when `WebApp.initData` is empty.

---

## Changelog notes

- All responses: `{ data, error }` envelope
- Use `/api/v1/` prefix
- TMA auth: `POST /auth/telegram` with raw `WebApp.initData`
- Access vs refresh tokens — only access on API routes
- Refresh rotation — save new tokens after every renew
