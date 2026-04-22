# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY courier/go.mod courier/go.sum ./courier/
WORKDIR /app/customer
COPY customer/go.mod ./
COPY customer/go.sum ./

RUN go mod download

WORKDIR /app
COPY pkg ./pkg
COPY courier ./courier
WORKDIR /app/customer
COPY customer/cmd ./cmd
COPY customer/internal ./internal
COPY customer/pkg ./pkg
COPY customer/config ./config

RUN CGO_ENABLED=0 GOOS=linux go build -o /customer-app ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /customer-app .
COPY customer/config ./config
COPY migrations ./migrations

EXPOSE 8091

CMD ["./customer-app"]
