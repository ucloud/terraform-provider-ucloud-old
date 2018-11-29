# Specify the provider and access details
provider "ucloud" {
  region = "${var.region}"
}

data "ucloud_zones" "default" {}

data "ucloud_kvstore_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  engine_version    = "4.0"
}

resource "ucloud_kvstore_instance" "master" {
  name              = "tf-example-redis-master"
  tag               = "tf-example"
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  instance_type     = "${var.instance_type}"
  password          = "${var.password}"

  # engine argument must same as the engine of instance type argument
  engine         = "redis"
  engine_version = "4.0"

  # kvstore can set customized parameter group as redis config for active-standby redis
  parameter_group_id = "${data.ucloud_kvstore_parameter_groups.default.parameter_groups.0.id}"

  # here is the begin time of daily backup progress
  backup_begin_time = 3
}

resource "ucloud_kvstore_slave" "slave" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  name              = "tf-example-redis-read-only-slave-${count.index + 1}"
  tag               = "tf-example"
  instance_type     = "${var.instance_type}"

  # current instance is the slave of master has the id
  master_id = "${ucloud_kvstore_instance.master.id}"

  # kvstore can set customized parameter group as redis config for active-standby redis
  parameter_group_id = "${data.ucloud_kvstore_parameter_groups.default.parameter_groups.0.id}"

  # easily scale out by count
  count = "${var.slave_count}"
}
