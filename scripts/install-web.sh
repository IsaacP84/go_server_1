#!/bin/sh

if [ -n "${INSTALL_DIR+set}"]; then
    INSTALL_DIR="./temp"
fi
echo Installing web data to $INSTALL_DIR/var/www/pi-server/

mkdir -p $INSTALL_DIR/var/www/pi-server/static
cp -r ./web/static $INSTALL_DIR/var/www/pi-server/