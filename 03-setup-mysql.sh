#!/bin/bash
#
# Sets up schemas in MySQL database for sample service.
# Prerequisites:
# - kubectl is installed and available on the PATH: https://kubernetes.io/docs/tasks/tools/
# - Kubernetes cluster is configured with kubectl and kubectl context is set to use this cluster
# - Port forwarding set up to MySQL pod using `kubectl -n mysql port-forward pod/mysql-0 <port>`

tmp_dir=$(mktemp -d)
trap "rm -rf ${tmp_dir}" EXIT

spire_bundle_file="${tmp_dir}/bundle.pem"
_goose="goose -certfile ${spire_bundle_file}"

fetch_bundle() {
    kubectl -n spire exec spire-server-0 -- /opt/spire/bin/spire-server bundle show > "${spire_bundle_file}"
}

setup_goose_config() {
    local mysql_pod_name="mysql-0"
    local mysql_address="127.0.0.1"
    local mysql_pod_forwarded_port="9999"
    local mysql_container_name="mysql"
    local mysql_user="root"
    local mysql_password=$(kubectl -n mysql logs "${mysql_pod_name}" "${mysql_container_name}" | grep "GENERATED ROOT PASSWORD" | sed 's/^.*GENERATED ROOT PASSWORD: \(.\+\)$/\1/g')
    local database_name="users"

    GOOSE_DBSTRING="${mysql_user}:${mysql_password}@tcp(${mysql_address}:${mysql_pod_forwarded_port})/${database_name}?parseTime=true"
    GOOSE_DRIVER=mysql
    export GOOSE_DRIVER GOOSE_DBSTRING

    fetch_bundle
}

setup_goose_config
$_goose -dir pkg/store/init up
$_goose -dir pkg/store/schema up
