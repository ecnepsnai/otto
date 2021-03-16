#!/bin/sh
set -e

if [[ -z "${VERSION}" ]]; then
    VERSION=${1:?Version required}
fi
export VERSION=${VERSION}

OTTO_PATH=$(realpath ../)

PRODUCT_NAME=otto
PACKAGE_NAME=otto-${VERSION}
COLOR_NC='\033[0m'
COLOR_GREEN='\033[0;32m'
LOG=${OTTO_PATH}/otto-install.log

cd ${OTTO_PATH}
git clean -qxdf
cd ${OTTO_PATH}/scripts
./install_backend.sh ${VERSION}
cd ${OTTO_PATH}/static/
rm -rf build/
echo -en "Building frontend... "
npm install >> ${LOG} 2>&1
npx webpack --config webpack.login.production.js >> ${LOG}
npx webpack --config webpack.app.production.js >> ${LOG}
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
cd ${OTTO_PATH}
rm -rf artifacts/
mkdir -p artifacts/

function build_server() {
    cd ${OTTO_PATH}/cmd/server
    CGO_ENABLED=0 GOOS=${1} GOARCH=${2} go build -ldflags="-s -w" -o ${3}
    NAME=${PRODUCT_NAME}-${VERSION}_${1}_${2}
    mv ${3} ../../
    cd ../../

    rm -rf ${PACKAGE_NAME}
    mkdir -p ${PACKAGE_NAME}/static
    mkdir -p ${PACKAGE_NAME}/clients
    mv ${3} ${PACKAGE_NAME}
    cp -r static/build ${PACKAGE_NAME}/static
    cp -r artifacts/ottoclient* ${PACKAGE_NAME}/clients
    tar -czf ${NAME}.tar.gz ${PACKAGE_NAME}/
    rm -rf ${PACKAGE_NAME}/
    mv ${NAME}.tar.gz artifacts/
}

function build_client() {
    cd ${OTTO_PATH}/cmd/client
    CGO_ENABLED=0 GOOS=${1} GOARCH=${2} go build -ldflags="-s -w"
    NAME=${PRODUCT_NAME}client-${VERSION}_${1}-${2}
    mv client ../../otto
    cd ../../
    mkdir ${PACKAGE_NAME}
    mv otto ${PACKAGE_NAME}
    cd ${PACKAGE_NAME}
    tar -czf ${NAME}.tar.gz *
    mv ${NAME}.tar.gz ../artifacts/
    cd ../
    rm -rf ${PACKAGE_NAME}
}

echo -en "Packaging client builds... "
for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd'; do
        build_client ${OS} ${ARCH}
    done
done
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"

echo -en "Packaging server build... "
for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd'; do
        build_server ${OS} ${ARCH} ${PRODUCT_NAME}
    done
done
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"

cd scripts/
./docker.sh ${VERSION}
