#!/bin/bash
set -e
set -x

pwd
dpkg-deb --build --root-owner-group ottoagent
mv ./ottoagent.deb ./ottoagent/otto-agent-${OTTO_VERSION}.${ARCH}.deb
