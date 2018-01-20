#!/bin/bash

echo "Generating Docker Compose Env file"
sed -e "s;%ID%;$1;g" -e "s;%KEY%;$2;g" \
    -e "s;%ENDPOINT%;$3;g" -e "s;%REGION%;$4;g" -e "s;%DSN%;$5;g" \
    env.template > .env
tail .env
