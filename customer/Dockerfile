# syntax=docker/dockerfile:1

FROM golang:1.25

WORKDIR /app

COPY customer/go.mod ./

RUN go mod download

COPY migrations ./migrations
COPY customer/cmd ./cmd
COPY customer/internal ./internal
COPY customer/models ./models
COPY customer/pkg ./pkg
COPY customer/config ./config

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /customer-app ./cmd/server

EXPOSE 8091

CMD ["/customer-app"]
