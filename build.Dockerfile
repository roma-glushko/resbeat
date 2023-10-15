# syntax=docker/dockerfile:1
ARG CUDA_VERSION=12.2.0
ARG BASE_DIST=ubuntu20.04

FROM nvidia/cuda:${CUDA_VERSION}-base-${BASE_DIST} as build

ARG GOLANG_VERSION=1.20.5

RUN apt-get update -y -q && apt-get upgrade -y -q
RUN DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git

RUN curl -s https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH $PATH:/usr/local/go/bin

WORKDIR /service

COPY go.mod go.sum main.go /service/
COPY ./pkg /service/pkg
