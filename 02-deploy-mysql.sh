#!/bin/bash
# Deploy MySQL to Kubernetes.
# Prerequisites:
# - kubectl is installed and available on the PATH: https://kubernetes.io/docs/tasks/tools/
# - Kubernetes cluster is configured with kubectl and kubectl context is set to use this cluster

# Deploy MySQL database
kubectl apply -k ./config/k8s/mysql

# Wait for DB to be up
kubectl wait --for=condition=ready pod -n mysql -l app=mysql-server
