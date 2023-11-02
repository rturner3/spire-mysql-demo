#!/bin/bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source "${SCRIPT_DIR}/common.sh"

bb=$(tput bold)
nn=$(tput sgr0)

echo "${bb}Creating registration entry for tls-reloader...${nn}"
spire_server entry create \
    -spiffeID spiffe://example.org/mysql/client/tls-reloader \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -x509SVIDTTL 120 \
    -hint mysql-client \
    -selector k8s:ns:mysql \
    -selector k8s:pod-label:app:mysql-server
