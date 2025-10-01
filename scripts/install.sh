#!/bin/sh
INSTALL_DIR="./temp"
echo Installing to $INSTALL_DIR

./scripts/build.sh
# mv the bin file to /usr/bin
# web pages go to /var/www/html/
# config goes to /usr/local/etc/

mkdir -p $INSTALL_DIR/usr/bin/pi-server
cp ./bin/app $INSTALL_DIR/usr/bin/pi-server


mkdir -p $INSTALL_DIR/var/www/pi-server/static
cp -r ./web/static $INSTALL_DIR/var/www/pi-server/

mkdir -p $INSTALL_DIR/etc/pi-server
cp config.json $INSTALL_DIR/etc/pi-server
echo Done

exit