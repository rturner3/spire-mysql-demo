#/bin/bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source "${SCRIPT_DIR}/common.sh"

bb=$(tput bold)
nn=$(tput sgr0)

echo "${bb}Creating registration entry for mysql...${nn}"
spire_server entry create \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -spiffeID spiffe://example.org/mysql/server \
    -x509SVIDTTL 120 \
    -hint mysql-server \
    -dns mysql.mysql.svc.cluster.local \
    -selector k8s:ns:mysql \
    -selector k8s:pod-label:app:mysql-server
