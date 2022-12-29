output "ping_lambda_function_name" {
  value = aws_lambda_function.ping.function_name
}

output "sns_topic_arn" {
  value = aws_sns_topic.ping_notification.arn
}
