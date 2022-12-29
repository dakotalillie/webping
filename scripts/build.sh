#!/usr/bin/env bash

cd lambda/ping || exit 1
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o webping .
