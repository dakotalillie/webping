#!/usr/bin/env bash

stack="$1"
if [ "$stack" != "dev" ] && [ "$stack" != "test" ] && [ "$stack" != "prod" ] ; then
  echo "usage: build.sh (dev, test, prod)" >&2
  exit 1
fi

mkdir -p ./bin
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/ping ./cmd/ping

cd ./bin || exit 1
zip ping.zip ping

aws s3 cp ping.zip "s3://dakotalillie-lambda-src/webping/${stack}/ping.zip"
