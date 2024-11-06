#!/bin/bash

set -e

HOST=$1
KEY=$2
FROM=$3

if [ -z $3 ]
then
    echo "Usage: $0 <postal-host> <postal-api-key> <sender-address>"
    exit 1
fi

JSON_DATA=$(jq -n \
  --arg from "$FROM" \
  --argjson to '["mail@example.invalid"]' \
  --arg sender "$FROM" \
  --arg subject 'Test mail' \
  --arg plain_body 'Test mail' \
  '$ARGS.named'
)

echo "[*] Sending this payload:"
echo $JSON_DATA | jq

echo "[*] Sending post request:"
curl -X POST -H "X-Server-API-Key: ${KEY}" -H "Content-Type: application/json" -d "${JSON_DATA}" "${HOST}/api/v1/send/message" | jq

echo
