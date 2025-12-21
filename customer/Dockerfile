# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY customer/go.mod ./

RUN go mod download

COPY customer/cmd ./cmd
COPY customer/internal ./internal
COPY customer/models ./models
COPY customer/pkg ./pkg
COPY customer/config ./config

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /customer-app ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /customer-app .
COPY customer/config ./config
COPY migrations ./migrations

EXPOSE 8091

CMD ["./customer-app"]
