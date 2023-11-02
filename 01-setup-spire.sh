#!/bin/bash
#
# Deploys SPIRE into a Kubernetes cluster.
# Prerequisites:
# - kubectl is installed and available on the PATH: https://kubernetes.io/docs/tasks/tools/
# - Kubernetes cluster is configured with kubectl and kubectl context is set to use this cluster

kubectl apply -k ./config/k8s/spire
echo "Waiting for SPIRE Server to be available..."
kubectl wait --for=condition=ready pod -n spire -l app=spire-server

# Create SPIRE registration entries
for entry_creation_script in ./scripts/create-*.sh; do
    bash "${entry_creation_script}"
done
