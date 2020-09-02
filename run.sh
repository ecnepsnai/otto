#!/bin/bash
set -e

ROOT=$(pwd)

cd $ROOT
cd scripts/codegen/
cbgen -n server -v dev
mv *.go $ROOT/server
cd $ROOT/cmd/client
cbgen -n main -v dev
cd $ROOT/cmd/server
EXE_NAME="otto_$(uname)_$(uname -m)"
go build -o $EXE_NAME
mv $EXE_NAME $ROOT
cd $ROOT
./$EXE_NAME --no-scheduler "$@"
