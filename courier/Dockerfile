# syntax=docker/dockerfile:1

FROM golang:1.25

WORKDIR /app

COPY customer /customer
COPY courier/go.mod ./

RUN go mod download

COPY migrations ./migrations
COPY courier/cmd ./cmd
COPY courier/internal ./internal
COPY courier/models ./models
COPY courier/pkg ./pkg
COPY courier/config ./config

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /courier-app ./cmd/server

EXPOSE 8080

CMD ["/courier-app"]
