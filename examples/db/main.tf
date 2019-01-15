# Specify the provider and access details
provider "ucloud" {
  region = "${var.region}"
}

# Query availability zone
data "ucloud_zones" "default" {}

# Create parameter group
data "ucloud_db_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  multi_az       = "false"
  engine            = "mysql"
  engine_version    = "5.7"
}

# Create database instance
resource "ucloud_db_instance" "master" {
  availability_zone  = "${data.ucloud_zones.default.zones.0.id}"
  name               = "tf-example-db-instance"
  instance_storage   = 20
  instance_type      = "mysql-ha-1"
  engine             = "mysql"
  engine_version     = "5.7"
  password           = "${var.db_password}"
  parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"

  # Backup policy
  backup_begin_time = 4
  backup_count      = 6
  backup_date       = "0111110"
  backup_black_list = "test.%"
}

# Create database slave
resource "ucloud_db_slave" "slave" {
  availability_zone  = "${data.ucloud_zones.default.zones.0.id}"
  name               = "f-example-db-slave"
  instance_storage   = 20
  instance_type      = "mysql-basic-1"
  password           = "${var.db_password}"
  parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
  master_id          = "${ucloud_db_instance.master.id}"
}
