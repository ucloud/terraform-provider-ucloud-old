---
layout: "ucloud"
page_title: "UCloud: ucloud_kvstore_backups"
sidebar_current: "docs-ucloud-datasource-kvstore-backups"
description: |-
  Provides a list of available redis backup in the current region.
---

# ucloud_kvstore_backups

Use this data source to get information about backup for master and slave in-memory caching for Redis, the automated backups and manual backups are both included. The data will be restored back to the specific time when the backup was taken. The backups can be also used in archival storage.

~> **Note** The current backup is only valid for master and slave in-memory caching for Redis.

## Example Usage

```hcl
data "ucloud_kvstore_backups" "example" {
}

output "first" {
    value = "${data.ucloud_kvstore_backups.example.backups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where in-memory caching backups are located. You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `ids` - (Optional) The ID of in-memory caching backups.
* `name_regex` - (Optional) A regex string to filter resulting backups by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `kvstore_instance_id` - (Optional) The ID of active-standby in-memory caching for Redis.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backups` - backups is a nested type. backups documented below.
* `total_count` - Total number of backup that satisfy the condition.

The attribute (`backups`) support the following:

* `availability_zone` - Availability zone where backups are located.
* `id` - The ID of backup.
* `name` - The name of backup.
* `kvstore_instance_id` - The ID of active-standby in-memory caching for Redis.
* `kvstore_instance_name` - The name of active-standby in-memory caching for Redis.
* `size` - The size of backup, measured in byte.
* `type` - The type of backup to be returned, Possible values are "Manual" as manual backups and "Auto" as automated backups.
* `status` - The status of backup, possible values are "Backuping", "Success", "Error" and "Expired".
* `create_time` - The time of creation for backup, in RFC3339 time string.
