#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8000}"
USERNAME="smoke_$(date +%s)"
PASSWORD="smoke-pass-123"

echo "==> Sign up"
curl -sf -X POST "$BASE_URL/auth/sign-up" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Smoke Test\",\"username\":\"$USERNAME\",\"email\":\"$USERNAME@test.local\",\"passwordHash\":\"$PASSWORD\"}"

echo
echo "==> Sign in"
TOKENS=$(curl -sf -X POST "$BASE_URL/auth/sign-in" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"passwordHash\":\"$PASSWORD\"}")
ACCESS=$(echo "$TOKENS" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)

echo "==> Create event"
EVENT=$(curl -sf -X POST "$BASE_URL/api/events/" \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke event","startTime":"2026-06-17T10:00:00Z","endTime":"2026-06-17T11:00:00Z","type":"work"}')

echo "==> List events"
curl -sf -X GET "$BASE_URL/api/events/" -H "Authorization: Bearer $ACCESS"

echo
echo "==> Free time"
curl -sf -X POST "$BASE_URL/api/events/free-time" \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{"date":"2026-06-17"}'

echo
echo "==> Create group"
curl -sf -X POST "$BASE_URL/api/groups/create" \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{"name":"Smoke group"}'

echo
echo "Smoke test passed"
