#!/bin/bash
DIR=$(dirname "$(dirname "$0")")
$DIR/scripts/generate-models.sh
$DIR/scripts/generate-proto.sh
