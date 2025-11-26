include config/.env

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
	$(MAKE) create-customer-db

clear-db:
	docker stop $(DB_CONTAINER_NAME) || true
	docker rm $(DB_CONTAINER_NAME) || true

# migrations
migrate-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(MIGRATIONS_DIR) up

migrate-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(MIGRATIONS_DIR) down

migrate-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(MIGRATIONS_DIR) status

# test data
seed-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) up

seed-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) down

seed-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(CUSTOMER_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) status

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

# db's
create-customer-db:
	docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(CUSTOMER_DB);" 2>/dev/null || true
	$(MAKE) migrate-up

# logs
logs-customer:
	docker logs -f $(CUSTOMER_CONTAINER_NAME)

debug-info-customer:
	mkdir -p debug-logs
	docker cp $(CUSTOMER_CONTAINER_NAME):/app/debug_customer.txt ./debug-logs/debug_customer.txt

check-db-customers:
	docker exec -it yafds-db psql -U postgres -d customer_db -c "SELECT * FROM CUSTOMERS;"

# API test commands
test-register:
	curl -X POST http://localhost:8081/register \
		-H "Content-Type: application/json" \
		-d '{"name":"John Doe","wallet_address":"0x123456789abc","address":"123 Zhopa, Mira"}'

test-login:
	curl -X POST http://localhost:8081/login \
		-H "Content-Type: application/json" \
		-d '{"wallet_address":"0x123456789abc"}'

test-api-info:
	curl -X GET http://localhost:8081/register
	curl -X GET http://localhost:8081/login

test-all-api: test-register test-login
	@echo "\nAll API tests completed!"

# all
build-all: launch-db build-customer

launch-all: launch-docker launch-db launch-customer

clear-all: clear-db clear-customer
	@echo "All containers removed!"

stop-all:
	docker stop $(DB_CONTAINER_NAME) $(CUSTOMER_CONTAINER_NAME) || true

logs-all: debug-info-customer logs-customer

check-all: check-db-customers

# minimalistic
stop:  stop-all
logs:  logs-all
clean: clear-all
check: check-all
build: clear-all build-all
launch: launch-all

# main commands
run: clear-customer build-all migrate-up launch-all
	@echo "Application launched successfully!"

dev: clear-customer build-all migrate-up seed-up launch-all
	@echo "Application launched with test data!"
