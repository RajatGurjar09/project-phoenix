COMPOSE_PROJECT_NAME ?= project-phoenix
MIGRATE_IMAGE ?= migrate/migrate:v4.17.1
POSTGRES_DB ?= phoenix
POSTGRES_USER ?= phoenix
POSTGRES_PASSWORD ?=
DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable
MIGRATIONS_PATH ?= /migrations

MIGRATE = docker run --rm \
	--network $(COMPOSE_PROJECT_NAME)_default \
	-v $(CURDIR)/backend/migrations:$(MIGRATIONS_PATH) \
	$(MIGRATE_IMAGE) \
	-path=$(MIGRATIONS_PATH) \
	-database "$(DATABASE_URL)"

.PHONY: migrate-up migrate-down migrate-version

migrate-up:
	@test -n "$(POSTGRES_PASSWORD)" || (echo "POSTGRES_PASSWORD must be set" && exit 1)
	$(MIGRATE) up

migrate-down:
	@test -n "$(POSTGRES_PASSWORD)" || (echo "POSTGRES_PASSWORD must be set" && exit 1)
	$(MIGRATE) down 1

migrate-version:
	@test -n "$(POSTGRES_PASSWORD)" || (echo "POSTGRES_PASSWORD must be set" && exit 1)
	$(MIGRATE) version
