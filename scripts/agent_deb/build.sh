#!/bin/bash
set -e

ROOT_PATH=$(realpath ../../)
OTTO_VERSION=${1:?Version required}
BUILD_DATE=$(date -R)
DOCKER_CMD=${DOCKER:-"podman"}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${ROOT_PATH}/logs/otto-install.log
ARCH="amd64"
if [[ $(uname -m) == 'aarch64' ]]; then
    ARCH="arm64"
fi

echo -en "Building agent deb... "

rm -rf build
mkdir -p build/DEBIAN
mkdir -p build/opt/otto-agent
mkdir -p build/usr/lib/systemd/system

cp ../../artifacts/ottoagent-${OTTO_VERSION}_linux-${ARCH}.tar.gz .
tar -xzf ottoagent-${OTTO_VERSION}_linux-${ARCH}.tar.gz
rm ottoagent-${OTTO_VERSION}_linux-${ARCH}.tar.gz
mv otto build/opt/otto-agent/agent

cp otto-agent.control.spec build/DEBIAN/control
perl -pi -e "s,%%VERSION%%,${OTTO_VERSION},g" build/DEBIAN/control
cp prerm.sh build/DEBIAN/prerm
cp postinst.sh build/DEBIAN/postinst
chmod +x build/DEBIAN/prerm
chmod +x build/DEBIAN/postinst

cp otto-agent.service build/usr/lib/systemd/system

podman build -t otto_build_deb:${OTTO_VERSION} . >> ${LOG} 2>&1
podman run --user root -v $(readlink -f build):/ottoagent:Z -e OTTO_VERSION=${OTTO_VERSION} -e ARCH=${ARCH} -it otto_build_deb:${OTTO_VERSION} >> ${LOG} 2>&1

cp build/*.deb ../../artifacts
rm -rf build

echo -e "${COLOR_GREEN}Finished${COLOR_NC}"