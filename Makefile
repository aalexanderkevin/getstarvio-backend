APP_NAME := getstarvio
APP_CMD := ./cmd/getstarvio
ENV_FILE ?= .env

.PHONY: help setup dep build up down logs-db stack-build stack-up stack-down logs-app logs-worker migrate migrate-docker migrate-down-one run-api run-worker test fmt lint tidy swagger swagger-docker

help:
	@echo "Available targets:"
	@echo "  make setup            # create .env from .env.example if missing"
	@echo "  make dep              # download Go dependencies"
	@echo "  make build            # build binary to ./bin/getstarvio"
	@echo "  make up               # start postgres via docker compose"
	@echo "  make down             # stop docker compose services"
	@echo "  make stack-build      # build backend image (plain progress)"
	@echo "  make stack-up         # start app + worker + db via docker compose"
	@echo "  make stack-down       # stop full docker compose stack"
	@echo "  make logs-db          # tail postgres logs"
	@echo "  make logs-app         # tail app logs"
	@echo "  make logs-worker      # tail worker logs"
	@echo "  make migrate          # run all up migrations"
	@echo "  make migrate-docker   # run migrations via docker compose app image"
	@echo "  make migrate-down-one # rollback one migration"
	@echo "  make run-api          # run Gin API server"
	@echo "  make run-worker       # run reminder worker process"
	@echo "  make test             # run go test ./..."
	@echo "  make fmt              # gofmt project files"
	@echo "  make lint             # go vet ./..."
	@echo "  make tidy             # go mod tidy"
	@echo "  make swagger          # generate Swagger docs using local swag binary"
	@echo "  make swagger-docker   # generate Swagger docs using Dockerized swag"

setup:
	@if [ ! -f .env ]; then cp .env.example .env; echo ".env created from .env.example"; else echo ".env already exists"; fi

dep:
	@echo ">> Downloading dependencies"
	@go mod download

build:
	@echo ">> Building binary"
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux GOMAXPROCS=1 GOMEMLIMIT=300MiB go build -p=1 -v -buildvcs=false -trimpath -ldflags="-s -w" -o ./bin/getstarvio ./cmd/getstarvio

up:
	@docker compose up -d db

down:
	@docker compose down

logs-db:
	@docker compose logs -f db

stack-build:
	@docker compose build --progress=plain app

stack-up:
	@docker compose up --build -d db app worker

stack-down:
	@docker compose down

logs-app:
	@docker compose logs -f app

logs-worker:
	@docker compose logs -f worker

migrate:
	@go run $(APP_CMD) migrate

migrate-docker:
	@docker compose run --rm app migrate

migrate-down-one:
	@go run $(APP_CMD) migrate --down-one

run-api:
	@go run $(APP_CMD) server

run-worker:
	@go run $(APP_CMD) worker

test:
	@go test ./...

fmt:
	@find cmd internal -name "*.go" -type f -print0 | xargs -0 gofmt -w

lint:
	@go vet ./...

tidy:
	@go mod tidy

swagger:
	@swag init -g cmd/getstarvio/main.go -o docs --parseDependency --parseInternal

swagger-docker:
	@docker run --rm -v "$$PWD":/code -w /code golang:1.26.2 sh -lc 'go install github.com/swaggo/swag/cmd/swag@v1.16.4 && /go/bin/swag init -g cmd/getstarvio/main.go -o docs --parseDependency --parseInternal'
