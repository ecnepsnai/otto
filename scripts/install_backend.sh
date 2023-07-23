#!/bin/bash
set -e

ROOT_PATH=$(realpath ../)
LOGS_PATH=${ROOT_PATH}/logs
OTTO_BACKEND_PATH=$(realpath ../otto)
OTTO_FRONTEND_PATH=$(realpath ../frontend)
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
mv -v *.go ${OTTO_BACKEND_PATH}/server >> ${LOG} 2>&1
mv -v *.ts ${OTTO_FRONTEND_PATH}/src/types >> ${LOG} 2>&1
cd ${OTTO_BACKEND_PATH}/
go build -v ./... >> ${LOG} 2>&1
go test -v ./... >> ${LOG} 2>&1
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
