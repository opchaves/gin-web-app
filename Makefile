MODULE := $(shell go list -m)
SHELL := /bin/bash

export PROJECT = ${MODULE}

DB_CONTAINER := kom-pg
DB_IMG := 15
DB_PORT := 5432
DB_USER := root
DB_PASS := secret
DB_NAME := kom-db
DB_VOL := kom-data
DB_URL=postgresql://root:secret@localhost:5432/kom-db?sslmode=disable
REDIS_PORT := 6380
REDIS_VOL := kom-redis-data

postgres:
	docker run --rm -d --name ${DB_CONTAINER} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASS} -v ${DB_VOL}:/var/lib/postgresql/data -p ${DB_PORT}:5432 postgres:${DB_IMG}

redis:
	docker run --rm -d --name kom-redis -v ${REDIS_VOL}:/data -p ${REDIS_PORT}:6379 redis:7.2-alpine redis-server --save 60 1

data-volume:
	docker volume create ${DB_VOL}
	docker volume create ${REDIS_VOL}

db-create:
	docker exec -it ${DB_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

db-sh:
	docker exec -it ${DB_CONTAINER} psql -U ${DB_USER} ${DB_NAME}

db-drop:
	docker exec -it ${DB_CONTAINER} dropdb ${DB_NAME}

db-recreate:
	make db-drop && make db-create

db-start:
	docker start ${DB_CONTAINER}

build:
	env GOOS=linux GOARCH=amd64 go build -o bin/server ${PROJECT}/cmd
	chmod +x bin/server

build-mac:
	env GOOS=darwin GOARCH=amd64 go build -o bin/server ${PROJECT}/cmd
	env GOOS=darwin GOARCH=amd64 go build -o bin/admin ${PROJECT}/cmd/admin
	chmod +x bin/server
	chmod +x bin/admin

run:
	go run ./cmd/main.go

start:
	./bin/server

tidy:
	go mod tidy
	go mod vendor

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# count= 
# ifdef COUNT
# 	count=$(COUNT)
# endif

migrate-up-all:
	migrate -database ${DB_URL} -path ./db/migrations up

COUNT=$(error missing "COUNT" variable)
# howto: make migrate-up COUNT=2
migrate-up:
	migrate -database ${DB_URL} -path ./db/migrations up $(COUNT)

migrate-down-all:
	migrate -database ${DB_URL} -path ./db/migrations down

# howto: make migrate-down COUNT=2
migrate-down:
	migrate -database ${DB_URL} -path ./db/migrations down $(COUNT)

MIGRATION=$(error missing "MIGRATION" variable)
# howto: make migrate-create MIGRATION=migration_name
migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(MIGRATION)

sqlc-install:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

sqlc-gen:
	sqlc generate

swag-install:
	go install github.com/swaggo/swag/cmd/swag@latest

swag:
	swag init
