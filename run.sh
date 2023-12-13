#!/bin/bash
set -e

ROOT=$(dirname $(readlink -f "$0"))

cd $ROOT
gengo -n server -c scripts/gengo -g otto/server -t frontend/src/types -q
cd otto/cmd/server
EXE_NAME=".otto_dev"
go build -o $EXE_NAME
mv $EXE_NAME $ROOT
cd $ROOT
./$EXE_NAME --no-scheduler --static-dir $(realpath frontend/build) "$@"
