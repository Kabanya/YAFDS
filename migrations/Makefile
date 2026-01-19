# Makefile for managing database migrations across all services
# Assumes each service has its own Makefile with migrate-up, migrate-down, etc.

.PHONY: up down migrate-status create-migration clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  up                 - Apply all migrations for all services"
	@echo "  down               - Rollback all migrations for all services"
	@echo "  migrate-status     - Show migration status for all services"
	@echo "  clean.             - Truncate all tables in the database (clear data)"
	@echo "  create-migration   - Create a new migration file (requires NAME and SERVICE)"
	@echo "  help               - Show this help"

# Apply migrations for all services
up: migrate-up-customer migrate-up-courier migrate-up-restaurant
	@echo "All migrations applied successfully."

# Rollback migrations for all services
down: migrate-down-customer migrate-down-courier migrate-down-restaurant
	@echo "All migrations rolled back."

# Show status for all services
migrate-status: migrate-status-customer migrate-status-courier migrate-status-restaurant
	@echo "Migration status checked for all services."

# Clean database (truncate all tables)
clean:
	@echo "Copying migrations/clear_db.sql into container and executing against yafds_db..."
	@docker cp clear_db.sql yafds-postgres:/tmp/clear_db.sql || docker cp ./clear_db.sql yafds-postgres:/tmp/clear_db.sql
	@docker exec yafds-postgres psql -U postgres -h localhost -p 5644 -d yafds_db -f /tmp/clear_db.sql
	@docker exec yafds-postgres rm -f /tmp/clear_db.sql
	@echo "Database cleared."

# Individual service targets
migrate-up-customer:
	@echo "Applying customer migrations..."
	cd ../customer && make migrate-up

migrate-down-customer:
	@echo "Rolling back customer migrations..."
	cd ../customer && make migrate-down

migrate-status-customer:
	@echo "Customer migration status:"
	cd ../customer && make migrate-status

migrate-up-courier:
	@echo "Applying courier migrations..."
	cd ../courier && make migrate-up

migrate-down-courier:
	@echo "Rolling back courier migrations..."
	cd ../courier && make migrate-down

migrate-status-courier:
	@echo "Courier migration status:"
	cd ../courier && make migrate-status

migrate-up-restaurant:
	@echo "Applying restaurant migrations..."
	cd ../restaurant && make migrate-up

migrate-down-restaurant:
	@echo "Rolling back restaurant migrations..."
	cd ../restaurant && make migrate-down

migrate-status-restaurant:
	@echo "Restaurant migration status:"
	cd ../restaurant && make migrate-status

# Create a new migration file
# Usage: make create-migration NAME="create_new_table" SERVICE=customer
create-migration:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make create-migration NAME='create_new_table' SERVICE=customer"; \
		exit 1; \
	fi
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE is required. Usage: make create-migration NAME='create_new_table' SERVICE=customer"; \
		exit 1; \
	fi
	@echo "Creating migration for $(SERVICE): $(NAME)"
	cd ../$(SERVICE) && make create-migration NAME="$(NAME)"