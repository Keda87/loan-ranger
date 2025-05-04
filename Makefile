include .env
export
export GO111MODULE        ?= on

run:
	go run cmd/api/main.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -installsuffix cgo -o app.backend cmd/rest/main.go

migration-status:
	goose -dir migration postgres "user=${DB_USER} password=${DB_PASS} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" status

.PHONY: migration
migration:
	goose -dir migration postgres "user=${DB_USER} password=${DB_PASS} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" create $(name) sql

migrate:
	goose -dir migration postgres "user=${DB_USER} password=${DB_PASS} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" up

migrate-back:
	goose -dir migration postgres "user=${DB_USER} password=${DB_PASS} host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} sslmode=disable" down

test:
	ENV=testing go test ./... -covermode=count -coverprofile=coverage.out ; go tool cover -func=coverage.out | grep total

test-cover:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: docs
docs:
	swag init -g cmd/rest/main.go