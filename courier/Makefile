include ../.env

# Variables

DB_CONNECTION_STRING = $(DB_CONNECTION_BASE) dbname=$(1)

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
		docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:17.6; \
	fi
	@sleep 1
	$(MAKE) create-courier-db

clear-db:
	docker stop $(DB_CONTAINER_NAME) || true
	docker rm $(DB_CONTAINER_NAME) || true

# migrations
migrate-courier-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -dir $(MIGRATIONS_DIR) up

migrate-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -dir $(MIGRATIONS_DIR) up

migrate-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -dir $(MIGRATIONS_DIR) down

# courier
build-courier:
	docker rmi courier-app || true
	docker build --no-cache --build-arg BUILD_TIME=$(shell date +%s) -t courier-app -f Dockerfile .

launch-courier:
	@if docker ps -a --format '{{.Names}}' | grep -q '^$(COURIER_CONTAINER_NAME)$$'; then \
		echo "Courier container exists. Starting..."; \
		docker start $(COURIER_CONTAINER_NAME); \
	else \
		echo "Courier container does not exist. Creating new one..."; \
		docker run --name $(COURIER_CONTAINER_NAME) -p $(COURIER_PORT):$(COURIER_PORT) -d courier-app; \
	fi

clear-courier:
	docker stop $(COURIER_CONTAINER_NAME) || true
	docker rm $(COURIER_CONTAINER_NAME) || true


#db's
create-courier-db:
	docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(COURIER_DB);" 2>/dev/null || true
	$(MAKE) migrate-courier-up

# logs
logs-courier:
	docker logs -f $(COURIER_CONTAINER_NAME)


# all
build-all: launch-db build-courier

launch-all: launch-docker launch-db launch-courier

clear-all: clear-db clear-courier
	@echo "rm all container!"

stop-all:
	docker stop $(DB_CONTAINER_NAME) $(COURIER_CONTAINER_NAME) || true

# minimalistic
stop:  stop-all
logs:  logs-courier
clean: clear-all
build: clear-all build-all # delete db and app, rebuild them
launch:launch-all          # launch app

# main command
run:   clear-courier build-all migrate-up launch-all # fully new app (exept db)
