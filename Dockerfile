# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-go-cook_it-api

EXPOSE 8080

CMD [ "/docker-go-cook_it-api" ]