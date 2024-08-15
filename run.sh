#!/bin/bash
set -e

ROOT=$(dirname $(readlink -f "$0"))

cd $ROOT
gengo -q -n server -c scripts/gengo -g otto/server -t frontend/src/types
cd $ROOT/otto/cmd/server
EXE_NAME=".otto_dev"
go build -ldflags="-X 'github.com/ecnepsnai/otto/server.DangerousSkipForcePasswordChangeForDefaultUser=true'" -o $EXE_NAME
mv $EXE_NAME $ROOT
cd $ROOT
./$EXE_NAME --no-scheduler --static-dir $(realpath frontend/build) "$@"
