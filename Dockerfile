#Dockerfile

FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o drunaServer ./cmd/main.go


FROM alpine:latest

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz \
    | tar xz -C /usr/local/bin migrate

WORKDIR /root/
COPY --from=builder /app/drunaServer .
COPY configs/config.docker.yaml ./configs/config.yaml
COPY migrations /migrations
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
