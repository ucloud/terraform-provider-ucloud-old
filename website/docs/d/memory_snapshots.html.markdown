---
layout: "ucloud"
page_title: "UCloud: ucloud_memory_snapshots"
sidebar_current: "docs-ucloud-datasource-memory-snapshots"
description: |-
  Provides a list of available redis snapshot in the current region.
---

# ucloud_memory_snapshots

此 Data Source 用于查询出主备版 Redis 的备份列表，包括由备份策略自动创建的备份，以及手动创建的备份。
备份可用于恢复 Redis 中的数据到该备份创建的时间，以及用于归档历史数据。
Use this data source to get information about snapshot for master and slave in-memory caching for Redis, the automated snapshots and manual snapshots are both included. The data will be restored back to the specific time when the snapshot was taken. The snapthosts can be also used in archival storage.

~> **使用限制** 注意，Redis 备份只能用于主备版 Redis，不能用于其它内存实例。
~> **Note** The current snapshot is only valid for master and slave in-memory caching for Redis.

## Example Usage

```hcl
data "ucloud_snapshot" "example" {
}

output "first" {
    value = "${data.ucloud_memory_snapshots.example.snapshots.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where in-memory caching snapshots are located. You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `ids` - (Optional) The ID of in-memory caching  snapshots.
* `name_regex` - (Optional) A regex string to filter resulting snapshots by name. Such as: .
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `memory_instance_id` - (Optional) 主备版 Redis 实例的 ID
* `memory_instance_id` - (Optional) The ID of master and slave in-memory caching for Redis.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshots` - snapshots is a nested type. snapshots documented below.
* `total_count` - Total number of snapshot that satisfy the condition.

The attribute (`snapshots`) support the following:

* `availability_zone` - 可用区
* `availability_zone` - Availability zone where snapshots are located.
* `id` - The ID of snapshot.
* `name` - The name of snapshot.
* `memory_instance_id` - 主备版 Redis 实例的 ID.
* `memory_instance_id` - The ID of master and slave in-memory caching for Redis.
* `memory_instance_name` - 主备版 Redis 实例的名称.
* `memory_instance_name` - The name of master and slave in-memory caching for Redis.
* `size` - 备份文件大小, 以字节为单位
* `size` - The size of snapshot, measured in byte.
* `type` - 备份类型: Manual 手动 Auto 自动
* `type` - The type of snapshots to be returned, Possible values are "Manual" as manual snapshots and "Auto" as automated snapshots.
* `status` - The status of snapshot, possible values are "Backuping", "Success", "Error" and "Expired".
* `create_time` - 备份时间 (UNIX时间戳)
* `create_time` - Provides the time when the snapshot was taken, in UNIX timestamp.
