#!/bin/sh
set -e

echo "==> gofmt"
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
  echo "These files need gofmt:"
  echo "$unformatted"
  exit 1
fi

echo "==> go vet"
JWT_SECRET=precommit-test-secret go vet ./...

echo "==> go test"
JWT_SECRET=precommit-test-secret go test ./pkg/... -count=1

echo "Pre-commit checks passed"
