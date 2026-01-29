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
FROM alpine:3.20

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /app/app /app/app

EXPOSE 8284

USER appuser

ENTRYPOINT ["/app/app"]