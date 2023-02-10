FROM golang:1.19 AS builder

ARG GOPRIVATE_USER="__token__"
ARG GOPRIVATE_PAT=""
ARG GOPRIVATE=""
ARG GOPRIVATE_SCHEMA="https"

RUN git config --global url."${GOPRIVATE_SCHEMA}://${GOPRIVATE_USER}:${GOPRIVATE_PAT}@${GOPRIVATE}/".insteadOf ${GOPRIVATE_SCHEMA}://${GOPRIVATE}/

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o app

FROM scratch

COPY --from=builder ["/build/app", "/"]

ENTRYPOINT ["/app"]
