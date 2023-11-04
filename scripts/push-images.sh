#!/bin/bash

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

"${SCRIPT_DIR}/push-spire-server-image.sh"
"${SCRIPT_DIR}/push-mysql-tls-images.sh"
"${SCRIPT_DIR}/push-sample-service-tls-image.sh"
