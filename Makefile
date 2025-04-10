SHELL := /bin/bash
GO := go
PACKAGE_NAME := beetle
GO_PATH := $(shell $(GO) env GOPATH)
MAIN_GO := ./cmd/main.go
IGNORED_DIRS := /swaggergenerated # regex
PACKAGE_DIRS := $(shell cd $(CURDIR) && $(GO) list ./... | grep -v $(IGNORED_DIRS) | grep -v "^$$")
BUILDFLAGS := ''
TESTS=$(PACKAGE_DIRS)

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
GOOSE := $(GO_PATH)/bin/goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING)
GOOSE_SEED := $(GO_PATH)/bin/goose -dir $(GOOSE_DIR_SEED) -no-versioning $(GOOSE_DRIVER) $(GOOSE_DBSTRING)

.PHONY: all build run test clean migrate-create migrate-up migrate-down migrate-status coverage coverage-report watch install docker-build docker-up docker-up-db

all: run

check: fmt lint compile test swag

# Also used by air
compile:
	$(GO) build -ldflags $(BUILDFLAGS) -o ./bin/$(PACKAGE_NAME) $(MAIN_GO)

build: swag compile

run: build
	@./bin/$(PACKAGE_NAME)

run-dev: build
	@BEETLE_CONFIG_FILE=$(DEVELOPMENT_CONFIG) ./bin/$(PACKAGE_NAME)

test:
	@$(GO) test $(TESTS)

clear-test-cache:
	@$(GO) clean -testcache

test-all:	clear-test-cache test

coverage-report:
	@$(GO) test $(TESTS) -coverprofile=coveragereport.out

coverage: coverage-report
	@go tool cover -html=coveragereport.out

watch:
	@BEETLE_CONFIG_FILE=$(DEVELOPMENT_CONFIG) $(GO_PATH)/bin/air

install:
	$(GO) mod download

swag:
	@$(GO_PATH)/bin/swag init -g ./internal/server/server.go -o ./swaggergenerated --parseInternal --generatedTime

fmt:
	@gofmt -s -w .

clean:
	rm -rf bin/*
	rm -rf swaggergenerated

staticcheck:
	@$(GO_PATH)/bin/staticcheck $(PACKAGE_DIRS)

vet:
	@$(GO) vet $(PACKAGE_DIRS)

fmt-check:
	@gofmt -l -s -e .

lint: staticcheck vet fmt-check

# Database commands
db-start:
	docker-compose --profile db up -d

db-stop:
	docker-compose --profile db down

db-restart: db-stop db-start

db-logs:
	docker-compose --profile db logs -f

# Migration commands
migrate-create: goose-create
migrate-up: goose-up
migrate-down: goose-down
migrate-status: goose-status

goose-create:
	@$(GOOSE) create temporary_title sql

goose-create-seed:
	@read -p "Enter seed name: " name; \
	$(GOOSE_SEED) create $$name sql

goose-create-go:
	@$(GOOSE) create temporary_title go

goose-up:
	@$(GOOSE) up

goose-seed-up:
	@$(GOOSE_SEED) up

goose-up-one:
	@$(GOOSE) up-by-one

goose-down:
	@$(GOOSE) down

goose-seed-down:
	@$(GOOSE_SEED) reset

goose-status:
	@$(GOOSE) status

goose-reset:
	@$(GOOSE) reset

goose-clean: goose-reset goose-up

regoose:
	@$(GOOSE) down
	@$(GOOSE) up-by-one

migrate: goose-up

# Docker commands
docker-build:
	docker build -t $(PACKAGE_NAME) .

docker-up:
	docker-compose --profile all up

docker-up-db:
	docker-compose --profile db up

# Create directories if they don't exist
init:
	mkdir -p $(GOOSE_DIR)
	mkdir -p $(GOOSE_DIR_SEED)
	mkdir -p ./bin
