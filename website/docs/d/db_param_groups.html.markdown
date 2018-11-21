---
layout: "ucloud"
page_title: "UCloud: ucloud_db_param_groups"
sidebar_current: "docs-ucloud-datasource-db-param-groups"
description: |-
  Provides a list of DB param group resources in the current region.
---

# ucloud_db_param_groups

This data source providers a list of DB param groups resources according to their availability zone, DB param group ID, class type, engine, engine version, region flag.

~> **使用限制** 当通过 DB param group id来查询配置参数时，availability zone 参数必填; 对于单可用区配置文件ParamGroupId只在当前可用区唯一，多可用区配置文件在当前region唯一。
~> **Note** The "availability zone" is mandatory required when querying parameter group via calling "ParamGroupId"；the "ParamGroupId" is 
unique in the current availability zone for single availability zone param groups, and in the current region for cross availability zone param groups.
## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where db param groups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)/summary/regionlist)当通过id来查询备份时，此参数必填
* `availability_zone` - (Optional) Availability zone where db param groups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)/summary/regionlist), this is mandatory required when quering parameter group via calling "ParamGroupId".
* `ids` - (Optional) The group of IDs of db param groups that require to be retrieved.
* `engine` - Database type, possible values are: "mysql", "percona", "postgresql".
* `engine_version` - The database engine version, possible values are: "5.1", "5.5", "5.6", "5.7", "9.4", "9.6", "10.4".
* `region_flag` - (Optional) 当请求没有填写Zone时，如果指定为true，表示只拉取跨可用区的相关配置文件，否则，拉取所有机房的配置文件（包括每个单可用区和跨可用区）
* `region_flag` - (Optional) Possible values are " true" and "false", when "availability_zone" is not included in the request, only the cross-availability zone DB parameter groups wil be returned if this is "true" , otherwise all the DB parameter groups will be returned (including single availability zone and cross-availability zone).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `param_groups` - db param groups is a nested type. param groups documented below.
* `total_count` - Total number of param groups that satisfy the condition.

The attribute (`param_groups`) support the following:

* `id` - The ID of db param group.
* `zone` - Availability zone where db param groups are located.
* `name` - The name of db param group
* `engine` - Database type.
* `engine_version` - The database engine version.
* `description` - The description of db param group
* `modifiable` - 参数组是否可修改
* `modifiable` - To determine whether the DB parameter group is modifiable or not.
* `param_member` - param member is a nested type. param groups documented below.

The attribute (`param_member`) support the following:

* `key` - The key of param.
* `value` - The value of param.
* `value_type` - 参数值应用类型，取值范围为unknown, int, string, bool
* `value_type` - The type of param value, Possible values are "unknown", "int", "string" and "bool".
* `allowed_value` - value允许的值(根据参数类型，用分隔符表示).
* `allowed_value` - The valid param value.
