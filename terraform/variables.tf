variable "aws_region" {
  description = "AWS region to create the infrastructure in"
  type        = string
}

variable "email" {
  default     = ""
  description = "The email address to notify when repeated failures are observed. An empty value means no email will be subscribed."
  type        = string
}

variable "enable_ping_cron" {
  default     = false
  description = "Whether or not to enable the cron for invoking the ping lambda"
  type        = bool
}

variable "endpoints" {
  description = "A list of endpoints to ping"
  type        = list(string)
}

variable "environment" {
  description = "The environment of the stack"
  type        = string

  validation {
    condition     = contains(["dev", "test", "prod"], var.environment)
    error_message = "Must be one of: (dev, test, prod)."
  }
}

variable "ping_lambda_schedule_expression" {
  default     = "rate(5 minutes)"
  description = "Schedule expression for triggering the ping lambda. See https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/ScheduledEvents.html"
  type        = string
}

variable "stack_name" {
  description = "The name of the stack"
  type        = string
}
