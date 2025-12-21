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
		-d $(REDIS_IMAGE) --requirepass $(REDIS_PASSWORD) --port $(REDIS_PORT)
	@timeout /t 2 /nobreak >nul 2>&1 || sleep 2

clear-redis:
	-docker stop $(REDIS_CONTAINER_NAME)
	-docker rm $(REDIS_CONTAINER_NAME)

# db
launch-db:
	-docker start $(DB_CONTAINER_NAME) 2>nul || docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):5432 -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d $(DB_IMAGE)
	@timeout /t 5 /nobreak >nul 2>&1 || sleep 5
	@docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) >nul 2>&1 || (timeout /t 10 /nobreak >nul 2>&1 || sleep 10)
	$(MAKE) create-restaurant-db
	$(MAKE) create-order-db

clear-db:
	-docker stop $(DB_CONTAINER_NAME)
	-docker rm $(DB_CONTAINER_NAME)

# migrations
migrate-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(MIGRATIONS_DIR) up

migrate-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(MIGRATIONS_DIR) down

migrate-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(MIGRATIONS_DIR) status

# test data
seed-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) up

seed-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) down

seed-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(RESTAURANT_DB)" -dir $(TESTDATA_MIGRATIONS_DIR) status

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
	-docker start $(RESTAURANT_CONTAINER_NAME) 2>nul || \
		docker run --name $(RESTAURANT_CONTAINER_NAME) \
			-e DB_HOST=host.docker.internal \
			-e REDIS_ADDRESS=host.docker.internal:$(REDIS_PORT) \
			-e REDIS_PASSWORD=$(REDIS_PASSWORD) \
			-e REDIS_DB=$(REDIS_DB) \
			-e SESSION_TTL=$(SESSION_TTL) \
			-p $(RESTAURANT_PORT):$(RESTAURANT_PORT) -d restaurant-app

clear-restaurant:
	-docker stop $(RESTAURANT_CONTAINER_NAME)
	-docker rm $(RESTAURANT_CONTAINER_NAME)
	-docker rmi restaurant-app

# db's
create-restaurant-db:
	-docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(RESTAURANT_DB);"
	@timeout /t 3 /nobreak >nul 2>&1 || sleep 3
	$(MAKE) migrate-up

create-order-db:
	-docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c "CREATE DATABASE $(ORDER_DB);"
	@timeout /t 3 /nobreak >nul 2>&1 || sleep 3
	$(MAKE) migrate-up-orders-db

# logs
logs-restaurant:
	docker logs -f $(RESTAURANT_CONTAINER_NAME)

debug-info-restaurant:
	docker exec $(RESTAURANT_CONTAINER_NAME) cat /app/restaurant_log_info.txt

db-check:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(RESTAURANT_DB) -c "SELECT * FROM RESTAURANTS;"

# all
build-all: launch-db build-restaurant

launch-all: launch-docker launch-redis launch-db launch-restaurant

clear-all: clear-redis clear-db clear-restaurant
	@echo "All containers removed!"

stop-all:
	-docker stop $(DB_CONTAINER_NAME) $(RESTAURANT_CONTAINER_NAME) $(REDIS_CONTAINER_NAME)

logs-all: debug-info-restaurant logs-restaurant

check-all: check-db-customers

# minimalistic
stop:  stop-all
logs:  logs-all
clean: clear-all clear-db
check: check-all
build: clear-all build-all
launch: launch-all

# main commands
run: clear-all build-all migrate-up launch-all
	@echo "Application launched successfully!"

dev: clear-all build-all migrate-up seed-up launch-all
	@echo "Application launched with test data!"
