variable "aws_region" {
  description = "AWS region to create the infrastructure in"
  type        = string
}

variable "email" {
  description = "The email address to notify when repeated failures are observed"
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
    condition     = contains(["dev", "prod"], var.environment)
    error_message = "Must be one of: (dev, prod)."
  }
}

variable "stack_name" {
  description = "The name of the stack"
  type        = string
}
