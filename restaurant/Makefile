include ../.env

# Variables

DB_CONNECTION_STRING = $(DB_CONNECTION_BASE) dbname=$(1)

MIGRATIONS_DIR := ../migrations/restaurant

RESTAURANT_DB := restaurant_db

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
	$(MAKE) create-restaurant-db

clear-db:
	-docker stop $(DB_CONTAINER_NAME)
	-docker rm $(DB_CONTAINER_NAME)

# migrations
migrate-restaurant-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(MIGRATIONS_DIR) up

migrate-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(MIGRATIONS_DIR) up

migrate-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(MIGRATIONS_DIR) down

# restaurant
build-restaurant:
	docker rmi restaurant-app || true
	docker build --no-cache --build-arg BUILD_TIME=$(shell date +%s) -t restaurant-app -f Dockerfile ..

launch-restaurant:
	@if docker ps -a --format '{{.Names}}' | grep -q '^$(RESTAURANT_CONTAINER_NAME)$$'; then \
		echo "Restaurant container exists. Starting..."; \
		docker start $(RESTAURANT_CONTAINER_NAME); \
	else \
		echo "Restaurant container does not exist. Creating new one..."; \
		docker run --name $(RESTAURANT_CONTAINER_NAME) -p $(RESTAURANT_PORT):$(RESTAURANT_PORT) -d restaurant-app; \
	fi

clear-restaurant:
	-docker stop $(RESTAURANT_CONTAINER_NAME)
	-docker rm $(RESTAURANT_CONTAINER_NAME)
	-docker rmi restaurant-app


#db's
create-restaurant-db:
	docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(RESTAURANT_DB);" 2>/dev/null || true
	$(MAKE) migrate-restaurant-up

# logs
logs-restaurant:
	docker logs -f $(RESTAURANT_CONTAINER_NAME)


# all
build-all: launch-db build-restaurant

launch-all: launch-docker launch-db launch-restaurant

clear-all: clear-db clear-restaurant
	@echo "rm all container!"

stop-all:
	docker stop $(DB_CONTAINER_NAME) $(RESTAURANT_CONTAINER_NAME) || true

# minimalistic
stop:  stop-all
logs:  logs-restaurant
clean: clear-all
build: clear-all build-all # delete db and app, rebuild them
launch:launch-all          # launch app

# main command
run:   clear-restaurant build-all migrate-up launch-all # fully new app (exept db)
