# Go Blog API – Makefile
# Use: make [target]. Default: make help

BINARY_NAME := api
MAIN_PATH   := ./cmd/api
GO          := go
DOCKER      := docker compose

.PHONY: help build run test tidy clean swag docker-up docker-down docker-build docker-rebuild install

help:
	@echo "Go Blog API – targets:"
	@echo "  make build        - Build binary ($(BINARY_NAME))"
	@echo "  make run          - Run API (go run)"
	@echo "  make test         - Run tests"
	@echo "  make tidy         - go mod tidy"
	@echo "  make swag         - Regenerate Swagger docs (docs/)"
	@echo "  make clean        - Remove binary and temp files"
	@echo "  make install      - go mod download"
	@echo "  make docker-up    - Start API + Postgres (docker compose up -d)"
	@echo "  make docker-build - Build and start (docker compose up --build -d)"
	@echo "  make docker-down  - Stop containers (docker compose down)"
	@echo "  make docker-rebuild - No-cache build and start"

build:
	$(GO) build -o $(BINARY_NAME) $(MAIN_PATH)

run:
	$(GO) run $(MAIN_PATH)

test:
	$(GO) test -v ./...

tidy:
	$(GO) mod tidy

swag:
	$(GO) run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -d . -o docs

clean:
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe 2>/dev/null || true
	@echo "Cleaned."

install:
	$(GO) mod download

docker-up:
	$(DOCKER) up -d

docker-build:
	$(DOCKER) up --build -d

docker-down:
	$(DOCKER) down

drb:
	$(DOCKER) build --no-cache
	$(DOCKER) up -d
