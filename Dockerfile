# syntax=docker/dockerfile:1

### Этап 1: сборка
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Если main.go в корне:
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/

### Этап 2: рантайм
FROM alpine:latest
COPY --from=builder /app/bot /bot
CMD ["/bot"]