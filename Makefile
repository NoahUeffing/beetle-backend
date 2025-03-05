SHELL := /bin/bash
GO := go
PACKAGE_NAME := beetle
GO_PATH := $(shell $(GO) env GOPATH)
EXTENDED_PATH := PATH=$$PATH:$(GO_PATH)/bin

# Development configuration
DEVELOPMENT_CONFIG := ./configs/dev.yaml

# Database configuration
DB_USER := postgres
DB_PASSWORD := postgres
DB_NAME := beetle
DB_HOST := localhost
DB_PORT := 5432
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Goose migration settings
GOOSE_DIR := ./assets/db/migrations
GOOSE_DIR_SEED := ./assets/db/seeds
GOOSE_DRIVER := postgres
GOOSE_DBSTRING := $(DB_URL)
GOOSE := $(EXTENDED_PATH) goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)"
GOOSE_SEED := $(EXTENDED_PATH) goose -dir $(GOOSE_DIR_SEED) -no-versioning $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)"

.PHONY: all build run test clean db-start db-stop db-restart db-logs migrate-create migrate-up migrate-down migrate-status coverage coverage-report watch install install-dev-tools docker-build docker-up docker-up-db

all: build

build:
	$(GO) build -o ./bin/$(PACKAGE_NAME) ./cmd/api

run: build
	./bin/$(PACKAGE_NAME)

test:
	$(GO) test ./...

clean:
	rm -rf bin/*
	rm -f coveragereport.out

# Coverage commands
coverage-report:
	@$(GO) test ./... -coverprofile=coveragereport.out

coverage: coverage-report
	@go tool cover -html=coveragereport.out

# Development commands
install-dev-tools:
	@go install github.com/air-verse/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest

watch: install-dev-tools
	@BEETLE_CONFIG_FILE=$(DEVELOPMENT_CONFIG) air

install:
	$(GO) mod download

# Docker commands
docker-build:
	docker build -t $(PACKAGE_NAME) .

docker-up:
	docker-compose --profile all up

docker-up-db:
	docker-compose --profile db up

# Database commands
db-start:
	docker-compose --profile db up -d

db-stop:
	docker-compose --profile db down

db-restart: db-stop db-start

db-logs:
	docker-compose --profile db logs -f

# Migration commands
goose-create:
	@read -p "Enter migration name: " name; \
	$(GOOSE) create $$name sql

migrate-up: goose-up
goose-up:
	$(GOOSE) up

migrate-down: goose-down
goose-down:
	$(GOOSE) down

migrate-status: goose-status
goose-status:
	$(GOOSE) status

goose-up-one:
	$(GOOSE) up-by-one

goose-reset:
	$(GOOSE) reset

goose-clean: goose-reset goose-up

# Seed commands
goose-create-seed:
	@read -p "Enter seed name: " name; \
	$(GOOSE_SEED) create $$name sql

goose-seed-up:
	$(GOOSE_SEED) up

goose-seed-down:
	$(GOOSE_SEED) reset

# Create directories if they don't exist
init:
	mkdir -p $(GOOSE_DIR)
	mkdir -p $(GOOSE_DIR_SEED)
	mkdir -p ./bin

# Swagger commands
swagger:
	$(EXTENDED_PATH) swag init -g ./internal/server/server.go -o ./api/swaggergenerated --parseInternal

# Linting
lint: fmt vet staticcheck

fmt:
	@$(GO) fmt ./...

vet:
	@$(GO) vet ./...

staticcheck:
	@staticcheck ./... 