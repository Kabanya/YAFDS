include config/.env

# Set PATH to include Go binaries
export PATH := $(USERPROFILE)/go/bin:$(PATH)

# other
launch-docker:
	@docker info >nul 2>&1 || echo "Warning: Docker might not be running"

# redis
launch-redis:
	-docker start $(REDIS_CONTAINER_NAME) 2>nul || \
		docker run --name $(REDIS_CONTAINER_NAME) -p $(REDIS_PORT):$(REDIS_PORT) \
		-d $(REDIS_IMAGE) --requirepass $(REDIS_PASSWORD)
	@timeout /t 2 /nobreak >nul 2>&1 || sleep 2

clear-redis:
	-docker stop $(REDIS_CONTAINER_NAME)
	-docker rm $(REDIS_CONTAINER_NAME)

# db
launch-db:
	-docker start $(DB_CONTAINER_NAME) 2>nul || docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d $(DB_IMAGE)
	@timeout /t 15 /nobreak >nul 2>&1 || sleep 15
	@docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) >nul 2>&1 || (timeout /t 10 /nobreak >nul 2>&1 || sleep 10)
	$(MAKE) create-restaurant-db
	$(MAKE) create-order-db

clear-db:
	-docker stop $(DB_CONTAINER_NAME)
	-docker rm $(DB_CONTAINER_NAME)

# migrations
migrate-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RETAURANT_DB)" -dir $(MIGRATIONS_DIR) up

migrate-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RETAURANT_DB)" -dir $(MIGRATIONS_DIR) down

migrate-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RETAURANT_DB)" -dir $(MIGRATIONS_DIR) status

# test data
seed-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RETAURANT_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) up

seed-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RETAURANT_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) down

seed-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RETAURANT_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) status

# migrations
migrate-up-orders-db:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(ORDER_DB)" -dir $(ORDERS_MIGRATIONS_DIR) up

migrate-down-orders-db:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(ORDER_DB)" -dir $(ORDERS_MIGRATIONS_DIR) down

migrate-status-orders-db:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(ORDER_DB)" -dir $(ORDERS_MIGRATIONS_DIR) status

# test data
seed-up-orders-db:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(ORDER_DB)" -dir $(TESTDATA_ORDERS_MIGRATIONS_DIR) up

seed-down-orders-db:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(ORDER_DB)" -dir $(TESTDATA_ORDERS_MIGRATIONS_DIR) down

seed-status-orders-db:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(ORDER_DB)" -dir $(TESTDATA_ORDERS_MIGRATIONS_DIR) status


# restaurant
build-restaurant:
	-docker rmi restaurant-app 2>nul
	docker build --build-arg BUILD_TIME=1 -t restaurant-app -f Dockerfile ..

launch-restaurant:
	-docker start $(RETAURANT_CONTAINER_NAME) 2>nul || \
		docker run --name $(RETAURANT_CONTAINER_NAME) \
			-e DB_HOST=host.docker.internal \
			-e REDIS_ADDRESS=host.docker.internal:$(REDIS_PORT) \
			-e REDIS_PASSWORD=$(REDIS_PASSWORD) \
			-e REDIS_DB=$(REDIS_DB) \
			-e SESSION_TTL=$(SESSION_TTL) \
			-p $(RETAURANT_PORT):$(RETAURANT_PORT) -d restaurant-app

clear-restaurant:
	-docker stop $(RETAURANT_CONTAINER_NAME)
	-docker rm $(RETAURANT_CONTAINER_NAME)
	-docker rmi restaurant-app

# db's
create-restaurant-db:
	-docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(RETAURANT_DB);"
	@timeout /t 3 /nobreak >nul 2>&1 || sleep 3
	$(MAKE) migrate-up

create-order-db:
	-docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(ORDER_DB);"
	@timeout /t 3 /nobreak >nul 2>&1 || sleep 3
	$(MAKE) migrate-up-orders-db

# logs
logs-restaurant:
	docker logs -f $(RETAURANT_CONTAINER_NAME)

debug-info-restaurant:
	docker exec $(RETAURANT_CONTAINER_NAME) cat /app/restaurant_log_info.txt

db-check:
	docker exec -it yafds-db psql -U postgres -d COURIER_DB -c "SELECT * FROM CUSTOMERS;"

# API test commands
test-register:
	curl -X POST http://localhost:$(RETAURANT_PORT)/register \
		-H "Content-Type: application/json" \
		-d '{"name":"John Doe","wallet_address":"0x123456789abc","address":"123 Zhopa, Mira"}'

test-login:
	curl -X POST http://localhost:$(RETAURANT_PORT)/login \
		-H "Content-Type: application/json" \
		-d '{"wallet_address":"0x123456789abc"}'

test-api-info:
	curl -X GET http://localhost:$(RETAURANT_PORT)/register
	curl -X GET http://localhost:$(RETAURANT_PORT)/login

test-all-api: test-register test-login
	@echo "\nAll API tests completed!"

# all
build-all: launch-db build-restaurant

launch-all: launch-docker launch-redis launch-db launch-restaurant

clear-all: clear-redis clear-db clear-restaurant
	@echo "All containers removed!"

stop-all:
	-docker stop $(DB_CONTAINER_NAME) $(RETAURANT_CONTAINER_NAME) $(REDIS_CONTAINER_NAME)

logs-all: debug-info-restaurant logs-restaurant

check-all: check-db-customers

# minimalistic
stop:  stop-all
logs:  logs-all
clean: clear-all
check: check-all
build: clear-all build-all
launch: launch-all

# main commands
run: clean build-all migrate-up launch-all
	@echo "Application launched successfully!"

dev: clean build-all migrate-up seed-up launch-all
	@echo "Application launched with test data!"
