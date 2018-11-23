---
layout: "ucloud"
page_title: "UCloud: ucloud_kvstore_instance"
sidebar_current: "docs-ucloud-resource-kvstore-instance"
description: |-
  Provides an KVStore Instance resource.
---

# ucloud_kvstore_instance

UCloud KV 存储实例是兼容 Redis 和 Memcache 协议的 Key-Value 在线存储服务。

~> **注意** Memcache 应仅用于缓存，不可用于数据持久化存储。当容灾切换，或扩容操作后，数据将被清空。

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where instance is located. such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `name` - (Required) KV 存储实例名称，6-63 个字符，支持字母、数字、中划线或下划线。
* `engine` - (Required) KV 存储实例的引擎名称，支持 memcache 或 redis。
* `instance_type` - (Required) KV 存储实例的机型，表示实例引擎类型和容量的组合，可参考下方 [KV 存储实例机型表](/docs/providers/ucloud/appendix/kvstore_instance_type.html.markdown)。
* `instance_charge_type` - (Optional) Charge type. Possible values are: "Year" as pay by year, "Month" as pay by month, "Dynamic" as pay by hour (specific permission required). The default value is "Month".
* `instance_duration` - (Optional) The duration that you will buy the resource, the default value is "1". It is not required when "Dynamic" (pay by hour), the value is "0" when pay by month and the instance will be vaild till the last day of that month.
* `tag` - (Optional) A tag to assign to the instance. The default value is "Default" (means no tag assigned).
* `port` - KV 存储实例的端口，默认 6379

主备版 Redis 特有参数（Required 表示对于主备版必填，而对其它机型选填）：

主备版特有的特性有：密码，备份，多 DB，指定版本

* `engine_version` - (Required) 主备版 Redis 实例的版本名称，支持 3.0, 3.2, 4.0
* `password` - (Optional) 主备版 Redis 实例的密码，6-36 个字符，支持字母、数字、中划线或下划线。
* `parameter_group_id` - (Optional) 主备版 Redis 配置组的 ID，如填写，则使用该配置组作为 Redis 配置文件。默认使用 Redis 默认配置文件。
* `backup_policy` - (Optional) 主备版 Redis 的备份策略。默认不启用备份策略。backup_policy is a nested type. backup_policy documented below.
* `backup_id` - (Optional) 主备版 Redis 备份的 ID，如填写，则表示当前实例是基于该备份恢复并创建的。默认不填写。

The argument (`backup_policy`) support the following:

* `backup_time` - 自动备份开始时间，范围 [0-23]，表示每天在几点钟触发自动备份，采用 UTC+8 时间。

~> **从库限制** Terraform 中，主备版 Redis 暂时不可以跨可用区。

~> **备份限制** 如果想使用备份策略/备份恢复特性，必须先在对应地域开通 UFile 产品。

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ip_set` - ip_set is a nested type. ip_set documented below.
* `create_time` - The time of creation for memory instance.
* `expire_time` - The expiration time for memory instance.
* `update_time` - The time whenever there is a change made to memory instance.
* `status` - KV 存储实例的状态

The attribute (`ip_set`) support the following:

* `ip` - KV 存储实例的虚拟 IP
* `port` - KV 存储实例的端口
