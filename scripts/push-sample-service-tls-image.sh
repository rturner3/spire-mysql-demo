#!/bin/bash
# Requires login to Docker Hub account using `docker login` before running this script.

username=$(docker info | sed '/Username:/!d;s/.* //')

docker build --tag "${username}/spire-mysql-sample-service:latest" -f Dockerfile.sampleservice .
docker push "${username}/spire-mysql-sample-service:latest"
