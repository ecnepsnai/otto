#!/bin/bash
set -e

OTTO_PATH=$(realpath ../)
VERSION=${1:?Version required}
DOCKER_CMD=${DOCKER:-"podman"}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${OTTO_PATH}/otto-install.log


echo -en "Packaging server container... "
rm -rf Docker
mkdir Docker
cp Dockerfile Docker/
cp entrypoint.sh Docker/entrypoint.sh
cd Docker

# Add service
cp ../../artifacts/otto-${VERSION}_linux_amd64.tar.gz .
tar -xzf otto-${VERSION}_linux_amd64.tar.gz
rm otto-${VERSION}_linux_amd64.tar.gz
mv otto-${VERSION} otto
${DOCKER_CMD} build -t otto:${VERSION} -t ghcr.io/ecnepsnai/otto:${VERSION} -t otto:latest -t ghcr.io/ecnepsnai/otto:latest . >> ${LOG} 2>&1
cd ../
rm -rf Docker
${DOCKER_CMD} save --quiet -o otto-${VERSION}_docker_amd64.tar otto:${VERSION} >> ${LOG} 2>&1
gzip otto-${VERSION}_docker_amd64.tar
mv otto-${VERSION}_docker_amd64.tar.gz ../artifacts
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
