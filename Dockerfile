# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

MAINTAINER maria jose davila <mariajdab@gmail.com>

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go get -d -v
RUN go build -v
RUN go build -o /post-api-rest

EXPOSE 8080

CMD ["./post-api-rest" ]