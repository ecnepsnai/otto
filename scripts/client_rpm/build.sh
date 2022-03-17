#!/bin/bash
set -e

ROOT_PATH=$(realpath ../../)
OTTO_VERSION=${1:?Version required}
DOCKER_CMD=${DOCKER:-"podman"}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${ROOT_PATH}/logs/otto-install.log

echo -en "Building client rpm... "

rm -rf build/
mkdir -p build/
rm -rf otto-${OTTO_VERSION}/
mkdir otto-${OTTO_VERSION}/

cp -r ../../otto otto-${OTTO_VERSION}/
cp otto.service otto-${OTTO_VERSION}/
tar -czf otto-${OTTO_VERSION}.tar.gz otto-${OTTO_VERSION}/
cp Dockerfile otto.spec build/
rm -rf otto-${OTTO_VERSION}/
mv otto-${OTTO_VERSION}.tar.gz build/

cd build/
perl -pi -e "s,##VERSION##,${OTTO_VERSION},g" otto.spec

GOLANG_ARCH="amd64"
if [[ $(uname -m) == 'aarch64' ]]; then
    GOLANG_ARCH="arm64"
fi

podman build --build-arg GOLANG_ARCH=${GOLANG_ARCH} -t otto_build . >> ${LOG} 2>&1
rm -rf rpms
mkdir -p rpms
podman run --user root -v $(readlink -f rpms):/root/rpmbuild/RPMS:Z -it otto_build >> ${LOG} 2>&1
cp rpms/*/*.rpm .
mv *.rpm $(ls *.rpm | sed 's/otto/ottoclient/')
mv *.rpm ../../../artifacts
cd ../
rm -rf build/

echo -e "${COLOR_GREEN}Finished${COLOR_NC}"