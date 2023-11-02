#!/bin/bash

spire_server() {
    kubectl exec -n spire spire-server-0 -- \
        /opt/spire/bin/spire-server "${@}"
}
