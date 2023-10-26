#!/bin/bash
# Requires login to chiragk25 Docker Hub account using `docker login` before running this script.

# Builds and pushes tls-bootstrap image to https://hub.docker.com/repository/docker/chiragk25/spire-server-tls-bootstrap
docker build --tag chiragk25/spire-mysql-tls-bootstrap:latest -f Dockerfile.tlsbootstrap .
docker push chiragk25/spire-server-tls-bootstrap:latest

# Builds and pushes tls-reload image to https://hub.docker.com/repository/docker/chiragk25/spire-server-tls-reload
docker build --tag chiragk25/spire-mysql-tls-reload:latest -f Dockerfile.tlsreload .
docker push chiragk25/spire-server-tls-reload:latest
