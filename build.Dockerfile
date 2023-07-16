# syntax=docker/dockerfile:1
FROM golang:1.20 as build

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3

WORKDIR /service

COPY go.mod go.sum main.go /service/
RUN go mod download
