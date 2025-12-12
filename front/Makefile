install:
	npm install

build:
	npm run build

preview:
	npm run preview

dev:
	npm run dev --prefix c:\Development\YAFDS\front

clean: clear-front
	if exist node_modules rmdir /s /q node_modules
	if exist dist rmdir /s /q dist

# Docker commands
build-docker: clean
	docker build -t front-app .

clear-front:
	-docker stop front-container
	-docker rm front-container
	-docker rmi front-app

launch-front:
	-docker start front-container 2>nul || docker run --name front-container -p 5173:5173 -d front-app

.PHONY: install build preview dev clean build-docker clear-front launch-front run

run: build-docker launch-front