DB_URL=postgres://sca_user:sca_pass@db:5432/sca_db?sslmode=disable

.PHONY: up down migrate-up migrate-down reset

up:
	docker-compose up -d db
	@echo "Waiting for Postgres to be ready..."
	@until docker-compose exec db pg_isready -U sca_user -d sca_db; do sleep 1; done
	@echo "Postgres is ready!"

down:
	docker-compose down

migrate-up:
	docker-compose run --rm migrate up

migrate-down:
	docker-compose run --rm migrate down 1

reset: down up migrate-up
