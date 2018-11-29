output "vpc_id" {
  value = "${ucloud_vpc.default.id}"
}

output "subnet_id" {
  value = "${ucloud_subnet.default.id}"
}

output "security_group_id" {
  value = "${ucloud_security_group.default.id}"
}

output "image_id" {
  value = "${data.ucloud_images.default.images.0.id}"
}

output "instance_id_list" {
  value = "${ucloud_instance.web.*.id}"
}
