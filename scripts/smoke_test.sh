#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8000}"
USERNAME="smoke_$(date +%s)"
PASSWORD="smoke-pass-123"

extract_data() {
  python3 - "$1" <<'PY'
import json, sys
payload = json.load(sys.stdin)
print(json.dumps(payload.get("data", payload)))
PY
}

echo "==> Sign up"
SIGNUP=$(curl -sf -X POST "$BASE_URL/auth/sign-up" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Smoke Test\",\"username\":\"$USERNAME\",\"email\":\"$USERNAME@test.local\",\"password\":\"$PASSWORD\"}")
extract_data <<<"$SIGNUP" >/dev/null

echo "==> Sign in"
TOKENS=$(curl -sf -X POST "$BASE_URL/auth/sign-in" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
TOKENS_DATA=$(extract_data <<<"$TOKENS")
ACCESS=$(echo "$TOKENS_DATA" | python3 -c "import json,sys; print(json.load(sys.stdin)['accessToken'])")

echo "==> Create event"
curl -sf -X POST "$BASE_URL/api/v1/events/" \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke event","startTime":"2026-06-17T10:00:00Z","endTime":"2026-06-17T11:00:00Z","type":"work"}' >/dev/null

echo "==> List events"
curl -sf -X GET "$BASE_URL/api/v1/events/" -H "Authorization: Bearer $ACCESS" >/dev/null

echo "==> Free time"
curl -sf -X POST "$BASE_URL/api/v1/events/free-time" \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{"date":"2026-06-17"}' >/dev/null

echo "==> Create group"
curl -sf -X POST "$BASE_URL/api/v1/groups/create" \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{"name":"Smoke group"}' >/dev/null

echo "Smoke test passed"
