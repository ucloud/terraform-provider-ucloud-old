---
layout: "ucloud"
page_title: "UCloud: ucloud_kvstore_slave"
sidebar_current: "docs-ucloud-resource-kvstore-slave"
description: |-
  Provides an KVStore slave resource.
---

# ucloud_kvstore_slave

UCloud KV 存储实例从库，提供了主备版 Redis 的只读实例用于提升读取性能。
The master-slave Redis with read only feature  is supported by the slave set of UCloud KV storage instance, in order to offering the exceptional read capability.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where instance is located. such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `name` - (Required) KV 存储实例名称，6-63 个字符，支持字母、数字、中划线或下划线。
* `name` - (Required) The name of KV storage instance which contains 6 to 63 characters, including letter，number，strikethrough and underline.
* `instance_type` - (Required) KV 存储实例的机型，表示实例引擎类型和容量的组合，可参考下方 [KV 存储实例机型表](/docs/providers/ucloud/appendix/kvstore_instance_type.html.markdown)。
* `instance_type` - (Required) The type of Key-Value storage instance which specifies the engine and storage size, please visit the instance type table for more details. (/docs/providers/ucloud/appendix/kvstore_instance_type.html.markdown).
* `master_id` - (Required) 主备版 Redis 主实例的 ID。
* `master_id` - (Required) The ID of master set of master-slave Redis.
* `password` - (Optional) 主备版 Redis 实例的密码，6-36 个字符，支持字母、数字、中划线或下划线。
* `password` - (Optional) The password of master-slave Redis which contains 6 to 36 characters, including letter，number，strikethrough and underline.
* `parameter_group_id` - (Optional) 主备版 Redis 配置组的 ID，如填写，则使用该配置组作为 Redis 配置文件。默认使用 Redis 默认配置文件。
* `parameter_group_id` - (Optional) The ID of param group of master-slave Redis, the specific param group applies if the ID is specified, otherwise the default param group applies.
* `tag` - (Optional) A tag to assign to the instance. The default value is "Default" (means no tag assigned).

~> **从库限制** Terraform 中，主备版 Redis 暂时不可以跨可用区。
~> **Note** The cross multiple zone of master-slave Redis is currently not supported by running Terrafrom provider.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ip_set` - ip_set is a nested type. ip_set documented below.
* `instance_type` - (Required) KV 存储实例的机型，表示实例引擎类型和容量的组合，可参考下方 [KV 存储实例机型表](/docs/providers/ucloud/appendix/kvstore_instance_type.html.markdown)。
* `instance_type` - (Required) The type of Key-Value storage instance which specifies the engine and storage size, please visit the instance type table for more details. (/docs/providers/ucloud/appendix/kvstore_instance_type.html.markdown).
* `create_time` - The time of creation for memory instance.
* `expire_time` - The expiration time for memory instance.
* `update_time` - The time whenever there is a change made to memory instance.
* `status` - KV 存储实例的状态
* `status` - The status of KV storage instance.

The attribute (`ip_set`) support the following:

* `ip` - KV 存储实例的虚拟 IP
* `ip` - The virtual ip of KV storage instance.
* `port` - KV 存储实例的端口
* `port` - The port on which Key-Value storage instance accepts connections, it is 6379 by default.
