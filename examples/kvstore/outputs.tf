output "master_id" {
  value = "${ucloud_kvstore_instance.master.id}"
}

output "slave_id_list" {
  value = ["${ucloud_kvstore_slave.slave.*.id}"]
}
