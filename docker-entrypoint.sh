#!/bin/sh
set -e

if [ -n "$DATABASE_URL" ]; then
  echo "Running database migrations..."
  migrate -path /migrations -database "$DATABASE_URL" up
fi

exec ./drunaServer
