#!/bin/bash

set -e
set -x
export PACKAGE_NAME="github.com/juandiegopalomino/cloudgrep"
export GOLANG_CROSS_VERSION="v1.18.3"

DIR="./bin"
mkdir -p $DIR
ls -lh

rm -f $DIR/*.zip

docker build --platform linux/amd64 -f ./Dockerfile-build --tag cloudgrep-build:latest .

docker run \
		--rm \
		--platform linux/amd64 \
		-e GO_VERSION=1.18 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/${PACKAGE_NAME} \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/${PACKAGE_NAME} \
		cloudgrep-build:latest \
		build --rm-dist --skip-validate "$@"

ls -lh $DIR
find $DIR -user root -exec sudo chown $USER: {} +
for directory in $(find $DIR -mindepth 1 -maxdepth 1 -type d)
do
  echo "Archiving ${directory}..."
  relative_name=`basename ${directory}`
  if [[ $relative_name == windows* ]]; then
    fin=$directory/cloudgrep.exe
  else
    fin=$directory/cloudgrep
  fi
  fout=$DIR/${relative_name}.zip
  zip -9 -q -j $fout $fin
done