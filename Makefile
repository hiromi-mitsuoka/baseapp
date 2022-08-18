# https://makefiletutorial.com/

# Variables
DBUSER:=root
DBPASSWORD:=password
DBPORT:=8000
DBNAME:=baseapp

SERVICE:=

.PHONY: build
	docker compose build

.PHONY: up
	docker compose up

.PHONY: down
	docker compose down

.PHONY: bash/service
	docker compose exec $(SERVICE) bin/bash