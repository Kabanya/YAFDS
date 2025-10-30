# syntax=docker/dockerfile:1

FROM golang:1.25

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go *.sql ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /customer-app .

EXPOSE 8081

CMD ["/customer-app"]