#!/bin/bash

SCRIPT_DIR=`dirname $0`

VERSION=`cat ${SCRIPT_DIR}/.version`

echo ${VERSION}


if ! test -f "kc_rslbe-$VERSION-amd64.tar.bz2"; then
    echo "kc_rslbe-$VERSION-amd64.tar.bz2 can not be found. Aborting ..."
    exit 1
fi

echo "Deploying kc_rslbe-$VERSION-amd64.tar.bz2 locally ..."


`sudo rm -rf  /data/tools/repository/micro-services/kc_rslbe-${VERSION}-amd64/bin`
echo "sudo rm -rf  /data/tools/repository/micro-services/kc_rslbe-${VERSION}-amd64/bin"


`sudo tar -jxf kc_rslbe-${VERSION}-amd64.tar.bz2  -C /data/tools/repository/micro-services/`
echo "sudo tar -jxf kc_rslbe-$VERSION-amd64.tar.bz2  -C /data/tools/repository/micro-services/"

`sudo chown -R cumulis:cumulis /data/tools/repository/micro-services/kc_rslbe-${VERSION}/`
echo "sudo chown -R cumulis:cumulis /data/tools/repository/micro-services/kc_rslbe-$VERSION-amd64/"

`sudo unlink /data/tools/repository/micro-services/kc_rslbe/bin`
echo "sudo unlink /data/tools/repository/micro-services/kc_rslbe/bin"

`sudo ln -s /data/tools/repository/micro-services/kc_rslbe-${VERSION}/bin/ /data/tools/repository/micro-services/kc_rslbe`
echo "sudo ln -s /data/tools/repository/micro-services/kc_rslbe-${VERSION}/bin/ /data/tools/repository/micro-services/kc_rslbe"

echo "Done. kc_rslbe-$VERSION-amd64.tar.bz2 has been deployed locally"
