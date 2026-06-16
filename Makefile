APP_NAME=drunaServer
DOCKER_COMPOSE=docker-compose -f docker-compose.yml
MIGRATE_DB=postgres://postgres:postgres@db:5432/druna_db?sslmode=disable

.PHONY: build up down migrate-up migrate-down test lint smoke hook-install dev-up

hook-install:
	git config core.hooksPath .githooks

dev-up:
	docker compose -f docker-compose.dev.yml up

build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

migrate-up:
	docker run --network=drunaserver_druna-net --rm \
		-v ${PWD}/migrations:/migrations \
		migrate/migrate \
		-path=/migrations -database "$(MIGRATE_DB)" up

migrate-down:
	docker run --network=drunaserver_druna-net --rm \
		-v ${PWD}/migrations:/migrations \
		migrate/migrate \
		-path=/migrations -database "$(MIGRATE_DB)" down

test:
	JWT_SECRET=test-secret-key go test ./... -count=1

lint:
	golangci-lint run ./...

smoke:
	bash scripts/smoke_test.sh
