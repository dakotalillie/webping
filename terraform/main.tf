terraform {
  required_version = "~> 1.3"

  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.48"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      ManagedBy   = "Terraform"
      Project     = "Webping"
      Source      = "github.com/dakotalillie/webping"
      Stack       = var.stack_name
    }
  }
}

locals {
  lambda_function_name = "webping-${var.stack_name}"
}

resource "aws_iam_role" "ping_lambda" {
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
  name                = "webping-${var.stack_name}-lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "ping_lambda" {
  name = "inline-policy"
  role = aws_iam_role.ping_lambda.id

  policy = jsonencode({
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
      }
    ]
  })
}

data "archive_file" "ping_source" {
  type        = "zip"
  source_file = "../lambda/ping/webping"
  output_path = "webping.zip"
}

resource "aws_lambda_function" "ping" {
  depends_on       = [aws_cloudwatch_log_group.ping]
  filename         = data.archive_file.ping_source.output_path
  function_name    = local.lambda_function_name
  handler          = "webping"
  role             = aws_iam_role.ping_lambda.arn
  runtime          = "go1.x"
  source_code_hash = data.archive_file.ping_source.output_base64sha256

  environment {
    variables = {
      DB_TABLE  = aws_dynamodb_table.ping_state.name
      ENDPOINTS = join(",", var.endpoints)
      SNS_TOPIC = aws_sns_topic.ping_notification.arn
    }
  }
}

resource "aws_cloudwatch_log_group" "ping" {
  name              = "/aws/lambda/${local.lambda_function_name}"
  retention_in_days = 7
}

resource "aws_sns_topic" "ping_notification" {
  name = "webping-${var.stack_name}"
}

resource "aws_sns_topic_subscription" "email" {
  endpoint  = var.email
  protocol  = "email"
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

#resource "aws_lambda_permission" "allow_cloudwatch" {
#  action        = "lambda:InvokeFunction"
#  function_name = aws_lambda_function.ping.function_name
#  principal     = "events.amazonaws.com"
#}

#resource "aws_cloudwatch_event_rule" "ping_cron" {}
#
#resource "aws_cloudwatch_event_target" "ping_cron" {
#  arn  = ""
#  rule = aws_cloudwatch_event_rule.ping_cron.name
#}
