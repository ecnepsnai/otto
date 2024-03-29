#!/bin/bash
set -e

ROOT=$(dirname $(readlink -f "$0"))

cd $ROOT
cd scripts/codegen/
cbgen -n server
mv *.go $ROOT/otto/server
mv *.ts $ROOT/frontend/src/types
cd $ROOT/otto/cmd/server
EXE_NAME=".otto_dev"
go build -o $EXE_NAME
mv $EXE_NAME $ROOT
cd $ROOT
./$EXE_NAME --no-scheduler --static-dir $(realpath frontend/build) "$@"
