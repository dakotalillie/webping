#!/usr/bin/env bash

stack="$1"
if [ "$stack" != "dev" ] && [ "$stack" != "prod" ] ; then
  echo "usage: deploy.sh (dev, prod)" >&2
  exit 1
fi

echo "deploying $stack stack"

./scripts/build.sh

cd terraform || exit 1

if ! terraform workspace list | grep "$stack" > /dev/null; then
  terraform workspace new "$stack"
else
  terraform workspace select "$stack"
fi

terraform apply -var-file "$stack.tfvars"
