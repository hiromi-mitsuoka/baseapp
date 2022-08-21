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

.PHONY: bash/service
bash/service:
	docker compose exec $(SERVICE) bin/bash