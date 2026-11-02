# CodeQuest dev tasks. `make test` is the bar for "done".
SHELL := /bin/bash
COMPOSE := docker compose
DATABASE_URL ?= postgres://codequest:codequest@localhost:5432/codequest?sslmode=disable
SQLC_IMAGE := sqlc/sqlc:1.30.0
LINT_IMAGE := golangci/golangci-lint:v2.12-alpine

.PHONY: up down db migrate content validate test test-api test-web lint tokens tokens-check sqlc api-types

up: ## Start Postgres, Redis, API and web (applies migrations and builds content/)
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

db: ## Start only Postgres and wait until healthy
	$(COMPOSE) up -d --wait postgres

migrate: db
	cd api && DATABASE_URL='$(DATABASE_URL)' go run ./cmd/cq migrate up

content: migrate ## Compile content/ into Postgres
	cd api && DATABASE_URL='$(DATABASE_URL)' go run ./cmd/cq build ../content

validate:
	cd api && go run ./cmd/cq validate ../content

test: test-api lint test-web ## Everything CI runs

test-api: db
	cd api && go vet ./... && \
	  CQ_REQUIRE_DB=1 TEST_DATABASE_URL='$(DATABASE_URL)' go test -count=1 ./...

lint:
	docker run --rm -v "$(CURDIR)/api":/src -w /src $(LINT_IMAGE) golangci-lint run ./...

test-web:
	cd web && npm run typecheck && npm test && npm run tokens:check

tokens: ## Regenerate web/src/styles/tokens.css and web/tailwind.config.ts from docs/design/tokens.json
	cd web && npm run tokens

tokens-check:
	cd web && npm run tokens:check

sqlc: ## Regenerate api/internal/store from db/queries
	docker run --rm -v "$(CURDIR)/api":/src -w /src $(SQLC_IMAGE) generate

api-types: ## Regenerate web/src/api/schema.d.ts from api/openapi.yaml
	cd web && npm run api-types
