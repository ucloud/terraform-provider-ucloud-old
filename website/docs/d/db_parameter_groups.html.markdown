---
layout: "ucloud"
page_title: "UCloud: ucloud_db_parameter_groups"
sidebar_current: "docs-ucloud-datasource-db-parameter-groups"
description: |-
  Provides a list of DB parameter group resources in the current region.
---

# ucloud_db_parameter_groups

This data source providers a list of DB parameter groups resources according to their availability zone, DB parameter group ID, engine, engine version, region flag.

~> **Note** The "availability zone" is mandatory required when querying parameter group via parameter group IDs; the parameter group ID is unique in the current availability zone for single availability zone parameter groups, and it is also unique in the current region for cross availability zone parameter groups.
## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where db parameter groups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist), this is mandatory required when quering parameter group via parameter group IDs.
* `ids` - (Optional) The group of IDs of db parameter groups that require to be retrieved.
* `engine` - Database type, possible values are: "mysql", "percona", "postgresql".
* `engine_version` - The database engine version, possible values are: "5.1", "5.5", "5.6", "5.7", "9.4", "9.6", "10.4".
* `multi_az` - (Optional) Possible values are " true" and "false"; when "availability_zone" is not included in the request, only the cross-availability zone DB parameter groups wil be returned if this is "true" , otherwise all the DB parameter groups will be returned (including single availability zone and cross-availability zone).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `parameter_groups` - DB parameter groups is a nested type. parameter groups documented below.
* `total_count` - Total number of parameter groups that satisfy the condition.

The attribute (`parameter_groups`) support the following:

* `id` - The ID of db parameter group.
* `name` - The name of db parameter group
* `engine` - Database type.
* `engine_version` - The database engine version.
* `description` - The description of db parameter group
* `modifiable` - To determine whether the DB parameter group is modifiable or not.
* `parameter_member` - parameter member is a nested type. parameter groups documented below.

The attribute (`parameter_member`) support the following:

* `key` - The key of parameter.
* `value` - The value of parameter.
* `value_type` - The type of parameter value, Possible values are "unknown", "int", "string" and "bool".
* `allowed_value` - The valid parameter value.
