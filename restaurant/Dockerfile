# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY customer /customer
COPY restaurant/go.mod ./

RUN go mod download

COPY restaurant/cmd ./cmd
COPY restaurant/internal ./internal
COPY restaurant/models ./models
COPY customer/pkg ./pkg
COPY restaurant/config ./config

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /restaurant-app ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /restaurant-app .
COPY restaurant/config ./config
COPY migrations ./migrations

EXPOSE 8092

CMD ["./restaurant-app"]
