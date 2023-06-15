# syntax=docker/dockerfile:1

FROM golang:1.20-alpine3.17 as build

WORKDIR /service

COPY go.mod go.sum main.go /service/
RUN go mod download

COPY ./pkg /service/pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /service/resbeat

FROM alpine:3.18 AS release

WORKDIR /service

COPY --from=build /service/resbeat /service/resbeat

EXPOSE 8000

# Run
ENTRYPOINT ["/service/resbeat", "--host", "0.0.0.0"]
