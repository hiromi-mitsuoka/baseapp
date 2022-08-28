# https://makefiletutorial.com/

# Variables
DBUSER:=root
DBPASSWORD:=password
DBPORT:=8000
DBNAME:=baseapp

SERVICE:=

.PHONY: build
build:
	docker compose build

.PHONY: up
up:
	docker compose up

.PHONY: down
down:
	docker compose down

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: bash/service
bash/service:
	docker compose exec $(SERVICE) bin/bash

.PHONY: dry-migrate
dry-migrate:
	mysqldef -u user -p password -h 127.0.0.1 -P 33306 baseapp --dry-run < ./backend/_tools/mysql/schema.sql

.PHONY: migrate
migrate:
	mysqldef -u user -p password -h 127.0.0.1 -P 33306 baseapp < ./backend/_tools/mysql/schema.sql