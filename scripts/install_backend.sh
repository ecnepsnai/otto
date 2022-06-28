#!/bin/bash
set -e

ROOT_PATH=$(realpath ../)
LOGS_PATH=${ROOT_PATH}/logs
OTTO_PATH=$(realpath ../otto)
SCRIPTS_PATH=$(realpath .)
COLOR_NC='\033[0m'
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'

LOG=${LOGS_PATH}/otto-install.log
mkdir -p ${LOGS_PATH}

echo -en "Building backend... "
cd ${SCRIPTS_PATH}/codegen/
cbgen -n server
mv *.go ${OTTO_PATH}/server
cd ${OTTO_PATH}/cmd/agent
cbgen -n main
cd ${OTTO_PATH}/
go get -t >> ${LOG} 2>&1
cd ${OTTO_PATH}/cmd/agent
go get -t >> ${LOG} 2>&1
cd ${OTTO_PATH}/cmd/server
go get -t >> ${LOG} 2>&1
cd ${OTTO_PATH}/server
go get -t >> ${LOG} 2>&1
cd ${OTTO_PATH}/
go build
go test -v >> ${LOG} 2>&1
cd ${OTTO_PATH}/server
go test -v >> ${LOG} 2>&1
cd ${OTTO_PATH}/
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
