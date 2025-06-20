APP_NAME=drunaServer
DOCKER_COMPOSE=docker-compose -f docker-compose.yml

.PHONY: build up down migrate

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
		-path=/migrations -database "postgres://postgres:qwerty@db:22000/druna_db?sslmode=disable" up

migrate-down:
	docker run --network=drunaserver_druna-net --rm \
		-v ${PWD}/migrations:/migrations \
		migrate/migrate \
		-path=/migrations -database "postgres://postgres:postgres@db:5432/druna_db?sslmode=disable" down
