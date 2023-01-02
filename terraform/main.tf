terraform {
  required_version = "~> 1.3"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.48"
    }
  }

  backend "s3" {
    bucket               = "dakotalillie-tfstate"
    dynamodb_table       = "dakotalillie-tfstate-lock"
    key                  = "terraform.tfstate"
    region               = "us-west-1"
    workspace_key_prefix = "webping"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      ManagedBy   = "Terraform"
      Project     = "Webping"
      Source      = "github.com/dakotalillie/webping-infra"
      Stack       = var.stack_name
    }
  }
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name
}

module "ping_lambda_function" {
  source = "./modules/lambda"

  function_name = "webping-${var.stack_name}-ping"
  handler       = "ping"
  s3_key        = "webping/prod/ping.zip"

  environment_variables = {
    DB_TABLE  = aws_dynamodb_table.ping_state.name
    ENDPOINTS = join(",", var.endpoints)
    SNS_TOPIC = aws_sns_topic.ping_notification.arn
  }

  iam_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["dynamodb:BatchWriteItem", "dynamodb:Query"]
        Effect   = "Allow"
        Resource = aws_dynamodb_table.ping_state.arn
      },
      {
        Action   = "sns:Publish"
        Effect   = "Allow"
        Resource = aws_sns_topic.ping_notification.arn
      },
    ]
  })
}

module "sms_lambda_function" {
  source = "./modules/lambda"

  function_name = "webping-${var.stack_name}-sms"
  handler       = "sms"
  s3_key        = "webping/prod/sms.zip"

  iam_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = ["ssm:GetParameter", "ssm:GetParameters"]
        Effect = "Allow"
        Resource = [
          "arn:aws:ssm:${local.region}:${local.account_id}:parameter/Personal/PhoneNumber",
          "arn:aws:ssm:${local.region}:${local.account_id}:parameter/Twilio/AccountSID",
          "arn:aws:ssm:${local.region}:${local.account_id}:parameter/Twilio/AuthToken",
          "arn:aws:ssm:${local.region}:${local.account_id}:parameter/Twilio/PhoneNumber",
        ]
      }
    ]
  })
}

resource "aws_sns_topic" "ping_notification" {
  name = "webping-${var.stack_name}"
}

resource "aws_sns_topic_subscription" "email" {
  count = var.email == "" ? 0 : 1

  endpoint  = var.email
  protocol  = "email"
  topic_arn = aws_sns_topic.ping_notification.arn
}

resource "aws_sns_topic_subscription" "sms" {
  endpoint  = module.sms_lambda_function.function_arn
  protocol  = "lambda"
  topic_arn = aws_sns_topic.ping_notification.arn
}

resource "aws_dynamodb_table" "ping_state" {
  billing_mode   = "PROVISIONED"
  hash_key       = "Endpoint"
  name           = "webping-${var.stack_name}"
  range_key      = "Timestamp"
  read_capacity  = 1
  write_capacity = 1

  attribute {
    name = "Endpoint"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "N"
  }

  ttl {
    attribute_name = "ExpirationTime"
    enabled        = true
  }
}

resource "aws_cloudwatch_event_rule" "ping_cron" {
  count = var.enable_ping_cron ? 1 : 0

  name                = "webping-${var.stack_name}-cron"
  schedule_expression = var.ping_lambda_schedule_expression
}

resource "aws_cloudwatch_event_target" "ping_cron" {
  count = var.enable_ping_cron ? 1 : 0

  arn  = module.ping_lambda_function.function_arn
  rule = aws_cloudwatch_event_rule.ping_cron[0].name
}

resource "aws_lambda_permission" "allow_cloudwatch_to_trigger_ping" {
  count = var.enable_ping_cron ? 1 : 0

  action        = "lambda:InvokeFunction"
  function_name = module.ping_lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.ping_cron[0].arn
  statement_id  = "AllowExecutionFromEventBridge"
}

resource "aws_lambda_permission" "allow_sns_to_trigger_sms" {
  action        = "lambda:InvokeFunction"
  function_name = module.sms_lambda_function.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.ping_notification.arn
  statement_id  = "AllowExecutionFromSNS"
}
