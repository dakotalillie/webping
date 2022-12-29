#!/usr/bin/env bash

cd terraform || exit 1
terraform workspace select dev
terraform apply -var-file local.tfvars
