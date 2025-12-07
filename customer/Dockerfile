# syntax=docker/dockerfile:1

FROM golang:1.25

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY models ./models
COPY pkg ./pkg
COPY config ./config

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /customer-app ./cmd/server

EXPOSE 8081

CMD ["/customer-app"]
