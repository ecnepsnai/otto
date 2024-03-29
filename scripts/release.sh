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
FRONTEND_PATH=$(realpath ../frontend)

PRODUCT_NAME=otto
PACKAGE_NAME=${PRODUCT_NAME}-${VERSION}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${ROOT_PATH}/logs/otto-install.log

cd ${ROOT_PATH}
git clean -qxdf
cd ${ROOT_PATH}/scripts
./install_backend.sh
cd ${FRONTEND_PATH}
rm -rf build/
echo -en "Building frontend... "
npm install >> ${LOG} 2>&1
node start_webpack.js --mode production >> ${LOG} 2>&1
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
cd ${ROOT_PATH}
rm -rf artifacts/
mkdir -p artifacts/

function build_server() {
    cd ${OTTO_PATH}/cmd/server
    CGO_ENABLED=0 GOOS=${1} GOARCH=${2} GOAMD64=v2 go build -ldflags="-s -w -X 'github.com/ecnepsnai/otto/server.Version=${VERSION}' -X 'github.com/ecnepsnai/otto/server.BuildDate=${BUILD_DATE}' -X 'github.com/ecnepsnai/otto/server.BuildRevision=${BUILD_REVISION}'" -trimpath -buildmode=exe -o ${PRODUCT_NAME} >> ${LOG} 2>&1
    
    cp -r ${FRONTEND_PATH}/build static
    mkdir agents
    cp ${ROOT_PATH}/artifacts/ottoagent* agents
    cp ${ROOT_PATH}/artifacts/otto-agent* agents

    NAME=${PRODUCT_NAME}-${VERSION}_${1}_${2}
    tar -czf ${NAME}.tar.gz ${PRODUCT_NAME} static agents
    mv ${NAME}.tar.gz ${ROOT_PATH}/artifacts/
    git clean -qxdf
    cd ${ROOT_PATH}
}

function build_agent() {
    cd ${OTTO_PATH}/cmd/agent
    CGO_ENABLED=0 GOOS=${1} GOARCH=${2} GOAMD64=v2 go build -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}' -X 'main.BuildRevision=${BUILD_REVISION}'" -trimpath -buildmode=exe -o ${PRODUCT_NAME} >> ${LOG} 2>&1
    NAME=${PRODUCT_NAME}agent-${VERSION}_${1}-${2}
    tar -czf ${NAME}.tar.gz otto
    mv ${NAME}.tar.gz ${ROOT_PATH}/artifacts/
    git clean -qxdf
    cd ${ROOT_PATH}
}

echo -en "Packaging agent builds... "
for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd'; do
        build_agent ${OS} ${ARCH}
    done
done
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"

cd ${ROOT_PATH}/scripts/agent_rpm
./build.sh ${VERSION}
cd ${ROOT_PATH}/scripts/agent_deb
./build.sh ${VERSION}

echo -en "Packaging server build... "
for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd'; do
        build_server ${OS} ${ARCH}
    done
done

echo -e "${COLOR_GREEN}Finished${COLOR_NC}"

cd ${ROOT_PATH}/scripts/
./docker.sh ${VERSION}
