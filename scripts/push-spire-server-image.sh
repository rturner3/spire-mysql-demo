#!/bin/bash
# Requires login to Docker Hub account using `docker login` before running this script.

username=$(docker info | sed '/Username:/!d;s/.* //')

docker build --tag "${username}/spire-server-mysql-demo:latest" -f Dockerfile.server .
docker push "${username}/spire-server-mysql-demo:latest"
