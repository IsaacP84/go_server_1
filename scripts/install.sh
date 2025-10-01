#!/bin/sh
INSTALL_DIR="./temp"
echo Installing to $INSTALL_DIR

./scripts/build.sh

mkdir -p $INSTALL_DIR/usr/bin/pi-server
cp ./bin/app $INSTALL_DIR/usr/bin/pi-server

./scripts/install-web.sh

mkdir -p $INSTALL_DIR/etc/pi-server
cp config.json $INSTALL_DIR/etc/pi-server
echo Done

exit