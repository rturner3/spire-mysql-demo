# syntax=docker/dockerfile:1

FROM golang:1.21.3-alpine AS builder
RUN mkdir -p src/github.com/rturner3/spire-mysql-demo
WORKDIR /go/src/github.com/rturner3/spire-mysql-demo
COPY go.mod go.sum .
RUN go mod download
COPY cmd cmd
RUN go build -o /dbcredentialcomposer ./cmd/plugin/credentialcomposer/dbcredentialcomposer

FROM ghcr.io/spiffe/spire-server:1.8.2 AS base
COPY --link --from=builder --chown=1000:1000 --chmod=755 /dbcredentialcomposer /opt/spire/bin
