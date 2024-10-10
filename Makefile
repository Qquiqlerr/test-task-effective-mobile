include .env
export $(shell sed 's/=.*//' .env)

MIGRATIONS_DIR := ./migrations
BIN_DIR := ./bin
CMD_DIR := ./cmd
DOCS_DIR := ./docs

DB_CONN := "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable"

.PHONY: all init build run migrate_up migrate_down clean doc_gen

all: init migrate_up build run

init:
	@echo "Initializing project..."
	go mod tidy
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/swaggo/swag/cmd/swag@latest

build:
	@echo "Building the application..."
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/main $(CMD_DIR)/main.go

run:
	@echo "Running the application..."
	$(BIN_DIR)/main

migrate_up:
	@echo "Applying migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_CONN) up

migrate_down:
	@echo "Rolling back the last migration..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_CONN) down

doc_gen:
	@echo "Generating Swagger documentation..."
	swag init -g $(CMD_DIR)/main.go --parseDependency --output $(DOCS_DIR)

clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
	rm -rf $(DOCS_DIR)