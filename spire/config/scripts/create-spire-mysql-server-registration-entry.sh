#/bin/bash

set -e

bb=$(tput bold)
nn=$(tput sgr0)


echo "${bb}Creating registration entry for mysql...${nn}"
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/mysql/server \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -x509SVIDTTL 120 \
    -hint mysql-server \
    -selector k8s:ns:mysql \
    -selector k8s:pod-label:app:mysql-server
