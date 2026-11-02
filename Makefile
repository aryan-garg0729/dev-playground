# CodeQuest dev tasks. `make test` is the bar for "done".
SHELL := /bin/bash
COMPOSE := docker compose
DATABASE_URL ?= postgres://codequest:codequest@localhost:5432/codequest?sslmode=disable
SQLC_IMAGE := sqlc/sqlc:1.30.0
LINT_IMAGE := golangci/golangci-lint:v2.12-alpine

REDIS_TEST_URL ?= redis://localhost:6379/15
RUNNER_IMAGE := codequest/python:3.12

.PHONY: up down db migrate content validate content-test runner-image test test-api test-web lint tokens tokens-check sqlc api-types

up: ## Start Postgres, Redis, API, worker and web (applies migrations and builds content/)
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

db: ## Start Postgres and Redis and wait until healthy
	$(COMPOSE) up -d --wait postgres redis

migrate: db
	cd api && DATABASE_URL='$(DATABASE_URL)' go run ./cmd/cq migrate up

content: migrate ## Compile content/ into Postgres
	cd api && DATABASE_URL='$(DATABASE_URL)' go run ./cmd/cq build ../content

validate:
	cd api && go run ./cmd/cq validate ../content

runner-image: ## Build the codequest/python grading image
	docker build -q -t $(RUNNER_IMAGE) runners/python

content-test: runner-image ## Content CI: every solution passes, every wrong/ file fails
	cd api && go run ./cmd/cq test ../content

test: test-api content-test lint test-web ## Everything CI runs

test-api: db runner-image
	cd api && go vet ./... && \
	  CQ_REQUIRE_DB=1 CQ_REQUIRE_DOCKER=1 TEST_DATABASE_URL='$(DATABASE_URL)' TEST_REDIS_URL='$(REDIS_TEST_URL)' \
	  go test -count=1 ./...

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
