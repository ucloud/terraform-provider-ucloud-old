---
layout: "ucloud"
page_title: "UCloud: ucloud_db_backups"
sidebar_current: "docs-ucloud-datasource-db-backups"
description: |-
  Provides a list of DB backup resources in the current region.
---

# ucloud_db_backups

This data source providers a list of DB backup resources according to their db backup id, availability zone, DB instance ID, class type, backup type, begin time, end time.

~> **使用限制** 当通过 DB backup id来查询备份时，availability zone 参数必填

## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where db backups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)当通过id来查询备份时，此参数必填
* `ids` - (Optional) The group of IDs of db backups that require to be retrieved.
* `db_instance_id` - (Optional) DB实例Id，如果指定，则只获取该db的备份信息 该值可以通过DescribeUDBInstance获取
* `class_type` - (Optional) DB种类，分为sql和postgresql,其中，sql代表mysql和percona,
* `backup_type` - (Optional) 备份类型,取值为0或1，0表示自动，1表示手动
* `begin_time` - (Optional) 过滤条件:起始时间
* `end_time` - (Optional) 过滤条件:结束时间
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_backups` - db backups is a nested type. backups documented below.
* `total_count` - Total number of backup that satisfy the condition.

The attribute (`db_backups`) support the following:

* `id` - The ID of db backup.
* `zone` - Availability zone where db backups are located.
* `backup_zone` - 跨可用区高可用备库所在可用区
* `name` - 备份名称
* `backup_size` - 备份文件大小(字节)
* `backup_type` - 备份类型,取值为0或1，0表示自动，1表示手动
* `db_instance_id` - DB实例Id
* `db_instance_name` - 	对应的db实例名称
* `backup_begin_time` - 备份开始时间
* `backup_end_time` - 备份完成时间