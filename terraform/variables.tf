variable "aws_region" {
  description = "AWS region to create the infrastructure in"
  type        = string
}

variable "email" {
  default     = ""
  description = "The email address to notify when repeated failures are observed. An empty value means no email will be subscribed."
  type        = string
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

variable "stack_name" {
  description = "The name of the stack"
  type        = string
}
