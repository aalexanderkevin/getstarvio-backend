FROM golang:1.26 AS builder
WORKDIR /app

COPY . .
RUN make build

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/getstarvio /app/getstarvio
COPY --from=builder /app/database/migrations /app/database/migrations
EXPOSE 8080
ENTRYPOINT ["/app/getstarvio"]
