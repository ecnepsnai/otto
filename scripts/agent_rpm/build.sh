#!/bin/bash
set -e

ROOT_PATH=$(realpath ../../)
OTTO_VERSION=${1:?Version required}
BUILD_DATE=$(date -R)
DOCKER_CMD=${DOCKER:-"podman"}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${ROOT_PATH}/logs/otto-install.log

echo -en "Building agent rpm... "

rm -rf build/
mkdir -p build/
rm -rf otto-agent-${OTTO_VERSION}/
mkdir otto-agent-${OTTO_VERSION}/

cp -r ../../otto otto-agent-${OTTO_VERSION}/
cp otto-agent.service otto-agent-${OTTO_VERSION}/
tar -czf otto-agent-${OTTO_VERSION}.tar.gz otto-agent-${OTTO_VERSION}/
cp Dockerfile otto-agent.spec entrypoint.sh build/
rm -rf otto-agent-${OTTO_VERSION}/
mv otto-agent-${OTTO_VERSION}.tar.gz build/

cd build/

GOLANG_ARCH="amd64"
if [[ $(uname -m) == 'aarch64' ]]; then
    GOLANG_ARCH="arm64"
fi
GOLANG_VERSION=$(curl -sS "https://go.dev/dl/?mode=json" | jq -r '.[0].version' | sed 's/go//')

podman build --build-arg GOLANG_ARCH=${GOLANG_ARCH} --build-arg GOLANG_VERSION=${GOLANG_VERSION} -t otto_build_rpm:${OTTO_VERSION} . >> ${LOG} 2>&1
rm -rf rpms
mkdir -p rpms
podman run --user root -v $(readlink -f rpms):/root/rpmbuild/RPMS:Z -e OTTO_VERSION=${OTTO_VERSION} -e BUILD_DATE="${BUILD_DATE}" -it otto_build_rpm:${OTTO_VERSION} >> ${LOG} 2>&1
cp rpms/*/*.rpm .
mv *.rpm ../../../artifacts
cd ../
rm -rf build/

echo -e "${COLOR_GREEN}Finished${COLOR_NC}"