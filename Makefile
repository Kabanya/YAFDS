run-db:
	docker run --name yafds-db -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres:13.3