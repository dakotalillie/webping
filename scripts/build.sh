#!/usr/bin/env bash

mkdir -p ./bin
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/ping ./cmd/ping
