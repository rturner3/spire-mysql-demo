#!/bin/bash
#
# Sets up schemas in MySQL database for sample service.
# Prerequisites:
# - kubectl is installed and available on the PATH: https://kubernetes.io/docs/tasks/tools/
# - mysql CLI is installed and available on the PATH
# - Kubernetes cluster is configured with kubectl and kubectl context is set to use this cluster
# - Port forwarding set up to MySQL pod using `kubectl -n mysql port-forward pod/mysql-0 <port>`

tmp_dir=$(mktemp -d)
trap "rm -rf ${tmp_dir}" EXIT

spire_bundle_file="${tmp_dir}/bundle.pem"

fetch_bundle() {
    kubectl -n spire exec spire-server-0 -- /opt/spire/bin/spire-server bundle show > "${spire_bundle_file}"
}

mysql_pod_name="mysql-0"
mysql_container_name="mysql"
mysql_password=$(kubectl -n mysql logs "${mysql_pod_name}" "${mysql_container_name}" | grep "GENERATED ROOT PASSWORD" | sed 's/^.*GENERATED ROOT PASSWORD: \(.\+\)$/\1/g')

_mysql() {
    local mysql_user="root"
    local mysql_address="127.0.0.1"
    local mysql_pod_forwarded_port="9999"
    mysql --user "${mysql_user}" \
        --password="${mysql_password}" \
        --host "${mysql_address}" \
        --port "${mysql_pod_forwarded_port}" \
        --ssl-ca "${spire_bundle_file}" \
        "$@"
}

fetch_bundle
for sql_file in ./pkg/store/init/*.sql; do
    _mysql < "${sql_file}"
done

for sql_file in ./pkg/store/schema/*.sql; do
    _mysql < "${sql_file}"
done
