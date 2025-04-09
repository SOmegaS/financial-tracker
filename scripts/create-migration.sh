#!/bin/bash
DIR=$(dirname "$(dirname "$0")")
migrate create -ext="sql" -dir="$DIR/internal/database/migrations" -seq "$1"
