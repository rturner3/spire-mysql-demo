#!/bin/bash
#
# Deploys sample service that accesses a MySQL database into a Kubernetes cluster.
# Prerequisites:
# - kubectl is installed and available on the PATH: https://kubernetes.io/docs/tasks/tools/
# - Kubernetes cluster is configured with kubectl and kubectl context is set to use this cluster

kubectl apply -k ./config/k8s/sample-service
