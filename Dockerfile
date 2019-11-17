FROM golang:1.13-alpine as env

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct
# unfortunate, but needed if it falls back to direct, can't exclude a particular package
ENV GOSUMDB=off

RUN mkdir -p /go/src/github.com/68696c6c/twilio
WORKDIR /go/src/github.com/68696c6c/twilio


################################################################################
# Local development stage.
FROM env as dev
RUN apk add --no-cache bash
RUN GOFLAGS="" go get -u github.com/go-delve/delve/cmd/dlv
RUN echo 'alias ll="ls -lah"' >> ~/.bashrc
