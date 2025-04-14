#!/bin/bash
DIR=$(dirname "$(dirname "$0")")
sqlc generate -f="$DIR/configs/sqlc.yaml"
