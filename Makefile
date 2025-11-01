# Variables
DB_HOST        := localhost
DB_USER        := postgres
DB_PASSWORD    := secret
DB_PORT        := 5432
DB_DRIVER      := postgres
MIGRATIONS_DIR := ./db/migrations

COURIER_DB     := courier_db
CUSTOMER_DB    := customer_db
RESTAURANT_DB  := restaurant_db

DB_CONNECTION_STRING = host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(1) sslmode=disable

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
	@if docker ps -a --format '{{.Names}}' | grep -q '^yafds-db$$'; then \
		echo "Container yafds-db exists. Starting..."; \
		docker start yafds-db; \
	else \
		echo "Container yafds-db does not exist. Creating new one..."; \
		docker run --name yafds-db -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres:13.3; \
	fi
	@sleep 1
	$(MAKE) create-courier-db
	$(MAKE) create-customer-db
	$(MAKE) create-restaurant-db

clear-db:
	docker stop yafds-db
	docker rm yafds-db

# migrations
migrate-courier-up:
	cd ./courier    && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(COURIER_DB)    sslmode=disable" -dir $(MIGRATIONS_DIR) up

migrate-customer-up:
	cd ./customer   && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(CUSTOMER_DB)   sslmode=disable" -dir $(MIGRATIONS_DIR) up

migrate-restaurant-up:
	cd ./restaurant && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(RESTAURANT_DB) sslmode=disable" -dir $(MIGRATIONS_DIR) up


migrate-up:
	cd ./courier    && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(COURIER_DB)    sslmode=disable" -dir $(MIGRATIONS_DIR) up
	cd ./customer   && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(CUSTOMER_DB)   sslmode=disable" -dir $(MIGRATIONS_DIR) up
	cd ./restaurant && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(RESTAURANT_DB) sslmode=disable" -dir $(MIGRATIONS_DIR) up

migrate-down:
	cd ./courier    && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(COURIER_DB)    sslmode=disable" -dir $(MIGRATIONS_DIR) down
	cd ./customer   && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(CUSTOMER_DB)   sslmode=disable" -dir $(MIGRATIONS_DIR) down
	cd ./restaurant && goose $(DB_DRIVER) "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(RESTAURANT_DB) sslmode=disable" -dir $(MIGRATIONS_DIR) down

# courier
build-courier:
	docker build -t courier-app -f courier/Dockerfile courier/

launch-courier:
	@if docker ps -a --format '{{.Names}}' | grep -q '^courier-container$$'; then \
		echo "Courier container exists. Starting..."; \
		docker start courier-container; \
	else \
		echo "Courier container does not exist. Creating new one..."; \
		docker run --name courier-container -p 8080:8080 -d courier-app; \
	fi

clear-courier:
	docker stop courier-container
	docker rm courier-container

# customer
build-customer:
	docker build -t customer-app -f customer/Dockerfile customer/

launch-customer:
	@if docker ps -a --format '{{.Names}}' | grep -q '^customer-container$$'; then \
		echo "Customer container exists. Starting..."; \
		docker start customer-container; \
	else \
		echo "Customer container does not exist. Creating new one..."; \
		docker run --name customer-container -p 8081:8081 -d customer-app; \
	fi

clear-customer:
	docker stop customer-container
	docker rm customer-container

# restaurant
build-restaurant:
	docker build -t restaurant-app -f restaurant/Dockerfile restaurant/

launch-restaurant:
	@if docker ps -a --format '{{.Names}}' | grep -q '^restaurant-container$$'; then \
		echo "Restaurant container exists. Starting..."; \
		docker start restaurant-container; \
	else \
		echo "Restaurant container does not exist. Creating new one..."; \
		docker run --name restaurant-container -p 8082:8082 -d restaurant-app; \
	fi

clear-restaurant:
	docker stop restaurant-container
	docker rm restaurant-container

#db's
create-courier-db:
	docker exec yafds-db psql -U $(DB_USER) -c "CREATE DATABASE $(COURIER_DB);" 2>/dev/null || true
	$(MAKE) migrate-courier-up

create-customer-db:
	docker exec yafds-db psql -U $(DB_USER) -c "CREATE DATABASE $(CUSTOMER_DB);" 2>/dev/null || true
	$(MAKE) migrate-customer-up

create-restaurant-db:
	docker exec yafds-db psql -U $(DB_USER) -c "CREATE DATABASE $(RESTAURANT_DB);" 2>/dev/null || true
	$(MAKE) migrate-restaurant-up

# logs
logs-courier:
	docker logs -f courier-container

logs-customer:
	docker logs -f customer-container

logs-restaurant:
	docker logs -f restaurant-container

# all
build-all: launch-db build-courier build-customer build-restaurant

launch-all: launch-db launch-courier launch-customer launch-restaurant

logs-all:
	@echo "=== Showing logs from all containers ===" && \
	($(MAKE) logs-courier &) && \
	($(MAKE) logs-customer &) && \
	($(MAKE) logs-restaurant &) && \
	wait

clear-all: clear-db clear-courier clear-customer clear-restaurant
	@echo "rm all container!"

# minimalist
run:   launch-docker launch-all migrate-up
clean: clear-all
