---
layout: "ucloud"
page_title: "UCloud: ucloud_db_backups"
sidebar_current: "docs-ucloud-datasource-db-backups"
description: |-
  Provides a list of DB backup resources in the current region.
---

# ucloud_db_backups

This data source providers a list of DB backup resources according to their db backup id, availability zone, db instance ID, class type, backup type, begin time, end time.

~> **Note** The "availablity zone" is mandatory required when querying backups via db backup ids.
## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where db backups are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist). It is mandatory required when querying backups via db backup ids.
* `ids` - (Optional) The group of IDs of db backups that require to be retrieved.
* `db_instance_id` - (Optional) The ID of database instance, the corresponding backup will be retrieved upon the specific argumnet defined.
* `backup_type` - (Optional) The type of backup, Possible values are "automatic" as automated backups and "manual" as manual backups.
* `begin_time` - (Optional) The time when start the backup, formatted by RFC3339 time string.
* `end_time` - (Optional) The time when finsih the backup, formatted by RFC3339 time string.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `total_count` - Total number of database instance that satisfy the condition.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backups` - db backups is a nested type. backups documented below.
* `total_count` - Total number of backup that satisfy the condition.

The attribute (`backups`) support the following:

* `id` - The ID of db backup.
* `zone` - Availability zone where db backups are located.
* `standby_zone` -  Availability zone where the standby database instance is located for the high availability database instance with multiple zone.
* `name` - The name of the backups.
* `size` - The size of backup, measured in byte.
* `type` - The type of backup, Possible values are "automatic" as automated backups and "manual" as manual backups.
* `db_instance_id` - The ID of DB instance.
* `db_instance_name` - The name of the DB instance.
* `begin_time` - The time when start the backup, formatted by RFC3339 time string.
* `end_time` - The time when finsih the backup, formatted by RFC3339 time string.
