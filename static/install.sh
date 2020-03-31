#!/bin/bash
set -e

BOOTSTRAP_VERSION="4.3.1"
FONTAWESOME_VERSION="5.10.2"
JQUERY_VERSION="3.4.1"
ANGULARJS_VERSION="1.7.8"

COLOR_NC='\033[0m'
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'

mkdir -p copy/assets/css
mkdir -p copy/assets/js

# Bootstrap
echo -e "${COLOR_BLUE}[INFO]${COLOR_NC} Building CSS"
BOOSTRAP_URL="https://github.com/twbs/bootstrap/releases/download/v${BOOTSTRAP_VERSION}/bootstrap-${BOOTSTRAP_VERSION}-dist.zip"
BOOTSTRAP_ZIP="bootstrap.zip"
BOOTSTRAP_DIR="bootstrap"
curl -SsL ${BOOSTRAP_URL} > ${BOOTSTRAP_ZIP}
unzip -q ${BOOTSTRAP_ZIP} -d ${BOOTSTRAP_DIR}
rm ${BOOTSTRAP_ZIP}
cp bootstrap/bootstrap-${BOOTSTRAP_VERSION}-dist/css/bootstrap.min.* copy/assets/css
cp bootstrap/bootstrap-${BOOTSTRAP_VERSION}-dist/js/bootstrap.bundle.min.* copy/assets/js
rm -rf bootstrap

# Fontawesome
echo -e "${COLOR_BLUE}[INFO]${COLOR_NC} Building fonts"
mkdir -p fonts/
curl -SsL "https://use.fontawesome.com/releases/v${FONTAWESOME_VERSION}/fontawesome-free-${FONTAWESOME_VERSION}-web.zip" > fontawesome.zip
unzip fontawesome.zip >/dev/null
rm -rf package
mv fontawesome-free-${FONTAWESOME_VERSION}-web package
rm fontawesome.zip
sed -e 's/..\/webfonts/..\/fonts/g' package/css/all.css > fonts/fontawesome.css
cp package/webfonts/* fonts/
rm -rf package

#Javascript
echo -e "${COLOR_BLUE}[INFO]${COLOR_NC} Building JS Assets"
CDNJS_URL="https://cdnjs.cloudflare.com/ajax/libs"
ANGULARJS_URL="${CDNJS_URL}/angular.js/${ANGULARJS_VERSION}"
cd copy/assets/js

curl -Ss -O "${ANGULARJS_URL}/angular-route.min.js"
curl -Ss -O "${ANGULARJS_URL}/angular-route.min.js.map"
curl -Ss -O "${CDNJS_URL}/angular-sanitize/${ANGULARJS_VERSION}/angular-sanitize.min.js"
curl -Ss -O "${CDNJS_URL}/angular-sanitize/${ANGULARJS_VERSION}/angular-sanitize.min.js.map"
curl -Ss -O "${ANGULARJS_URL}/angular.min.js"
curl -Ss -O "${ANGULARJS_URL}/angular.min.js.map"
curl -Ss -O "${CDNJS_URL}/angular-moment/1.2.0/angular-moment.min.js"
curl -Ss -O "${CDNJS_URL}/moment.js/2.22.0/moment.min.js"
curl -Ss -O "${CDNJS_URL}/jquery/${JQUERY_VERSION}/jquery.min.js"
curl -Ss -O "${CDNJS_URL}/jquery/${JQUERY_VERSION}/jquery.min.map"

cd ../

