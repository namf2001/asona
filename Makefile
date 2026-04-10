# Simple Makefile for a Go project

-include .env
export

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go &
	@cd frontend && pnpm install --frozen-lockfile
	@cd frontend && pnpm dev

# Start all containers (app + postgres + redis + frontend)
docker-run:
	@docker compose up --build

# Start only infrastructure services (postgres + redis)
docker-db-run:
	@docker compose up -d postgres redis

# Stop only infrastructure services (postgres + redis)
docker-db-down:
	@docker compose stop postgres redis

# Shutdown all containers
docker-down:
	@docker compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integration Tests for the database layer
itest:
	@echo "Running integration tests..."
	@go test ./internal/service/database -v

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	swag init -g cmd/api/main.go -o docs/swagger || go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs/swagger

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Clear binary and build artifacts (main binary, air tmp dir, build logs)
clear:
	@echo "Clearing build artifacts..."
	@rm -f main
	@rm -rf tmp/
	@rm -f build-errors.log
	@echo "Done."

# Live Reload
watch: be-dev

# Run backend with Air hot-reload (uses .air.toml config)
be-dev:
	@echo "Starting backend with Air hot-reload..."
	@if command -v air > /dev/null; then \
		air -c .air.toml; \
	else \
		echo "Air is not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
		air -c .air.toml; \
	fi

# ── Postgres Migrations ───────────────────────────────────────────────────────
MIGRATIONS_DIR  = migrations
DB_URL          = "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

# Run all migrations (UP)
migrate-up:
	@echo "Running Postgres migrations (UP)..."
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up; \
	else \
		read -p "golang-migrate (migrate) is not installed. Install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
			migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up; \
		else \
			echo "Manual psql fallback..."; \
			for file in $(MIGRATIONS_DIR)/*.up.sql; do \
				echo "Applying $$file..."; \
				PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f $$file; \
			done; \
		fi; \
	fi

# Rollback migrations (DOWN)
migrate-down:
	@echo "Running Postgres migrations (DOWN)..."
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1; \
	else \
		echo "Manual psql fallback (DOWN)..."; \
		for file in $$(ls -r $(MIGRATIONS_DIR)/*.down.sql); do \
			echo "Rolling back $$file..."; \
			PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f $$file; \
			break; \
		done; \
	fi

# Show migration status
migrate-status:
	@echo "Migration files in $(MIGRATIONS_DIR)/:"
	@ls $(MIGRATIONS_DIR)/*.sql | sort -V

.PHONY: all build run test clean clear watch be-dev docker-run docker-down docker-db-run docker-db-down itest migrate-up migrate-down migrate-status swagger
