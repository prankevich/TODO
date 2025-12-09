POSTGRES_USER=postgres
POSTGRES_PASSWORD=1234
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=postgres

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable


migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-reset:
	migrate -path ./migrations -database "$(DB_URL)" down -all

migrate-create:
	migrate create -ext sql -dir ./migrations -seq -digits 4 $(name)

format:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run ./...