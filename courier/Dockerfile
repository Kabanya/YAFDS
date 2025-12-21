# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY customer /customer
COPY courier/go.mod ./

RUN go mod download

COPY courier/cmd ./cmd
COPY courier/internal ./internal
COPY courier/models ./models
COPY courier/pkg ./pkg
COPY courier/config ./config

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /courier-app ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /courier-app .
COPY courier/config ./config
COPY migrations ./migrations

EXPOSE 8090

CMD ["./courier-app"]
