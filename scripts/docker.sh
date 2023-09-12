#!/bin/bash
set -e

CPU_ARCH=$(uname -m)
ARCH="amd64"
ALPINE_HASH="c5c5fda71656f28e49ac9c5416b3643eaa6a108a8093151d6d1afc9463be8e33" # https://hub.docker.com/_/alpine/tags
if [[ ${CPU_ARCH} == 'aarch64' ]]; then
    ARCH="arm64"
    ALPINE_HASH="b312e4b0e2c665d634602411fcb7c2699ba748c36f59324457bc17de485f36f6"
fi

ROOT_PATH=$(realpath ../)
VERSION=${1:?Version required}
DOCKER_CMD=${DOCKER:-"podman"}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${ROOT_PATH}/logs/otto-install.log
REVISION=$(git rev-parse HEAD)
DATETIME=$(date --rfc-3339=seconds)

echo -en "Packaging server container... "
rm -rf Docker
mkdir Docker
cp Dockerfile Docker/Dockerfile
cd Docker

# Add service
mkdir otto
cd otto
cp ${ROOT_PATH}/artifacts/otto-${VERSION}_linux_${ARCH}.tar.gz .
tar -xzf otto-${VERSION}_linux_${ARCH}.tar.gz
rm otto-${VERSION}_linux_${ARCH}.tar.gz
cd ../
${DOCKER_CMD} build \
    --squash \
    --no-cache \
    --build-arg "ALPINE_HASH=${ALPINE_HASH}" \
    --label "org.opencontainers.image.created=${DATETIME}" \
    --label "org.opencontainers.image.version=${VERSION}" \
    --label "org.opencontainers.image.revision=${REVISION}" \
    -t otto:${VERSION} \
    -t otto:latest \
    -t ghcr.io/ecnepsnai/otto:${VERSION} \
    -t ghcr.io/ecnepsnai/otto:latest \
    . >> ${LOG} 2>&1
cd ../
rm -rf Docker
${DOCKER_CMD} save --quiet -o otto-${VERSION}_docker_${ARCH}.tar otto:${VERSION} >> ${LOG} 2>&1
gzip otto-${VERSION}_docker_${ARCH}.tar
mv otto-${VERSION}_docker_${ARCH}.tar.gz ${ROOT_PATH}/artifacts
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
