FROM golang:1.13-alpine as env

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct
# unfortunate, but needed if it falls back to direct, can't exclude a particular package
ENV GOSUMDB=off

RUN apk add --no-cache git gcc python bash openssh

RUN mkdir -p /go/src/github.com/68696c6c/twilio-client
WORKDIR /go/src/github.com/68696c6c/twilio-client

################################################################################
# Local development stage.
FROM env as dev
RUN GOFLAGS="" go get -u github.com/go-delve/delve/cmd/dlv
RUN echo 'alias ll="ls -lah"' >> ~/.bashrc


################################################################################
# Stage for running tests, skipping build for speed.
FROM env as base

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go mod vendor


################################################################################
# Builder stage, compile the app.
FROM base as builder

RUN go build -o api
