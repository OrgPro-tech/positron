# Makefile for database migrations

# Variables
DB_URL ?= postgres://username:password@localhost:5432/database_name?sslmode=disable
MIGRATIONS_DIR = ./migrations

# Check if golang-migrate is installed
MIGRATE := $(shell command -v migrate 2> /dev/null)

# Phony targets
.PHONY: install-migrate migrate-create migrate-up migrate-down migrate-force

# Install golang-migrate if not already installed
install-migrate:
ifndef MIGRATE
	@echo "Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
else
	@echo "golang-migrate is already installed."
endif

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $${name}

# Run all up migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Roll back all migrations
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

# Roll back one migration
migrate-down-one:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Force set version
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $$version

# Show current migration version
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Help target
help:
	@echo "Available targets:"
	@echo "  install-migrate   - Install golang-migrate tool"
	@echo "  migrate-create    - Create a new migration file"
	@echo "  migrate-up        - Run all up migrations"
	@echo "  migrate-down      - Roll back all migrations"
	@echo "  migrate-down-one  - Roll back one migration"
	@echo "  migrate-force     - Force set a specific version"
	@echo "  migrate-version   - Show current migration version"
	@echo ""
	@echo "Usage:"
	@echo "  make [target] DB_URL='your_database_url'"
	@echo "  Default DB_URL is 'postgres://username:password@localhost:5432/database_name?sslmode=disable'"