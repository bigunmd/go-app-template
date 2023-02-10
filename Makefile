APP_NAME = app
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/infrastructure/migrations
POSTGRES_HOST = 127.0.0.1
POSTGRES_PORT = 5432
POSTGRES_USER = postgres
POSTGRES_PASSWORD = postgres
POSTGRES_DB = postgres
POSTGRES_SSL_MODE=disable
POSTGRES_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE)
GOPRIVATE_USER = "__token__"
GOPRIVATE_PAT = ""
GOPRIVATE = ""
GOPRIVATE_SCHEMA = "https"

clean:
	rm -rf ./build
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
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go
run: swag build
	$(BUILD_DIR)/$(APP_NAME)
docker.app.build:
	docker build --rm \
	-t $(APP_NAME) \
	--build-arg GOPRIVATE=$(GOPRIVATE) \
	--build-arg GOPRIVATE_USER=$(GOPRIVATE_USER) \
	--build-arg GOPRIVATE_PAT=$(GOPRIVATE_PAT) \
	--build-arg GOPRIVATE_SCHEMA=$(GOPRIVATE_SCHEMA) \
	.
docker.app: docker.app.build
	docker run --rm -d \
		--name $(APP_NAME) \
		--network dev-network \
		-e POSTGRES_HOST=$(APP_NAME)-postgres \
		-e REDIS_HOST=$(APP_NAME)-redis \
		-p 8000:8000 \
		$(APP_NAME)
docker.stop.app:
	docker stop $(APP_NAME)
docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network
docker.redis:
	docker run --rm -d \
		--name $(APP_NAME)-redis \
		--network dev-network \
		-p 6379:6379 \
		redis:7-alpine
docker.stop.redis:
	docker stop $(APP_NAME)-redis
docker.postgres:
	docker run --rm -d \
		--name $(APP_NAME)-postgres \
		--network dev-network \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		postgres
docker.stop.postgres:
	docker stop $(APP_NAME)-postgres
docker.stop: docker.stop.redis docker.stop.postgres docker.stop.app
docker.run: docker.network docker.redis docker.postgres docker.app
migrate.create:
	migrate create -dir $(MIGRATIONS_FOLDER) -ext .sql -seq $(migration_name)
migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(POSTGRES_URL)" up
migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(POSTGRES_URL)" down
migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(POSTGRES_URL)" force $(version)
swag:
	swag fmt -d .,./user && swag init -d .,./user -pd fiber
