# syntax=docker/dockerfile:1
FROM golang:1.20 as build

WORKDIR /service

COPY go.mod go.sum main.go /service/
RUN go mod download
