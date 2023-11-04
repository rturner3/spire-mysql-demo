#!/bin/bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source "${SCRIPT_DIR}/common.sh"

bb=$(tput bold)
nn=$(tput sgr0)

forty_eight_hours_in_seconds=172800

echo "${bb}Creating registration entry for spire-mysql-client...${nn}"
spire_server entry create \
    -spiffeID spiffe://example.org/mysql/client/spire-mysql-client \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -x509SVIDTTL "${forty_eight_hours_in_seconds}" \
    -hint mysql-client \
    -selector k8s:ns:default \
    -selector k8s:pod-label:app:spire-mysql-client
