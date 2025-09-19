DB_URL=postgres://sca_user:sca_pass@localhost:5432/sca_db?sslmode=disable

.PHONY: up down migrate-up migrate-down reset

up:
	docker-compose up -d db
	@echo "Postgres is starting... wait a few seconds before migrating."

down:
	docker-compose down

migrate-up:
	docker-compose run --rm migrate up

migrate-down:
	docker-compose run --rm migrate down 1

reset: down up migrate-up
