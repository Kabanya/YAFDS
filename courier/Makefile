-include config/.env

# Set PATH to include Go binaries
ifeq ($(OS),Windows_NT)
export PATH := $(USERPROFILE)/go/bin:$(PATH)
else
GOPATH ?= $(shell go env GOPATH)
export PATH := $(GOPATH)/bin:$(PATH)
endif

# other
launch-docker:
	@docker info >/dev/null 2>&1 || echo "Warning: Docker might not be running"

copy-env:
	cp config/.env.example config/.env

# redis
launch-redis:
	-docker start $(REDIS_CONTAINER_NAME) 2>/dev/null || \
		docker run --name $(REDIS_CONTAINER_NAME) -p $(REDIS_PORT):$(REDIS_PORT) \
		-d $(REDIS_IMAGE) --port $(REDIS_PORT) --requirepass $(REDIS_PASSWORD)
	@sleep 2

clear-redis:
	-docker stop $(REDIS_CONTAINER_NAME) 2>/dev/null
	-docker rm $(REDIS_CONTAINER_NAME) 2>/dev/null

# db
launch-db:
	-docker start $(DB_CONTAINER_NAME) 2>/dev/null || docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d $(DB_IMAGE) -c port=$(DB_PORT)
	@sleep 5
	@docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) -p $(DB_PORT) >/dev/null 2>&1 || sleep 5
	$(MAKE) create-courier-db

clear-db:
	-docker stop $(DB_CONTAINER_NAME) 2>/dev/null
	-docker rm $(DB_CONTAINER_NAME) 2>/dev/null

# migrations
migrate-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -table goose_courier_version -dir $(MIGRATIONS_DIR) up

migrate-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -table goose_courier_version -dir $(MIGRATIONS_DIR) down

migrate-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -table goose_courier_version -dir $(MIGRATIONS_DIR) status

# test data
seed-up:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -table goose_courier_version -dir $(TESTDATA_MIGRATIONS_DIR) up

seed-down:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -table goose_courier_version -dir $(TESTDATA_MIGRATIONS_DIR) down

seed-status:
	goose $(DB_DRIVER) "$(DB_CONNECTION_BASE) dbname=$(COURIER_DB)" -table goose_courier_version -dir $(TESTDATA_MIGRATIONS_DIR) status

# courier
build-courier:
	-docker rmi courier-app 2>/dev/null
	docker build --build-arg BUILD_TIME=1 -t courier-app -f Dockerfile ..

launch-courier:
	-docker start $(COURIER_CONTAINER_NAME) 2>/dev/null || \
		docker run --add-host host.docker.internal:host-gateway \
			--name $(COURIER_CONTAINER_NAME) \
			-e DB_HOST=host.docker.internal \
			-e REDIS_ADDRESS=host.docker.internal:$(REDIS_PORT) \
			-e REDIS_PASSWORD=$(REDIS_PASSWORD) \
			-e REDIS_DB=$(REDIS_DB) \
			-e SESSION_TTL=$(SESSION_TTL) \
			-p $(COURIER_PORT):$(COURIER_PORT) -d courier-app

clear-courier:
	-docker stop $(COURIER_CONTAINER_NAME) 2>/dev/null
	-docker rm $(COURIER_CONTAINER_NAME) 2>/dev/null
	-docker rmi courier-app 2>/dev/null

# db's
create-courier-db:
	-docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -p $(DB_PORT) -c "CREATE DATABASE $(COURIER_DB);"
	@sleep 3
	$(MAKE) migrate-up

# logs
logs-courier:
	docker logs -f $(COURIER_CONTAINER_NAME)

debug-info-courier:
	docker exec $(COURIER_CONTAINER_NAME) cat /app/courier_log_info.txt

db-check:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(COURIER_DB) -p $(DB_PORT) -c "SELECT * FROM COURIERS;"


# all
build-all: launch-db build-courier

launch-all: launch-docker launch-redis launch-db launch-courier

clear-all: clear-redis clear-courier
	@echo "All containers removed!"

stop-all:
	-docker stop $(DB_CONTAINER_NAME) $(COURIER_CONTAINER_NAME) $(REDIS_CONTAINER_NAME)

logs-all: debug-info-courier logs-courier

check-all: db-check

# minimalistic
# Test credentials: 0x5647445a3e564b189fc35a75079fc924 / aad9D6C7e9R7
stop:  stop-all
logs:  logs-all
clean: clear-all clear-db
check: check-all
build: clear-all build-all
launch: launch-all

# main commands
run: clear-all build-all migrate-up launch-all
	@echo "Courier launched successfully!"

dev: clear-all build-all migrate-up launch-all seed-up
	@echo "Courier launched with test data!"