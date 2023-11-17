#!/bin/sh
set -e

if [[ -z "${VERSION}" ]]; then
    VERSION=${1:?Version required}
fi
export VERSION=${VERSION}
BUILD_DATE=$(date -R)
BUILD_REVISION=$(git rev-parse HEAD)

ROOT_PATH=$(realpath ../)
OTTO_PATH=$(realpath ../otto)

CPU_ARCH=$(uname -m)
ARCH="amd64"
if [[ ${CPU_ARCH} == 'aarch64' ]]; then
    ARCH="arm64"
fi

PRODUCT_NAME=otto
PACKAGE_NAME=${PRODUCT_NAME}-${VERSION}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${ROOT_PATH}/logs/otto-install.log

cd ${ROOT_PATH}
mkdir -p ${ROOT_PATH}/logs
mkdir -p ${ROOT_PATH}/artifacts

echo -en "Building agent... "
cd ${OTTO_PATH}/
go build -v ./... >> ${LOG} 2>&1
go test -v ./... >> ${LOG} 2>&1

cd ${OTTO_PATH}/cmd/agent
CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} GOAMD64=v2 go build -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}' -X 'main.BuildRevision=${BUILD_REVISION}'" -trimpath -buildmode=exe -o ${PRODUCT_NAME} >> ${LOG} 2>&1
NAME=${PRODUCT_NAME}agent-${VERSION}_linux-${ARCH}
tar -czf ${NAME}.tar.gz otto
mv ${NAME}.tar.gz ${ROOT_PATH}/artifacts/
git clean -qxdf
cd ${ROOT_PATH}
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"

cd ${ROOT_PATH}/scripts/agent_rpm
./build.sh ${VERSION}
cd ${ROOT_PATH}/scripts/agent_deb
./build.sh ${VERSION}
