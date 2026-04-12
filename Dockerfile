FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/getstarvio ./cmd/getstarvio

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/getstarvio /app/getstarvio
COPY --from=builder /src/database/migrations /app/database/migrations
EXPOSE 8080
ENTRYPOINT ["/app/getstarvio"]
