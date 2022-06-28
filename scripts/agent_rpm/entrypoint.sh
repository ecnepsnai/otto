#!/bin/bash

/usr/bin/rpmbuild -ba --define "_version ${OTTO_VERSION}" --define "_date ${BUILD_DATE}" otto.spec