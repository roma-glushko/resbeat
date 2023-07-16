FROM nvidia/cuda:12.2.0-runtime-ubuntu20.04

RUN apt-get update -y -q && apt-get upgrade -y -q
RUN DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git

RUN curl -s https://storage.googleapis.com/golang/go1.20.4.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH $PATH:/usr/local/go/bin

WORKDIR /service
