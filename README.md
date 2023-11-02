# MySQL mTLS Authentication using SPIRE

Demonstration of using SPIRE-issued X.509-SVIDs in MySQL database and clients, providing client authentication
and connection encryption in MySQL via mutual-TLS (mTLS), deployed on Kubernetes. This project assumes readers to have 
basic knowledge of the following technologies and concepts:

1. [SPIFFE](https://spiffe.io/) and [SPIRE](https://github.com/spiffe/spire)
2. [Kubernetes](https://kubernetes.io/)
3. [MySQL](https://www.mysql.com/)
4. [Transport Layer Security](https://en.wikipedia.org/wiki/Transport_Layer_Security)

## [Getting Started](#getting-started)

This section guides deploying and setting up components involved in this project on a Kubernetes cluster.
Refer [design](#design) for details on the components involved.

### Prerequisites

1. `kubectl` is installed and available on the PATH: https://kubernetes.io/docs/tasks/tools/ 
2. Kubernetes cluster is configured with `kubectl` and `kubectl context` is set to use this cluster
3. `curl` is installed and available on the PATH: https://everything.curl.dev/get

### Deploy and Setup

Deploy SPIRE server and SPIRE agent in the `spire` namespace on the cluster.
The script also creates necessary registrations in SPIRE server.
```
./01-deploy-spire.sh
```

Deploy MySQL server along with a persistent volume in the `mysql` namespace on the cluster
```
./02-deploy-mysql.sh
```

In a separate terminal, start port forwarding to MySQL server pod
```
kubectl -n mysql port-forward pod/mysql-0 3306:3306
```

Setup MySQL users and database, by connecting to the MySQL server via port-forwarding.
```
./03-setup-mysql.sh
```

Deploy `sample-service` in the default namespace.
```
./04-deploy-service.sh
```

### Verify DB Access using Sample Service

This section verifies DB access via mTLS using Sample Service.

In a separate terminal, start port forwarding to the `sample-service` 
```
kubectl port-forward service/sample-service  8888:8888
```

Use `curl` to verify the `GET /api/v1/users` endpoint on `localhost:8888`. The output should contain sample user data which was created
during MySQL DB setup.
```
curl -s -X GET http://localhost:8888/api/v1/users
```

Use `curl` to create a new user 
```
curl -s -X POST http://localhost:8888/api/v1/users -H 'Content-Type: application/json' -d '{"Name":"David"}'
```

Verify the newly created user is showing up in the `GET /api/v1/users` request
```
curl -s -X GET http://localhost:8888/api/v1/users
```

Verify the DB connection made by `sample-service` in the MySQL server general log. It should have a `Connect` log message from 
`spire-mysql-client` user (representing sample-service) connecting to the `spiredemo` DB using SSL/TLS.
```
kubectl exec -it mysql-0 -n mysql -c mysql -- cat /var/lib/mysql/general.log
```

### Cleanup 

Cleanup the environment using the cleanup script
```
./cleanup.sh
```

## [Design](#design)

MySQL server uses its X.509-SVIDs as server certificate and encryption whereas clients use the X.509-SVIDs to 
authenticate to the MySQL server. 

### Architecture

![Architecture](./docs/img/architecture.png)

SPIRE Server and SPIRE Agent are deployed on the Kubernetes Cluster along with registration entries created in the
SPIRE Server. MySQL is deployed with an init-container that fetches X.509-SVID for the MySQL server and 
writes to a tmpfs volume shared with the MySQL server container. Upon init, MySQL server container reads
the X.509-SVID from the tmpfs volume for SSL configuration. 

The sample-service also fetches its X.509 SVID from SPIRE and uses it to authenticate to MySQL server and  
to verify MySQL server's certificate (mTLS). The design also uses SPIRE CredentialComposer plugin to configure
X.509-SVIDs issued to sample-service have a unique Subject with Common Name - `/C=US/O=SPIRE/CN=spire-mysql-client`,
thereby mapping this service's identity to the user created in MySQL -
`CREATE USER 'spire-mysql-client' REQUIRE SUBJECT '/C=US/O=SPIRE/CN=spire-mysql-client'` 

### Certificate Auto-Rotation

![Rotation](./docs/img/rotation.png)

For security reasons, SPIRE server is configured to issue X.509-SVIDs with a relatively low TTL. This poses a challenge wrt 
MySQL server as it constantly requires updating the server SSL configuration with new X.509-SVID content from SPIRE 
agent. To achieve auto-rotation, MySQL server pod has a custom sidecar container called TLS reloader, 
which is responsible for rotating the MySQL server's SSL configuration. It does so by fetching X.509-SVID updates 
from SPIRE agent, writes them to shared tmpfs volume and executes `ALTER INSTANCE RELOAD TLS` query on the MySQL server, 
to force MySQL server to reload  the SSL configuration from disk.

