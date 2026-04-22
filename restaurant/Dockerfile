# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY customer/go.mod customer/go.sum ./customer/
WORKDIR /app/restaurant
COPY restaurant/go.mod ./
COPY restaurant/go.sum ./

RUN go mod download

WORKDIR /app
COPY pkg ./pkg
COPY customer ./customer
WORKDIR /app/restaurant
COPY restaurant/cmd ./cmd
COPY restaurant/internal ./internal
COPY restaurant/models ./models
COPY restaurant/config ./config

RUN CGO_ENABLED=0 GOOS=linux go build -o /restaurant-app ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /restaurant-app .
COPY restaurant/config ./config
COPY migrations ./migrations

EXPOSE 8092

CMD ["./restaurant-app"]
