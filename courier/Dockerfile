# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY customer/go.mod customer/go.sum ./customer/
WORKDIR /app/courier
COPY courier/go.mod ./
COPY courier/go.sum ./

RUN go mod download

WORKDIR /app
COPY pkg ./pkg
COPY customer ./customer
WORKDIR /app/courier
COPY courier/cmd ./cmd
COPY courier/internal ./internal
COPY courier/pkg ./pkg
COPY courier/config ./config

RUN CGO_ENABLED=0 GOOS=linux go build -o /courier-app ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /courier-app .
COPY courier/config ./config
COPY migrations ./migrations

EXPOSE 8093

CMD ["./courier-app"]
