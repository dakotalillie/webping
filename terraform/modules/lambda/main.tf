resource "aws_iam_role" "this" {
  name = "${var.function_name}-lambda"

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

resource "aws_iam_role_policy" "inline" {
  name   = "InlinePolicy"
  policy = var.iam_role_policy
  role   = aws_iam_role.this.id
}

resource "aws_iam_role_policy" "logging" {
  name = "LoggingPolicy"
  role = aws_iam_role.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["logs:CreateLogStream", "logs:PutLogEvents"]
        Effect   = "Allow"
        Resource = aws_cloudwatch_log_group.this.arn
      }
    ]
  })
}

resource "aws_lambda_function" "this" {
  depends_on = [aws_cloudwatch_log_group.this]

  function_name = var.function_name
  handler       = var.handler
  role          = aws_iam_role.this.arn
  runtime       = "go1.x"
  s3_bucket     = "dakotalillie-lambda-src"
  s3_key        = var.s3_key

  dynamic "environment" {
    for_each = var.environment_variables != null ? toset([1]) : toset([])
    content {
      variables = var.environment_variables
    }
  }
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${var.function_name}"
  retention_in_days = 7
}
