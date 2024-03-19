include .env

migration_dir ?= "common/db/migrations"
migration_source ?= "file://$(migration_dir)"
migration_destination ?= "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

# Takes the first target as command
Command := $(firstword $(MAKECMDGOALS))
# Skips the first word
Arguments := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

migrate-create:
	migrate create -ext sql -dir $(migration_dir)  $(Arguments)

migrate-up:
	migrate -source $(migration_source) -database $(migration_destination) up

migrate-down:
	migrate -source $(migration_source) -database $(migration_destination) down

docker-up:
	docker compose up -d

docker-down:
	docker compose down
