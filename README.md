# Webping Infra

This repo contains Terraform code defining the AWS infrastructure needed for 
my Webping project, which alerts me if any of my websites go down. It also 
contains tests for this Terraform, written using
[Terratest](https://terratest.gruntwork.io/).

## Getting Started

1. Clone this repo.
2. Copy `.env.sample` to `.env`. The requisite values can be found via the 
   Twilio UI.
3. Copy `terraform/dev.tfvars.sample` to `terraform/dev.tfvars`. If you want 
   to test the email integration, specify a value for the `email` variable.
4. Run `scripts/deploy.sh dev` to create the resources in AWS.

## Contributing

Tests can be run by running `go test -v ./tests`. These tests will spin up 
actual infrastructure in AWS and send an SMS. Note that the SMS lambda 
relies upon some SSM parameter values which are not controlled via this 
project and were added manually via the AWS console.

Deployment to prod is handled by merging to the `main` branch, whereupon a 
GitHub Actions workflow will kick off to run any final checks and deploy any 
updates.

When you're done developing, `cd` into `terraform` directory, switch to the 
dev workspace using `terraform workspace select dev`, and run `terraform 
destroy`. None of the resources provisioned using this project should cost 
much money, but it's still good to clean up after oneself.

Source code for the Lambda functions is located in the
[webping-lambda repo](https://github.com/dakotalillie/webping-lambda).
