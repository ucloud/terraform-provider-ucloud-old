---
layout: "ucloud"
page_title: "UCloud: ucloud_kvstore_param_groups"
sidebar_current: "docs-ucloud-datasource-kvstore-parameter-groups"
description: |-
  Provides a list of KVStore parameter group resources in the current region.
---

# ucloud_kvstore_param_groups

This data source providers a list of KVStore parameter groups resources according to their availability zone, KVStore parameter group ID, class type, engine, engine version, region flag.

## Example Usage

```hcl
data "ucloud_kvstore_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  engine_version    = "4.0"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where kvstore parameter groups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist), this is mandatory required when quering parameter group via calling "ParamGroupId".
* `ids` - (Optional) The group of IDs of kvstore parameter groups that require to be retrieved.
* `name_regex` - (Optional) A regex string to filter resulting parameter groups by name. Such as: "^redis-[34].0$" means redis 3.0 or redis 4.0.
* `engine_version` - The kvstore engine version, possible values are: "3.0", "3.2", "4.0".
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `parameter_groups` - kvstore parameter groups is a nested type. parameter groups documented below.
* `total_count` - Total number of parameter groups that satisfy the condition.

The attribute (`parameter_groups`) support the following:

* `availability_zone` - Availability zone where kvstore parameter groups are located.
* `id` - The ID of kvstore parameter group.
* `name` - The name of kvstore parameter group
* `description` - The description of kvstore parameter group
* `engine_version` - The kvstore engine version.
* `status` - The status of parameter group
* `create_time` - The time of creation for parameter group, in RFC3339 time string.
* `update_time` - The time whenever there is a change made to parameter group, in RFC3339 time string.
