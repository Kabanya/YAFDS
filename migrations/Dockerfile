# syntax=docker/dockerfile:1

# Image to run Goose migrations from the root /migrations folder.
# Usage example:
#   docker build -t yafds-migrator ./migrations
#   docker run --rm --network host yafds-migrator \
#     postgres "host=localhost user=postgres password=secret sslmode=disable dbname=courier_db" \
#     -dir /migrations/courier up

FROM golang:1.23-alpine AS build

RUN apk add --no-cache git ca-certificates
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/goose /usr/local/bin/goose

WORKDIR /migrations
COPY . .

ENTRYPOINT ["goose"]
