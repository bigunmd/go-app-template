# App [![Mentioned in Awesome Fiber](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/gofiber/awesome-fiber)

Golang application template. Based on clean architecture principles
- `domain/*` - business related logic with described usecases, repositories, presenters and services
- `infrastructure/*` - selected frameworks that are used to drive the implementation of buisness rules for usecase, repository and presenter interfaces
- `pkg/*` - additional, non-dependant logic that can be also used externaly, outside of current project scope
- `service/*` - implementation of the domain's http service interface (can be in form of cli service, websocket service, etc.)
- `user/*` - implementation of the domain's user usecase, repository and presenter interfaces
- `main.go` - application entrypoint
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
    - [Application](#application)
    - [Docker](#docker)
      - [Application](#application-1)
      - [Postgres](#postgres)
      - [Redis](#redis)
  - [Configuration](#configuration)
    - [Http](#http)
      - [yaml](#yaml)
      - [env](#env)
    - [Postgres](#postgres-1)
      - [yaml](#yaml-1)
      - [env](#env-1)
    - [Redis](#redis-1)
      - [yaml](#yaml-2)
      - [env](#env-2)
    - [Fiber](#fiber)
      - [yaml](#yaml-3)
      - [env](#env-3)
    - [Logger](#logger)
      - [yaml](#yaml-4)
      - [env](#env-4)
    - [CLI args](#cli-args)
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
Swag converts Go annotations to Swagger Documentation 2.0.
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
### Application
To test, build and run application
```bash
$ make run
```
### Docker
To start project within docker network
```bash
$ make docker.run
```
To stop
```bash
$ make docker.stop
```
#### Application
To build application docker image
```bash
$ make docker.app.build
```
You can also override build args to access GOPRIVATE repositories inside container
```
$ GOPRIVATE="" GOPRIVATE_USER="" GOPRIVATE_PAT="" GOPRIVATE_SCHEMA="" make docker.app.build
```
To start application
```bash
$ make docker.app
```
To stop application
```bash
$ make docker.stop.app
```
#### Postgres
To start postgres instance
```bash
$ make docker.postgres
```
To stop
```bash
$ make docker.stop.postgres
```
#### Redis
To start redis instance
```bash
$ make docker.redis
```
To stop
```bash
$ make docker.stop.redis
```
## Configuration
### Http
#### yaml
```yaml
http:
  host: 0.0.0.0
  port: 8000
  prefix: ""
  apiPath: "/api"
```
#### env
```bash
HTTP_HOST="0.0.0.0"
HTTP_PORT=8000
HTTP_PREFIX=""
HTTP_API_PATH="/api"
```
### Postgres
#### yaml
```yaml
postgres:
  host: 127.0.0.1
  port: 5432
  sslMode: disable
  db: postgres
  user: postgres
  password: postgres
  maxOpenConn: 0
  maxIdleConn: 2
  connMaxIdleTime: 1m
  connMaxLifetime: 1h
```
#### env
```bash
POSTGRES_HOST="127.0.0.1"
POSTGRES_PORT=5432
POSTGRES_SSL_MODE="disable"
POSTGRES_DB="postgres"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="postgres"
POSTGRES_MAX_OPEN_CONN=0
POSTGRES_MAX_IDLE_CONN=2
POSTGRES_CONN_MAX_IDLE_TIME="1m"
POSTGRES_CONN_MAX_LIFETIME="1h"
```
### Redis
#### yaml
```yaml
redis:
  auth: false
  host: 127.0.0.1
  port: 6379
  database: 0
  poolSize: 10
  username: ""
  password: ""
```
#### env
```bash
REDIS_AUTH=false
REDIS_HOST="127.0.0.1"
REDIS_PORT=6379
REDIS_DATABASE=0
REDIS_POOL_SIZE=10
REDIS_USERNAME=""
REDIS_PASSWORD=""
```
### Fiber
#### yaml
```yaml
fiber:
  prefork: false
  serverHeader: app
  strictRouting: true
  caseSensitive: true
  immutable: false
  unescapePath: false
  etag: true
  bodyLimit: 4194304
  concurrency: 262144
  readTimeout: 4s
  writeTimeout: 4s
  idleTimeout: 4s
  readBufferSize: 4096
  writeBufferSize: 4096
  compressedFileSuffix: .fiber.gz
  proxyHeader: ""
  getOnly: false
  disableKeepalive: false
  disableDefaultDate: false
  disableDefaultContentType: false
  disableHeaderNormalizing: false
  disableStartupMessage: false
  appName: app
  streamRequestBody: true
  disablePreParseMultipartForm: false
  reduceMemoryUsage: false
  network: tcp4
  enableTrustedProxyCheck: false
  trustedProxies: []
  enableIpValidation: true
  enablePrintRoutes: false
  encryptCookie:
    key: "secret-thirty-2-character-string"
  csrf:
    enable: true
    cookieSecure: true
    cookieHttpOnly: true
  cache:
    expiration: 1m
    control: true
    header: X-Cache
  limiter:
    max: 5
    expiration: 1m
    skipFailedRequests: false
    skipSuccessfulRequests: false
  cors:
    allowOrigins:
      - "*"
    allowMethods:
      - GET
      - POST
      - PUT
      - DELETE
      - PATCH
      - HEAD
    allowHeaders: []
    exposeHeaders: []
    allowCredentials: true
    maxAge: 0
  pprof:
    prefix: ""
  prometheus:
    serviceName: app
    path: /metrics
```
#### env
```bash
FIBER_PREFORK=false
FIBER_SERVER_HEADER=""
FIBER_STRICT_ROUTING=true
FIBER_CASE_SENSITIVE=true
FIBER_IMMUTABLE=false
FIBER_UNESCAPE_PATH=false
FIBER_ETAG=true
FIBER_BODY_LIMIT=4194304 # 4 Mb
FIBER_CONCURRENCY=262144 # max parallel connections
FIBER_READ_TIMEOUT="4s"
FIBER_WRITE_TIMEOUT="4s"
FIBER_IDLE_TIMEOUT="4s"
FIBER_READ_BUFFER_SIZE=4096
FIBER_WRITE_BUFFER_SIZE=4096
FIBER_COMPRESSED_FILE_SUFFIX=".fiber.gz"
FIBER_PROXY_HEADER=""
FIBER_GET_ONLY=false
FIBER_DISABLE_KEEPALIVE=false
FIBER_DISABLE_DEFAULT_DATE=false
FIBER_DISABLE_DEFAULT_CONTENT_TYPE=false
FIBER_DISABLE_HEADER_NORMALIZING=false
FIBER_DISABLE_STARTUP_MESSAGE=false
FIBER_STREAM_REQUEST_BODY=true
FIBER_DISABLE_PRE_PARSE_MULTIPART_FORM=false
FIBER_REDUCE_MEMORY_USAGE=false
FIBER_NETWORK="tcp4"
FIBER_ENABLE_TRUSTED_PROXY_CHECK=false
FIBER_TRUSTED_PROXIES="" # example: "a b c d"
FIBER_ENABLE_IP_VALIDATION=true
FIBER_ENABLE_PRINT_ROUTES=false

FIBER_ENCRYPT_COOKIE_KEY="secret-thirty-2-character-string"

FIBER_CSRF_ENABLE=false
FIBER_CSRF_COOKIE_SECURE=true
FIBER_CSRF_COOKIE_HTTP_ONLY=true

FIBER_CACHE_EXPIRATION="1m"
FIBER_CACHE_CONTROL=true
FIBER_CACHE_HEADER="X-Cache"

FIBER_LIMITER_MAX=5
FIBER_LIMITER_EXPIRATION="1m"
FIBER_LIMITER_SKIP_FAILED_REQUESTS=false
FIBER_LIMITER_SKIP_SUCCESSFUL_REQUESTS=false

FIBER_CORS_ALLOW_ORIGINS="*" # example: "a b c d"
FIBER_CORS_ALLOW_METHODS="GET POST HEAD PUT DELETE PATCH" # example: "a b c d"
FIBER_CORS_ALLOW_HEADERS="" # example: "a b c d"
FIBER_CORS_EXPOSE_HEADERS="" # example: "a b c d"
FIBER_CORS_ALLOW_CREDENTIALS=true
FIBER_CORS_MAX_AGE=0

FIBER_PPROF_PREFIX="" # access pprof information on $(prefix)/debug/pprof

FIBER_PROMETHEUS_SERVICE_NAME="app"
FIBER_PROMETHEUS_PATH="/metrics"
```
### Logger
#### yaml
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
#### env
```bash
LOGGER_LEVEL="info" # [LOG_LEVEL, LOGGING]
LOGGER_WRITE_TO_FILE=false # [LOG_WRITE_TO_FILE]
LOGGER_FILE_PATH="~/" # [LOG_FILE_PATH]
LOGGER_FILE_NAME="app.log" # [LOG_FILE_NAME]
LOGGER_FILE_MAX_AGE="24h" # [LOG_FILE_MAX_AGE]
LOGGER_FILE_ROTATION_TIME="168h" # [LOG_FILE_ROTATION_TIME]
```
### CLI args
```bash
$ ./build/app --help
```
