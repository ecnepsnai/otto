#!/bin/sh
set -e

if [[ -z "${VERSION}" ]]; then
    VERSION=${1:?Version required}
fi
export VERSION=${VERSION}

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
./install_backend.sh ${VERSION}
cd ${FRONTEND_PATH}
rm -rf build/
echo -en "Building frontend... "
npm install >> ${LOG} 2>&1
npx webpack --config webpack.login.production.js >> ${LOG}
npx webpack --config webpack.app.production.js >> ${LOG}
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"
cd ${ROOT_PATH}
rm -rf artifacts/
mkdir -p artifacts/

function build_server() {
    cd ${OTTO_PATH}/cmd/server
    CGO_ENABLED=0 GOOS=${1} GOARCH=${2} go build -ldflags="-s -w" -trimpath -buildmode=exe -o ${PRODUCT_NAME} >> ${LOG} 2>&1
    
    cp -r ${FRONTEND_PATH}/build static
    mkdir clients
    cp ${ROOT_PATH}/artifacts/ottoclient* clients

    NAME=${PRODUCT_NAME}-${VERSION}_${1}_${2}
    tar -czf ${NAME}.tar.gz ${PRODUCT_NAME} static clients
    mv ${NAME}.tar.gz ${ROOT_PATH}/artifacts/
    git clean -qxdf
    cd ${ROOT_PATH}
}

function build_client() {
    cd ${OTTO_PATH}/cmd/client
    CGO_ENABLED=0 GOOS=${1} GOARCH=${2} go build -ldflags="-s -w" -trimpath -buildmode=exe -o ${PRODUCT_NAME} >> ${LOG} 2>&1
    NAME=${PRODUCT_NAME}client-${VERSION}_${1}-${2}
    tar -czf ${NAME}.tar.gz otto
    mv ${NAME}.tar.gz ${ROOT_PATH}/artifacts/
    git clean -qxdf
    cd ${ROOT_PATH}
}

echo -en "Packaging client builds... "
for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd'; do
        build_client ${OS} ${ARCH}
    done
done
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"

cd ${ROOT_PATH}/scripts/client_rpm
./build.sh ${VERSION}

echo -en "Packaging server build... "
for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd'; do
        build_server ${OS} ${ARCH}
    done
done

cd ${OTTO_PATH}
git checkout -- cmd/client/cbgen_version.go server/cbgen_version.go || true
echo -e "${COLOR_GREEN}Finished${COLOR_NC}"

cd ${ROOT_PATH}/scripts/
./docker.sh ${VERSION}
