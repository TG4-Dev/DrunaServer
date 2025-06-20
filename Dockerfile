#Dockerfile

FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o drunaServer ./cmd/main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/drunaServer .
COPY configs ./configs
COPY .env .

CMD ["./drunaServer"]