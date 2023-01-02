output "ping_lambda_function_name" {
  value = module.ping_lambda_function.function_name
}

output "sms_lambda_function_name" {
  value = module.sms_lambda_function.function_name
}

output "sns_topic_arn" {
  value = aws_sns_topic.ping_notification.arn
}
