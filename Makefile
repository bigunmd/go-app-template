APP_NAME = app
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/infrastructure/migrations
DATABASE_URL = postgres://postgres:postgres@127.0.0.1/postgres?sslmode=disable

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
	docker build --rm -t $(APP_NAME) .
docker.app: docker.app.build
	docker run --rm -d \
		--name $(APP_NAME) \
		--network dev-network \
		-p 8000:8000 \
		$(APP_NAME)
docker.app.stop:
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
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=postgres \
		-p 5432:5432 \
		postgres
docker.stop.postgres:
	docker stop $(APP_NAME)-postgres
docker.stop: docker.stop.redis docker.stop.postgres
docker.run: docker.network docker.redis docker.postgres
migrate.create:
	migrate create -dir $(MIGRATIONS_FOLDER) -ext .sql -seq $(migration_name)
migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up
migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down
migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)
swag:
	swag fmt -d .,./service && swag init -d .,./service -pd fiber
