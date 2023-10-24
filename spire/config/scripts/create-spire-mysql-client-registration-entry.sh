#!/bin/bash

set -e

bb=$(tput bold)
nn=$(tput sgr0)

echo "${bb}Creating registration entry for spire-mysql-client...${nn}"
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/mysql/client/spire-mysql-client \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -ttl 2m \
    -hint mysql-client \
    -selector k8s:ns:default \
    -selector k8s:pod-label:app:spire-mysql-client
