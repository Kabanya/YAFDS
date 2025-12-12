# syntax=docker/dockerfile:1

FROM golang:1.25

ARG BUILD_TIME
RUN echo "Build time: $BUILD_TIME"

WORKDIR /app

COPY courier/go.mod ./

RUN go mod download

COPY migrations ./migrations
COPY courier/ .

RUN date > /build-time.txt

RUN CGO_ENABLED=0 GOOS=linux go build -o /courier-app ./cmd/server

EXPOSE 8080

CMD ["/courier-app"]