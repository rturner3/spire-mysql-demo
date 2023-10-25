#/bin/bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source "${SCRIPT_DIR}/common.sh"

bb=$(tput bold)
nn=$(tput sgr0)

list_entries_resp=$(spire_server entry show -output json)
entry_ids=$(echo "${list_entries_resp}" | jq -r '.entries[].id')
for entry_id in ${entry_ids}; do
    spire_server entry delete -entryID "${entry_id}"
done
