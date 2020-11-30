output "instance_id" {
  value = aws_instance.example.id
}

output "ssm_parameter" {
  value = aws_ssm_parameter.instance_id.id
}
