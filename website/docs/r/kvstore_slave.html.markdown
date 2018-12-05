---
layout: "ucloud"
page_title: "UCloud: ucloud_kvstore_slave"
sidebar_current: "docs-ucloud-resource-kvstore-slave"
description: |-
  Provides an KVStore slave resource.
---

# ucloud_kvstore_slave

The active-standby Redis with read only feature is supported by the slave set of UCloud KV storage instance, in order to offering the exceptional read capability.

## Example Usage

```hcl
data "ucloud_zones" "default" {}

data "ucloud_kvstore_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  engine_version    = "4.0"
}

resource "ucloud_kvstore_slave" "slave" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  name              = "tf-example-redis-read-only-slave-${count.index + 1}"
  tag               = "tf-example"
  instance_type     = "redis.master.1"

  # current instance is the slave of master has the id
  master_id = "..."

  # kvstore can set customized parameter group as redis config for active-standby redis
  parameter_group_id = "${data.ucloud_kvstore_parameter_groups.default.parameter_groups.0.id}"

  # easily scale out by count
  count = 1
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where instance is located. such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `name` - (Required) The name of KV storage instance which contains 6 to 63 characters, including letter，number，strikethrough and underline.
* `instance_type` - (Required) The type of Key-Value storage instance which specifies the emgine and storage size, please visit the [instance type table](/docs/providers/ucloud/appendix/kvstore_instance_type.html.markdown) for more details.
* `master_id` - (Required) The ID of master set of Redis.
* `password` - (Optional) The password of Redis which contains 6 to 36 characters, including letter，number，strikethrough and underline.
* `parameter_group_id` - (Optional) The ID of parameter group of Redis, the specific parameter group applies if the ID is specified, otherwise the default parameter group applies.
* `tag` - (Optional) A tag to assign to the instance. The default value is "Default" (means no tag assigned).

~> **Note** The active-standby Redis doesn't support to be created on multiple zones with Terraform.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ip_set` - ip_set is a nested type. ip_set documented below.
* `create_time` - The time of creation for kvstore instance, in RFC3339 time string.
* `expire_time` - The expiration time for kvstore instance, in RFC3339 time string.
* `update_time` - The time whenever there is a change made to kvstore instance, in RFC3339 time string.
* `status` - The status of KV storage instance.
* `vpc_id` - The ID of VPC linked to the instance.
* `subnet_id` - The ID of subnet linked to the instance.
* `instance_charge_type` - The charge type of instance, possible values are: "Year", "Month" and "Dynamic" as pay by hour.

The attribute (`ip_set`) support the following:

* `ip` - The virtual ip of KV storage instance.
* `port` - The port on which Key-Value storage instance accepts connections, it is 6379 by default.
