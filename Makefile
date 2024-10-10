include .env
export $(shell sed 's/=.*//' .env)
start:
	goose postgres -dir ./migrations "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up
	export GOBIN=$(pwd)/bin
	go build -o ./bin/main cmd/main.go
	./bin/main

init:
	go mod tidy
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate_down:
	goose postgres -dir ./migrations "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" down
doc_gen:
	swag init -g cmd/main.go --parseDependency --output ./docs

