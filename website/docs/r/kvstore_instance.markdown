---
layout: "ucloud"
page_title: "UCloud: ucloud_kvstore_instance"
sidebar_current: "docs-ucloud-resource-kvstore-instance"
description: |-
  Provides an KVStore instance resource.
---

# ucloud_kvstore_instance

The UCloud Key-Value storage instance is an online storage service which is compatiable with Redis and Memcache protocol.

~> **Note** The Memcache applies to in-memory cache, and doesn't apply to data persistence storage when there is a downtime switching or storage extension taking place, since all the data will be wiped out .

## Example Usage

```hcl
data "ucloud_zones" "default" {}

data "ucloud_kvstore_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  engine_version    = "4.0"
}

resource "ucloud_kvstore_instance" "master" {
  name              = "tf-example-redis-master"
  tag               = "tf-example"
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  instance_type     = "redis.master.1"
  password          = "2018_KVStore"

  # engine argument must same as the engine of instance type argument
  engine         = "redis"
  engine_version = "4.0"

  # kvstore can set customized parameter group as redis config for active-standby redis
  parameter_group_id = "${data.ucloud_kvstore_parameter_groups.default.parameter_groups.0.id}"

  # here is the begin time of daily backup progress
  backup_begin_time = 3
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where instance is located. such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `name` - (Required) The name of the Key-Value storage instance which contains 6 to 63 characters, including letter，number，strikethrough and underline.
* `engine` - (Required) The type of engine, possible values are: "memcache" and "redis".
* `instance_type` - (Required) The type of Key-Value storage instance which specifies the emgine and storage size, please visit the [instance type table](/docs/providers/ucloud/appendix/kvstore_instance_type.html.markdown) for more details.
* `instance_charge_type` - (Optional) Charge type. Possible values are: "Year" as pay by year, "Month" as pay by month, "Dynamic" as pay by hour (specific permission required). The default value is "Month".
* `instance_duration` - (Optional) The duration that you will buy the resource, the default value is "1". It is not required when "Dynamic" (pay by hour), the value is "0" when pay by month and the instance will be vaild till the last day of that month.
* `tag` - (Optional) A tag to assign to the instance. The default value is "Default" (means no tag assigned).
* `vpc_id` - (Optional) The ID of VPC linked to the instance.
* `subnet_id` - (Optional) The ID of subnet linked to the instance.

The unique arguments for the active-standby KV storage for Redis ("Required" stands for the mandatory required for active-standby Redis and optional for all the other KV stroage instance type):

* `engine_version` - (Required) The version of engine of active-standby Redis, epossible values are: 3.0, 3.2 and 4.0.
* `password` - (Optional) The password for active-standby Redis instance which contains 6 to 36 characters, including letter，number，strikethrough and underline.
* `parameter_group_id` - (Optional) The ID of active-standby Redis parameter group, the specific parameter group will be applied if the ID is specified, otherwise the default parameter group will be applied.
* `backup_begin_time` - Specifies what time to start the auto backup, range from 0 to 23, format in UTC-8.
* `backup_id` - (Optional) The ID of active-standby Redis backup set, The instance is created based on a backup set if the ID is specified, otherwise the ID is set to "null".

~> **Note** The active-standby Redis doesn't support to be created on multiple zones with Terraform.

~> **Note** You have to enable the "UFile" service of required region before enabling the backup policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ip_set` - ip_set is a nested type. ip_set documented below.
* `create_time` - The time of creation for kvstore instance, in RFC3339 time string.
* `expire_time` - The expiration time for kvstore instance, in RFC3339 time string.
* `update_time` - The time whenever there is a change made to kvstore instance, in RFC3339 time string.
* `status` - The status of KV storage instance.

The attribute (`ip_set`) support the following:

* `ip` - The virtual ip of KV storage instance.
* `port` - The port on which Key-Value storage instance accepts connections, it is 6379 by default.
