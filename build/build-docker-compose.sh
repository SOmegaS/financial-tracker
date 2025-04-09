#!/bin/bash
DIR=$(dirname $(dirname "$0"))
docker-compose -f="$DIR/deployments/docker-compose.yaml" --env-file="$DIR/configs/.env" up
