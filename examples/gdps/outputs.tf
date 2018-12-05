output "primary_region" {
  value = "${var.primary_region}"
}

output "replica_region" {
  value = "${var.replica_region}"
}

output "primary_instance_id_list" {
  value = ["${module.app-primary.instance_id_list}"]
}

output "replica_instance_id_list" {
  value = ["${module.app-replication-cross-region.instance_id_list}"]
}
