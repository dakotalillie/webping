#!/usr/bin/env bash

cd lambda/ping || exit 1
go test -v .
