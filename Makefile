APP_NAME = app
BUILD_DIR = $(PWD)/build
CONFIG_FILE = .$(APP_NAME).config.dev.yaml

POSTGRES_HOST = 127.0.0.1
POSTGRES_PORT = 5432
POSTGRES_USER = postgres
POSTGRES_PASSWORD = postgres
POSTGRES_DB = postgres
POSTGRES_SSL_MODE = disable
POSTGRES_SEARCH_PATH = public
POSTGRES_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE)&search_path=$(POSTGRES_SEARCH_PATH)

MIGRATIONS_FOLDER = migrations/$(APP_NAME)

DOCKER_PATH = ./docker/Dockerfile
DOCKER_TAG = test
DOCKER_NETWORK = dev-network

PRIVATE_REGISTRY_PATH = ""

GOPRIVATE_USER = "__token__"
GOPRIVATE_PAT = ""
GOPRIVATE = ""
GOPRIVATE_SCHEMA = "https"

clean:
	rm -rf $(BUILD_DIR)
critic:
	gocritic check -enableAll main
security:
	gosec ./...
lint:
	golangci-lint run ./...
test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover -p 1 ./...
	go tool cover -html=cover.out -o coverage.html
build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)/main.go
run: build
	$(BUILD_DIR)/$(APP_NAME)
run.go:
	go run ./cmd/$(APP_NAME)/main.go -c $(CONFIG_FILE)
swag:
	swag fmt -d ./internal && swag init -d ./cmd/$(APP_NAME),./internal/$(APP_NAME),./internal/controller/http -pd fiber

migrate.create:
	migrate create -dir $(MIGRATIONS_FOLDER) -ext .sql -seq $(NAME) -v 
migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(POSTGRES_URL)&x-migrations-table=$(POSTGRES_SEARCH_PATH)_migrations" -verbose up
migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(POSTGRES_URL)&x-migrations-table=$(POSTGRES_SEARCH_PATH)_migrations" -verbose down
migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(POSTGRES_URL)&x-migrations-table=$(POSTGRES_SEARCH_PATH)_migrations" force $(VERSION) -v

gen.clean:
	rm -rf gen/*
gen.sqlc:
	sqlc generate
gen: gen.clean gen.sqlc

docker.network:
	docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || \
	docker network create -d bridge $(DOCKER_NETWORK)
docker.run.postgres:
	docker run --rm -d \
		--name $(APP_NAME)-postgres \
		--network $(DOCKER_NETWORK) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		postgres
docker.stop.postgres:
	docker stop $(APP_NAME)-postgres
docker.run.redis:
	docker run --rm -d \
		--name $(APP_NAME)-redis \
		--network $(DOCKER_NETWORK) \
		-p 6379:6379 \
		redis:7-alpine
docker.stop.redis:
	docker stop $(APP_NAME)-redis
docker.remote.push:
	docker buildx build --rm \
	-t $(PRIVATE_REGISTRY_PATH)/$(APP_NAME):$(DOCKER_TAG) \
	--build-arg APP_NAME=$(APP_NAME) \
	--build-arg GOPRIVATE=$(GOPRIVATE) \
	--build-arg GOPRIVATE_USER=$(GOPRIVATE_USER) \
	--build-arg GOPRIVATE_PAT=$(GOPRIVATE_PAT) \
	--build-arg GOPRIVATE_SCHEMA=$(GOPRIVATE_SCHEMA) \
	-f $(DOCKER_PATH) . --push
docker.build.app:
	docker build --rm \
	-t $(APP_NAME):$(DOCKER_TAG) \
	--build-arg APP_NAME=$(APP_NAME) \
	--build-arg GOPRIVATE=$(GOPRIVATE) \
	--build-arg GOPRIVATE_USER=$(GOPRIVATE_USER) \
	--build-arg GOPRIVATE_PAT=$(GOPRIVATE_PAT) \
	--build-arg GOPRIVATE_SCHEMA=$(GOPRIVATE_SCHEMA) \
	-f $(DOCKER_PATH) .
docker.run.app: docker.build.app
	docker run --rm \
		--name $(APP_NAME) \
		--network $(DOCKER_NETWORK) \
		-e REDIS_HOST=$(APP_NAME)-redis \
		-p 8000:8000 \
		$(APP_NAME):$(DOCKER_TAG)
docker.stop.app:
	docker stop $(APP_NAME)
docker.stop: docker.stop.postgres docker.stop.redis
docker.run: docker.network docker.run.postgres docker.run.redis
