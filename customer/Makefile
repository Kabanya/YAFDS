include config/.env

# Set PATH to include Go binaries
export PATH := $(USERPROFILE)/go/bin:$(PATH)

# other
launch-docker:
	@docker info >nul 2>&1 || echo "Warning: Docker might not be running"

# redis
launch-redis:
	-docker start $(REDIS_CONTAINER_NAME) 2>nul || \
		docker run --name $(REDIS_CONTAINER_NAME) -p $(REDIS_PORT):6379 \
		-d $(REDIS_IMAGE) --requirepass $(REDIS_PASSWORD)
	@timeout /t 2 /nobreak >nul 2>&1 || sleep 2

clear-redis:
	-docker stop $(REDIS_CONTAINER_NAME)
	-docker rm $(REDIS_CONTAINER_NAME)

# db
launch-db:
	-docker start $(DB_CONTAINER_NAME) 2>nul || docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d $(DB_IMAGE) -c port=$(DB_PORT)
	@timeout /t 2 /nobreak >nul 2>&1 || sleep 5
	@docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) -p $(DB_PORT) >nul 2>&1 || (timeout /t 5 /nobreak >nul 2>&1 || sleep 5)
	$(MAKE) create-customer-db

clear-db:
	-docker stop $(DB_CONTAINER_NAME)
	-docker rm $(DB_CONTAINER_NAME)

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
	-docker rmi customer-app 2>nul
	docker build --build-arg BUILD_TIME=1 -t customer-app -f Dockerfile ..

launch-customer:
	-docker start $(CUSTOMER_CONTAINER_NAME) 2>nul || \
		docker run --name $(CUSTOMER_CONTAINER_NAME) \
			-e DB_HOST=host.docker.internal \
			-e REDIS_ADDRESS=host.docker.internal:$(REDIS_PORT) \
			-e REDIS_PASSWORD=$(REDIS_PASSWORD) \
			-e REDIS_DB=$(REDIS_DB) \
			-e SESSION_TTL=$(SESSION_TTL) \
			-p $(CUSTOMER_PORT):$(CUSTOMER_PORT) -d customer-app

clear-customer:
	-docker stop $(CUSTOMER_CONTAINER_NAME)
	-docker rm $(CUSTOMER_CONTAINER_NAME)
	-docker rmi customer-app

# db's
create-customer-db:
	-docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -p $(DB_PORT) -c "CREATE DATABASE $(CUSTOMER_DB);"
	@timeout /t 3 /nobreak >nul 2>&1 || sleep 3
	$(MAKE) migrate-up

# logs
logs-customer:
	docker logs -f $(CUSTOMER_CONTAINER_NAME)

debug-info-customer:
	docker exec $(CUSTOMER_CONTAINER_NAME) cat /app/customer_log_info.txt

db-check:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(CUSTOMER_DB) -p $(DB_PORT) -c "SELECT * FROM CUSTOMERS;"


# all
build-all: launch-db build-customer

launch-all: launch-docker launch-redis launch-db launch-customer

clear-all: clear-redis  clear-customer
	@echo "All containers removed!"

stop-all:
	-docker stop $(DB_CONTAINER_NAME) $(CUSTOMER_CONTAINER_NAME) $(REDIS_CONTAINER_NAME)

logs-all: debug-info-customer logs-customer

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
