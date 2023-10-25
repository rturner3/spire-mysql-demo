#!/bin/bash

undeploy_k8s_components() {
    dir=$1
    kustomize build "${dir}" | kubectl delete -f -
}

undeploy_k8s_components ./config/k8s/sample-service
undeploy_k8s_components ./config/k8s/mysql
./scripts/delete-entries.sh
undeploy_k8s_components ./config/k8s/spire
