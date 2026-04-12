# syntax=docker/dockerfile:1.7
FROM golang:1.26.2 AS builder
WORKDIR /app

COPY go.mod go.sum* Makefile ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make dep

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOMAXPROCS=1 GOMEMLIMIT=300MiB make build

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/getstarvio /app/getstarvio
COPY --from=builder /app/database/migrations /app/database/migrations
EXPOSE 8080
ENTRYPOINT ["/app/getstarvio"]
