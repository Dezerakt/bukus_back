include .env
export

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: run
run: ## Run the application
	go run ./cmd/app

.PHONY: run-migrate
run-migrate: ## Run database migrations
	go run ./cmd/app migrate

.PHONY: compose-up
compose-up: ## Start docker compose services (Postgres)
	docker compose up -d

.PHONY: compose-down
compose-down: ## Stop docker compose services
	docker compose down

# ==================================================================================== #
# QUALITY
# ==================================================================================== #

.PHONY: lint
lint: ## Run linter
	golangci-lint run ./...

.PHONY: test
test: ## Run tests
	go test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# ==================================================================================== #
# BUILD
# ==================================================================================== #

.PHONY: build
build: ## Build the application
	go build -o ./bin/app ./cmd/app

# ==================================================================================== #
# CODEGEN
# ==================================================================================== #

.PHONY: generate
generate: ## Run go generate
	go generate ./...

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=<migration_name>)
	migrate create -ext sql -dir migrations -seq $(name)
