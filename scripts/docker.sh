#!/bin/bash
set -e

CPU_ARCH=$(uname -m)
ARCH="amd64"
if [[ ${CPU_ARCH} == 'aarch64' ]]; then
    ARCH="arm64"
fi

OTTO_PATH=$(realpath ../)
VERSION=${1:?Version required}
DOCKER_CMD=${DOCKER:-"podman"}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${OTTO_PATH}/otto-install.log


echo -en "Packaging server container... "
rm -rf Docker
mkdir Docker
cp ${ARCH}.Dockerfile Docker/Dockerfile
cp entrypoint.sh Docker/entrypoint.sh
cd Docker

# Add service
cp ../../artifacts/otto-${VERSION}_linux_${ARCH}.tar.gz .
tar -xzf otto-${VERSION}_linux_${ARCH}.tar.gz
rm otto-${VERSION}_linux_${ARCH}.tar.gz
mv otto-${VERSION} otto
${DOCKER_CMD} build -t otto:${VERSION} -t ghcr.io/ecnepsnai/otto:${VERSION} -t otto:latest -t ghcr.io/ecnepsnai/otto:latest . >> ${LOG} 2>&1
cd ../
rm -rf Docker
${DOCKER_CMD} save --quiet -o otto-${VERSION}_docker_${ARCH}.tar otto:${VERSION} >> ${LOG} 2>&1
gzip otto-${VERSION}_docker_${ARCH}.tar
mv otto-${VERSION}_docker_${ARCH}.tar.gz ../artifacts
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
