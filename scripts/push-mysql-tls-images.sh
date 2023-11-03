#!/bin/bash
# Requires login to Docker Hub account using `docker login` before running this script.

username=$(docker info | sed '/Username:/!d;s/.* //')

docker build --tag "${username}/spire-mysql-tls-bootstrap:latest" -f Dockerfile.tlsbootstrap .
docker push "${username}/spire-mysql-tls-bootstrap:latest"

docker build --tag "${username}/spire-mysql-tls-reload:latest" -f Dockerfile.tlsreload .
docker push "${username}/spire-mysql-tls-reload:latest"
