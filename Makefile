-include .env
export

db-up:
	docker compose up -d postgres

down:
	docker compose down

migrations-up:
	migrate -path migrations -database ${DB_URL} up

migrations-down:
	migrate -path migrations -database ${DB_URL} down

run-app:
	go run main.go

dev:
	docker compose up -d postgres
	@echo "Waiting for postgres..."
	@timeout /t 5 > NUL
	make migrations-up
	go run main.go
