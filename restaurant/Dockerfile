# syntax=docker/dockerfile:1

FROM golang:1.25

WORKDIR /app

COPY customer /customer
COPY restaurant/go.mod ./

RUN go mod download

COPY migrations ./migrations
COPY restaurant/cmd ./cmd
COPY restaurant/internal ./internal
COPY restaurant/models ./models
COPY restaurant/pkg ./pkg
COPY restaurant/config ./config

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /restaurant-app ./cmd/server

EXPOSE 8082

CMD ["/restaurant-app"]
