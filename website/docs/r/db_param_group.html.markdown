---
layout: "ucloud"
page_title: "UCloud: ucloud_db_param_group"
sidebar_current: "docs-ucloud-resource-db-param-group"
description: |-
  Provides a DB param group resource.
---

# ucloud_db_param_group

Provides a DB param group resource. 
~> Note The "availability zone" is mandatory required when querying parameter group via calling "ParamGroupId"；the "ParamGroupId" is unique in the current availability zone for single availability zone param groups, and it is also unique in the current region for multiple zones param groups.
## Example Usage

```hcl
# Query availability zone
data "ucloud_zones" "default" {}

# Query parameter group
data "ucloud_db_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  region_flag       = "false"
  engine            = "mysql"
  engine_version    = "5.7"
}

# Create parameter group
resource "ucloud_db_parameter_group" "example" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  name              = "tf-example-parameter-group"
  description       = "this is a test"
  src_group_id      = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
  engine            = "mysql"
  engine_version    = "5.7"

  parameter_input {
    key   = "max_connections"
    value = "3000"
  }

  parameter_input {
    key   = "slow_query_log"
    value = "1"
  }
}
```
## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where db param groups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist),The "availability zone" is mandatory required when querying parameter group via calling "ParamGroupId".
* `name` - (Required) The name of db param group.
* `description` - (Optional) The description of db param group.
* `src_group_id` - (Required) The ID of source DB param group.
* `engine` - (Required) Database type, possible values are: "mysql", "percona", "postgresql".
* `engine_version` - (Required) The database engine version, possible values are: "5.1", "5.5", "5.6", "5.7", "9.4", "9.6", "10.4".
* `region_flag` - (Optional) Possible values are " true" and "false", when "availability_zone" is not included in the request, only the multiple zones DB parameter groups wil be returned if this is "true" , otherwise all the DB parameter groups will be returned (including single availability zone and multiple zones).
* `parameter_input` - (Optional) parameter input is a nested type. parameter input documented below.

The attribute (`parameter_input`) support the following:
* `key` - (Required) The key of param.
* `value` - (Required) The value of param.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `parameter_output` - db param groups is a nested type. param groups documented below.

The attribute (`parameter_output`) support the following:

* `key` - (Required) The key of param.
* `value` - (Required) The value of param.
