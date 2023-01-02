/* REQUIRED */

variable "function_name" {
  description = "The name of the Lambda function"
  type        = string
}

variable "handler" {
  description = "The name of the handler for the Lambda function"
  type        = string
}

variable "iam_role_policy" {
  description = "The policy to be attached to the Lambda function's IAM role"
  type        = string
}

variable "s3_key" {
  description = "The key of the S3 object containing the Lambda's source code"
  type        = string
}

/* OPTIONAL */

variable "environment_variables" {
  default     = null
  description = "Environment variables for the Lambda function"
  type        = map(string)
}
