---
layout: "ucloud"
page_title: "UCloud: ucloud_db_slave"
sidebar_current: "docs-ucloud-resource-db-slave"
description: |-
  Provides a Database slave resource.
---

# ucloud_db_slave

Provides a Database slave resource.

~> **Note** Please try to keep the same settings for both master and slave databases, otherwise it's likely to have an issue of data loss when making synchronization. The slave database creation is not supported for mysql and percona in 5.5 version. In addition, the slave is a basic database (normal version), it takes around 5 mins to shut down when making upgrade/degrade(incloud the memory of instance_type and instance_storage), please make the necessary arrangements to your business in advance to prevent any loss of data. Up to five slave databases can be created from the same master database.
## Example Usage

```hcl
# Query availability zone
data "ucloud_zones" "default" {}

# Create parameter group
data "ucloud_db_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  region_flag       = "false"
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
  password           = "2018_dbSlave"
  parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
  tag                = "tf-example"

  # Backup policy
  backup_begin_time = 4
  backup_count      = 6
  backup_date       = "0111110"
  backup_black_list = ["test.%"]
}

# Create database slave
resource "ucloud_db_slave" "slave" {
  availability_zone  = "${data.ucloud_zones.default.zones.0.id}"
  name               = "f-example-db-slave"
  instance_storage   = 20
  instance_type      = "mysql-basic-1"
  password           = "2018_dbSlave"
  parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
  master_id          = "${ucloud_db_instance.master.id}"
}
```

## Argument Reference

The following arguments are supported:

* `master_id` - (Required) The ID of master DB instance, it is mandatory required to request when creating database slave.
* `name` - (Required) The name of the database slave, should have 1 - 63 characters and only support chinese, english, numbers, '-', '_', '.'.
* `parameter_group_id` - (Required) The ID of database parameter group. Note: The "parameter_group_id" of the multiple zone database should be included in the request for the slave of high availability database instance with multiple zone. When it is changed, the database slave will reboot to make the change take effect.
* `instance_storage` - (Optional) Specifies the allocated storage size in gigabytes (GB), range from 20 to 3000GB. The volume adjustment must be a multiple of 10 GB. The maximum disk volume for SSD type are：
    - 500GB if the memory chosen is equal or less than 8GB;
    - 1000GB if the memory chosen is from 12 to 24GB;
    - 2000GB if the memory chosen is 32GB;
    - 3000GB if the memory chosen is equal or more than 48GB.
* `instance_type` - (Optional) Specifies the type of database slave with format "engine-type-memory", Possible values are:
    - "mysql","percona" and "postgresql" for engine;
    - "basic" is the only one type as normal version.
    - possible values for memory are: 1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64GB.
* `is_lock` - (Optional) Specifies whether need to set master DB to read only when creating database slave, possible values are "true" and "false", it is "false" by default.
* `password` - (Optional) The password for the database slave which should have 8-30 characters. It must contain at least 3 items of Capital letters, small letter, numbers and special characters. The special characters include <code>`()~!@#$%^&*-+=_|{}\[]:;'<>,.?/</code>.
* `port` - (Optional) The port on which the DB accepts connections, the default port is 3306 for mysql and percona and 5432 for postgresql.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone` - Availability zone where database slaves are located.
* `status` - Specifies the status of database slave, possible values are: "Init","Fail", "Starting", "Running", "Shutdown", "shutoff", "Delete", "Upgrading", "Promoting", "Recovering" and "Recover fail".
* `create_time` - The creation time of database slave, formatted by RFC3339 time string.
* `expire_time` - The expiration time of database slave, formatted by RFC3339 time string.
* `modify_time` - The modification time of database slave, formatted by RFC3339 time string.
* `instance_charge_type` - The charge type of slave, possible values are: "Year", "Month" and "Dynamic" as pay by hour.
* `vpc_id` - The ID of VPC linked to the slave.
* `subnet_id` - The ID of subnet linked to the slave.
