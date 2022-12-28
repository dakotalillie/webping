#!/usr/bin/env bash

cd lambda/ping || exit
GOARCH=amd64 GOOS=linux go build -o webping .

cd ../../terraform || exit
terraform apply -var-file local.tfvars
