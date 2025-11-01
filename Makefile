#db
run-db:
	@if docker ps -a --format '{{.Names}}' | grep -q '^yafds-db$$'; then \
		echo "Container yafds-db exists. Starting..."; \
		docker start yafds-db; \
	else \
		echo "Container yafds-db does not exist. Creating new one..."; \
		docker run --name yafds-db -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres:13.3; \
	fi

clear-db:
	docker stop yafds-db
	docker rm yafds-db

# courier
build-courier:
	docker build -t courier-app -f courier/Dockerfile courier/

run-courier:
	@if docker ps -a --format '{{.Names}}' | grep -q '^courier-container$$'; then \
		echo "Courier container exists. Starting..."; \
		docker start courier-container; \
	else \
		echo "Courier container does not exist. Creating new one..."; \
		docker run --name courier-container -p 8080:8080 -d courier-app; \
	fi

logs-courier:
	docker logs -f courier-container

clear-courier:
	docker stop courier-container
	docker rm courier-container

# customer
build-customer:
	docker build -t customer-app -f customer/Dockerfile customer/

run-customer:
	@if docker ps -a --format '{{.Names}}' | grep -q '^customer-container$$'; then \
		echo "Customer container exists. Starting..."; \
		docker start customer-container; \
	else \
		echo "Customer container does not exist. Creating new one..."; \
		docker run --name customer-container -p 8081:8081 -d customer-app; \
	fi

logs-customer:
	docker logs -f customer-container

clear-customer:
	docker stop customer-container
	docker rm customer-container

# restaurant
build-restaurant:
	docker build -t restaurant-app -f restaurant/Dockerfile restaurant/

run-restaurant:
	@if docker ps -a --format '{{.Names}}' | grep -q '^restaurant-container$$'; then \
		echo "Restaurant container exists. Starting..."; \
		docker start restaurant-container; \
	else \
		echo "Restaurant container does not exist. Creating new one..."; \
		docker run --name restaurant-container -p 8082:8082 -d restaurant-app; \
	fi

logs-restaurant:
	docker logs -f restaurant-container

clear-restaurant:
	docker stop restaurant-container
	docker rm restaurant-container

# all
build-all: build-courier build-customer build-restaurant

run-all: run-db run-courier run-customer run-restaurant

logs-all: #надо че нить сделать
	@echo "=== Showing logs from all containers ===" && \
	(docker logs -f courier-container &) && \
	(docker logs -f customer-container &) && \
	(docker logs -f restaurant-container &) && \
	wait

clear-all: clear-db clear-courier clear-customer clear-restaurant
	@echo "rm all container!"

