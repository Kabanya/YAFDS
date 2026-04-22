install:
	npm ci || npm install

build: install
	npm run build

preview:
	npm run preview

dev: install
	npm run dev

clean: clear-front
	rm -rf node_modules dist

stop-dev:
	-lsof -ti tcp:5173 | xargs -r kill

# Docker commands
build-docker: clean
	docker build -t front-app .

clear-front:
	-docker stop yafds-front-service
	-docker rm yafds-front-service
	-docker rmi front-app

launch-front:
	-docker start yafds-front-service 2>/dev/null || docker run --name yafds-front-service -p 5173:5173 -d front-app

clear-all: clear-front
	@echo "All containers removed!"

.PHONY: install build preview dev clean stop-dev build-docker clear-front launch-front run run-local clear-all

run: build-docker launch-front

run-local: dev
