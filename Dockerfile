FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
RUN make build

FROM debian:bullseye

WORKDIR /app
RUN apt-get update && apt-get -y upgrade

COPY --from=build-stage /app/build/ /app
