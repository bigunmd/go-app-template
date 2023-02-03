# App
## List of contents
- [App](#app)
  - [List of contents](#list-of-contents)
  - [Project requirements](#project-requirements)
  - [Install Go helpers](#install-go-helpers)
    - [gocritic](#gocritic)
    - [golangci-lint](#golangci-lint)
    - [gosec](#gosec)
    - [swag](#swag)
    - [migrate](#migrate)
  - [Makefile](#makefile)
    - [Postgres migrations](#postgres-migrations)
  - [Logger](#logger)
    - [yaml](#yaml)
    - [env](#env)
    - [args](#args)
## Project requirements
- Go 1.19
- Docker
## Install Go helpers
### gocritic
Highly extensible Go source code linter providing checks currently missing from other linters.
```bash
$ go install github.com/go-critic/go-critic/cmd/gocritic@latest
```
### golangci-lint
Fast Go linters runner. It runs linters in parallel, uses caching, supports yaml config, has integrations with all major IDE and has dozens of linters included.
```bash
$ go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
### gosec
Inspects source code for security problems by scanning the Go AST.
```bash
$ go install github.com/securego/gosec/v2/cmd/gosec@latest
```
### swag
Swag converts Go annotations to Swagger Documentation 2.0. We've created a variety of plugins for popular Go web frameworks. This allows you to quickly integrate with an existing Go project (using Swagger UI).
```bash
$ go install github.com/swaggo/swag/cmd/swag@latest
```
### migrate
Database migrations written in Go. Use as CLI or import as library.
```bash
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
## Makefile
### Postgres migrations
Default path for migrations folder is `infrastructure/migrations`

To create new migration
```bash
$ migration_name=create_user_table make migrate.create
```
To apply migrations
```bash
$ make migrate.up
```
To reverse migrations
```bash
$ make migrate.down
```
To force migrations up to a certain version
```bash
$ version=000003 make migrate.force
```
## Logger
### yaml
```yaml
logger:
  level: debug
  writeToFile: true
  file:
    path: ./
    name: app.log
    maxAge: 24h
    rotationTime: 168h
```
### env
```bash
LOGGER_LEVEL="info" # [LOG_LEVEL, LOGGING]
LOGGER_WRITE_TO_FILE=false # [LOG_WRITE_TO_FILE]
LOGGER_FILE_PATH="~/" # [LOG_FILE_PATH]
LOGGER_FILE_NAME="app.log" # [LOG_FILE_NAME]
LOGGER_FILE_MAX_AGE="24h" # [LOG_FILE_MAX_AGE]
LOGGER_FILE_ROTATION_TIME="168h" # [LOG_FILE_ROTATION_TIME]
```
### args
```bash
$ ./build/app --logger.level info
```