output "db_instance_id" {
  value = "${ucloud_db_instance.master.id}"
}

output "db_slave_id" {
  value = "${ucloud_db_slave.slave.id}"
}
