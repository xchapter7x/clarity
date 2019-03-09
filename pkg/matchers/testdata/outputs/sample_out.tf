output "rds_address" {
  value = "${element(concat(aws_db_instance.rds.*.address, list("")), 0)}"
}

output "rds_port" {
  value = "${element(concat(aws_db_instance.rds.*.port, list("")), 0)}"
}

output "rds_username" {
  value = "${element(concat(aws_db_instance.rds.*.username, list("")), 0)}"
}

output "rds_password" {
  sensitive = true
  value     = "${element(concat(aws_db_instance.rds.*.password, list("")), 0)}"
}

output "rds_db_name" {
  value = "${element(concat(aws_db_instance.rds.*.name, list("")), 0)}"
}
