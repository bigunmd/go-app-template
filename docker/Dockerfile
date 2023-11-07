ARG APP_NAME="app"

# Step 1: Modules caching
FROM golang:1.21.1-alpine3.18 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21.1-alpine3.18 as builder

ARG APP_NAME

COPY --from=modules /go/pkg /go/pkg
COPY . /${APP_NAME}
WORKDIR /${APP_NAME}
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/${APP_NAME} ./cmd/${APP_NAME}

# Step 3: Final
FROM scratch

ARG APP_NAME

COPY --from=builder /bin/${APP_NAME} /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]