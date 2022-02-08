#!/bin/bash
set -e

ROOT=$(dirname $(readlink -f "$0"))

cd $ROOT
cd scripts/codegen/
cbgen -n server -v dev
mv *.go $ROOT/otto/server
cd $ROOT/otto/cmd/client
cbgen -n main -v dev
cd $ROOT/otto/cmd/server
EXE_NAME=".otto_dev"
go build -o $EXE_NAME
mv $EXE_NAME $ROOT
cd $ROOT
./$EXE_NAME --no-scheduler --static-dir $(realpath frontend/build) "$@"
