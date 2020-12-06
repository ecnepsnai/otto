#!/bin/bash
set -e

OTTO_PATH=$(realpath ../)
SCRIPTS_PATH=$(realpath .)
STATIC_DIR=$(realpath ./../static)
COLOR_NC='\033[0m'
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'
OTTO_VERSION=${1:-dev}

LOG=${OTTO_PATH}/otto-install.log

echo -en "Building frontend... "
cd ${STATIC_DIR}
npm install >> "${LOG}" 2>&1
npx webpack --config webpack.app.development.js >> "${LOG}" 2>&1
npx webpack --config webpack.login.development.js >> "${LOG}" 2>&1
cd ../
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
