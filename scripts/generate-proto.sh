#!/bin/bash
DIR=$(dirname "$(dirname "$0")")
protoc --go_out="$DIR/pkg" --go-grpc_out=$DIR/pkg $DIR/pkg/api/*.proto
