FROM alpine:latest as os

ARG APP_VER=1.0.0
ARG APP_NAME=integrasi-mitra-grpc

ENV APP_NAME=$APP_NAME
ENV APP_VER=$APP_VER
ENV APP_PORT=8080
ENV APP_MODE=release
ENV GRPC_SERVER_MITRA_INTEGRASI=localhost:50051
ENV GRPC_SERVER_PENGERAHAN=localhost:50052
ENV COCKCROACH_DB_CONN_STR=postgresql://taperadev:%23Kecoanuklir2020%23@10.172.31.2:26257/sitara_dev

FROM golang:1.14-alpine AS build

ENV GO111MODULE=auto

WORKDIR /app/src/internal
COPY ./internal/. .

WORKDIR /app/src/grpc
#COPY ./grpc/go.build.mod ./go.mod
COPY ./grpc/go.mod .
COPY ./grpc/go.sum .

RUN go mod tidy -v
RUN go mod download -x

COPY ./grpc/. .
RUN rm .env

RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main

RUN rm -rf /app/src

FROM os
WORKDIR /app
COPY --from=build /app .
ENTRYPOINT ["./main"]

EXPOSE $APP_PORT