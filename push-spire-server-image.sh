#!/bin/bash
# Builds and pushes custom spire-server image to https://hub.docker.com/repository/docker/rturner0676/spire-server-mysql-demo
# Requires login to rturner0676 Docker Hub account using `docker login` before running this script.

docker build --tag rturner0676/spire-server-mysql-demo:latest -f Dockerfile.server .
docker push rturner0676/spire-server-mysql-demo:latest
