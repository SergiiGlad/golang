#!/bin/bash

echo "Generating Docker Compose Env file"
cat <<End-of-env-template  > .env
ID=$1
KEY=$2
ENDPOINT=$3
REGION=$4
DSN=$5
End-of-env-template
cat .env
