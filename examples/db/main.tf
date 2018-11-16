# Specify the provider and access details
provider "ucloud" {
  region = "${var.region}"
}

resource "ucloud_db_instance" "master" {
  availability_zone    = "cn-bj2-03"
  name                 = ""
  engine               = "mongodb"
  engine_version       = "3.2"
  username             = ""
  password             = ""
  port                 = 3306
  cpu                  = 1
  memory               = 1
  instance_type        = "primary"
  instance_duration    = 1
  instance_charge_type = "Month"
  instance_storage     = 20

  # vpc
  vpc_id    = ""
  subnet_id = ""

  # param group
  param_group_id = ""

  # MySQL master/slave
  is_slave = false
}

resource "ucloud_db_instance" "slave" {
  engine         = "mysql"
  engine_version = "5.7"

  availability_zone    = "cn-bj2-03"
  name                 = ""
  engine               = "mongodb"
  engine_version       = "3.2"
  username             = ""
  password             = ""
  port                 = 3306
  cpu                  = 1
  memory               = 1
  instance_type        = "primary"
  instance_duration    = 1
  instance_charge_type = "Month"
  instance_storage     = 20

  # vpc
  vpc_id    = ""
  subnet_id = ""

  is_slave  = true
  master_id = "${ucloud_db_instance.master.id}"
}
