include ../.env

# Variables

DB_CONNECTION_STRING = $(DB_CONNECTION_BASE) dbname=$(1)

MIGRATIONS_DIR := ./cmd/database/migrations

CUSTOMER_DB := customer_db

# other
launch-docker:
	@if ! docker info > /dev/null 2>&1; then \
		echo "Docker is not running. Starting Docker Desktop..."; \
		open -a Docker; \
		sleep 5; \
		until docker info > /dev/null 2>&1; do \
			echo "Waiting for Docker to start..."; \
			sleep 2; \
		done; \
		echo "Docker started successfully!"; \
	else \
		echo "Docker is already running"; \
	fi

# db
launch-db:
	@if docker ps -a --format '{{.Names}}' | grep -q '^$(DB_CONTAINER_NAME)$$'; then \
		echo "Container $(DB_CONTAINER_NAME) exists. Starting..."; \
		docker start $(DB_CONTAINER_NAME); \
	else \
		echo "Container $(DB_CONTAINER_NAME) does not exist. Creating new one..."; \
		docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:13.3; \
	fi
	@sleep 1
	$(MAKE) create-customer-db

clear-db:
	docker stop $(DB_CONTAINER_NAME) || true
	docker rm $(DB_CONTAINER_NAME) || true

# migrations
migrate-customer-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(MIGRATIONS_DIR) up

migrate-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(MIGRATIONS_DIR) up

migrate-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(MIGRATIONS_DIR) down

# customer
build-customer:
	docker rmi customer-app || true
	docker build --no-cache --build-arg BUILD_TIME=$(shell date +%s) -t customer-app -f Dockerfile .

launch-customer:
	@if docker ps -a --format '{{.Names}}' | grep -q '^$(CUSTOMER_CONTAINER_NAME)$$'; then \
		echo "Customer container exists. Starting..."; \
		docker start $(CUSTOMER_CONTAINER_NAME); \
	else \
		echo "Customer container does not exist. Creating new one..."; \
		docker run --name $(CUSTOMER_CONTAINER_NAME) -p $(CUSTOMER_PORT):$(CUSTOMER_PORT) -d customer-app; \
	fi

clear-customer:
	docker stop $(CUSTOMER_CONTAINER_NAME) || true
	docker rm $(CUSTOMER_CONTAINER_NAME) || true


#db's
create-customer-db:
	docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(CUSTOMER_DB);" 2>/dev/null || true
	$(MAKE) migrate-customer-up

# logs
logs-customer:
	docker logs -f $(CUSTOMER_CONTAINER_NAME)


# all
build-all: launch-db build-customer

launch-all: launch-docker launch-db launch-customer

clear-all: clear-db clear-customer
	@echo "rm all container!"

stop-all:
	docker stop $(DB_CONTAINER_NAME) $(CUSTOMER_CONTAINER_NAME) || true

# minimalistic
stop:  stop-all
logs:  logs-customer
clean: clear-all
build: clear-all build-all # delete db and app, rebuild them
launch:launch-all          # launch app

# main command
run:   clear-customer build-all migrate-up launch-all # fully new app (exept db)
