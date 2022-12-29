#!/usr/bin/env bash

printenv

cd e2e-tests || exit 1
go test -v .
