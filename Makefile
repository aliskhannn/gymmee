ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DB_DRIVER ?= sqlite3
DB_STRING ?= ./gymlog.db
MIGRATION_DIR ?= ./migrations

.PHONY: help build run test clean migrate-up migrate-down migrate-status migrate-create

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build:
	@echo "Compilation..."
	go build -o bin/gymmee cmd/bot/main.go
	@echo "Done! The binary is located in bin/gymmee"

run:
	@echo "Launching the bot..."
	go run cmd/bot/main.go

test:
	go test -v -race ./...

clean:
	rm -rf bin/

migrate-up:
	goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DB_STRING)" up

migrate-down:
	goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DB_STRING)" down

migrate-status:
	goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DB_STRING)" status

migrate-create:
	@if [ -z "$(name)" ]; then echo "Error: Please provide a migration name, for example: make migrate-create name=init"; exit 1; fi
	goose -dir $(MIGRATION_DIR) create $(name) sql