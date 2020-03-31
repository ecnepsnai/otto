#!/bin/bash
set -e

VERSION=${1:?Version required}
DOCKER_CMD=${DOCKER:-"docker"}

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
${DOCKER_CMD} build -t otto:${VERSION} .
cd ../
rm -rf Docker
${DOCKER_CMD} save otto:${VERSION} > otto-${VERSION}_docker.tar
gzip otto-${VERSION}_docker.tar
mv otto-${VERSION}_docker.tar.gz ../artifacts
