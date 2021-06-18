#!/bin/bash
set -e

cd ../../cmd/client
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o otto
mv otto ../../scripts/client
cd ../../scripts/client
podman build -t ottoclient:latest --squash .
rm -f otto