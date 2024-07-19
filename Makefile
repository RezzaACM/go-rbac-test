include .env

# DEFINE VARIABLES
MIGRATION_DIR=db/migrations

all:help

# create-migration:
# 	migrate create -ext=sql -dir=db/migrations -seq $(CMD)

migrate-up:
	migrate -path=db/migrations -database ${POSTGRESQL_URL} -verbose up

migrate-down:
	migrate -path=db/migrations -database ${POSTGRESQL_URL} -verbose down

migration:
ifndef MIGRATION_NAME
	@echo "Error: MIGRATION_NAME is required"
	@exit 1
endif
	@echo "Creating migration: $(MIGRATION_NAME)"
	migrate create -ext=sql -dir=$(MIGRATION_DIR) -seq $(MIGRATION_NAME)


# Help target
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Clean up the binary"
	@echo "  migration     - Create a new migration file"
	@echo "	 	Example: make migrate MIGRATION_NAME='create_users_table'"
	@echo "  migrate-up    - Run database migrations up"
	@echo "  migrate-down  - Run database migrations down"
	@echo "  deps          - Install required dependencies"
	@echo "  test          - Run tests"
	@echo "  run           - Build and run the binary"
	@echo "  rerun         - Rebuild and run the binary"
	@echo "  help          - Display this help message"

# Default target is help
.DEFAULT_GOAL := help