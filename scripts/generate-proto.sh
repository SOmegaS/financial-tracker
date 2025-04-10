#!/bin/bash
DIR=$(dirname "$(dirname "$0")")
protoc --go_out="$DIR/pkg/api" $DIR/pkg/api/*.proto
