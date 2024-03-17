include .env

migration_source ?= "file://db/migrations"
migration_destination ?= "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

# Takes the first target as command
Command := $(firstword $(MAKECMDGOALS))
# Skips the first word
Arguments := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(Arguments)

migrate-up:
	migrate -source $(migration_source) -database $(migration_destination) up

migrate-down:
	migrate -source $(migration_source) -database $(migration_destination) down

docker-up:
	docker compose up -d

docker-down:
	docker compose down
