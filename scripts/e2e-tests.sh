#!/usr/bin/env bash

cd e2e-tests || exit 1
go test -v .
