# CodeQuest dev tasks. `make test` is the bar for "done".
SHELL := /bin/bash
COMPOSE := docker compose
DATABASE_URL ?= postgres://codequest:codequest@localhost:5432/codequest?sslmode=disable
SQLC_IMAGE := sqlc/sqlc:1.30.0
LINT_IMAGE := golangci/golangci-lint:v2.12-alpine

REDIS_TEST_URL ?= redis://localhost:6379/15
RUNNER_IMAGE := codequest/python:3.12
PW_IMAGE := mcr.microsoft.com/playwright:v1.64.1-noble

.PHONY: up down db migrate content validate content-test runner-image e2e lighthouse test test-api test-web test-sandbox lint tokens tokens-check sqlc api-types cli

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

e2e: ## Browser tests (sandbox isolation, timeouts, Run speed, content CI in the browser, screenshots)
	$(COMPOSE) up -d --build --wait
	cd api && DATABASE_URL='$(DATABASE_URL)' go run ./cmd/cq build ../content
	docker run --rm --network host --ipc host -v "$(CURDIR)":/repo -w /repo/web $(PW_IMAGE) npx playwright test

# Lighthouse on the landing page (production build, served by vite preview and
# proxied to the Compose API). The host's node_modules hold macOS binaries, so
# the build runs on a Linux copy inside the Playwright image.
lighthouse: ## Landing page Lighthouse audit; fails below 90 in any category
	$(COMPOSE) up -d --build --wait api
	docker run --rm --network host --ipc host -v "$(CURDIR)":/repo:ro -v "$(CURDIR)/web/e2e/.results":/out -v codequest_npm_cache:/root/.npm $(PW_IMAGE) sh -ec '\
	  mkdir /tmp/web && cd /repo/web && tar --exclude=node_modules --exclude=dist --exclude=e2e -cf - . | tar -xf - -C /tmp/web; \
	  cd /tmp/web && npm ci --no-audit --no-fund --loglevel=error && npm run build >/dev/null; \
	  (npx vite preview --port 4173 --strictPort >/tmp/preview.log 2>&1 &); \
	  for i in $$(seq 1 30); do node -e "fetch(\"http://localhost:4173/\").then(()=>process.exit(0),()=>process.exit(1))" && break; sleep 1; done; \
	  node scripts/lighthouse.mjs http://localhost:4173/ /out'

test: test-api content-test lint test-web test-sandbox e2e lighthouse ## Everything CI runs

test-api: db runner-image
	cd api && go vet ./... && \
	  CQ_REQUIRE_DB=1 CQ_REQUIRE_DOCKER=1 TEST_DATABASE_URL='$(DATABASE_URL)' TEST_REDIS_URL='$(REDIS_TEST_URL)' \
	  go test -count=1 ./...

lint:
	docker run --rm -v "$(CURDIR)/api":/src -w /src $(LINT_IMAGE) golangci-lint run ./...

test-web:
	cd web && npm run typecheck && npm test && npm run tokens:check

test-sandbox:
	cd sandbox && npm run assets && npm run typecheck && npm test

tokens: ## Regenerate web/src/styles/tokens.css and web/tailwind.config.ts from docs/design/tokens.json
	cd web && npm run tokens

tokens-check:
	cd web && npm run tokens:check

cli: ## Build the learner CLI into bin/codequest (set API_URL / APP_URL for a release build)
	cd api && go build -trimpath -ldflags "-X main.version=$$(git describe --always --dirty 2>/dev/null || echo dev)$(if $(API_URL), -X main.defaultAPIURL=$(API_URL))$(if $(APP_URL), -X main.defaultAppURL=$(APP_URL))" -o ../bin/codequest ./cmd/codequest

sqlc: ## Regenerate api/internal/store from db/queries
	docker run --rm -v "$(CURDIR)/api":/src -w /src $(SQLC_IMAGE) generate

api-types: ## Regenerate web/src/api/schema.d.ts from api/openapi.yaml
	cd web && npm run api-types
