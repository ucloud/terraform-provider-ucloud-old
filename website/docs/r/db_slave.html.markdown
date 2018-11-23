---
layout: "ucloud"
page_title: "UCloud: ucloud_db_slave"
sidebar_current: "docs-ucloud-resource-db-param-group"
description: |-
  Provides a Database slave resource.
---

# ucloud_db_slave

Provides a Database slave resource.

~> **使用限制** 注意，请尽量保持从库与主库配置的一致性，否则在数据同步时可能会出现数据丢失。"5.5"版的"mysql","percona"不支持创建从库
## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where database slaves are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `master_id` - (Required) 主库实例的id，当创建从库时必传；对于已创建的从库，如果此参数置空，则将当前从库提升为主库，与原主库分离
* `is_lock` - (Optional) 当创建从库时是否锁主库，默认为true，代表锁主库
* `password` - (Optional) 管理员密码.
* `name` - (Optional)  实例名称
* `instance_storage` - (Optional) 磁盘空间(GB), 暂时支持20G - 3000G
* `parameter_group_id` - (Optional) DB实例使用的配置参数组id.
* `instance_type` - (Required) 数据库机型.基本格式为"engine-type-memory",其中 engine 可以为"mysql","percona","postgresql"；type可以为"basic","ha",分别代表普通版和高可用版，高可用版实例采用双主热备架构，可以彻底解决因宕机或硬件故障造成的数据库不可用，mysql与percona只支持高可用版，postgresql现只支持普通版；memory可以为1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64,单位GB
* `port` - (Optional) 端口号，mysql与percona默认3306，postgresql默认5432

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - DB状态标记 Init：初始化中，Fail：安装失败，Starting：启动中，Running：运行，Shutdown：关闭中，Shutoff：已关闭，Delete：已删除，Upgrading：升级中，Promoting：提升为独库进行中，Recovering：恢复中，Recover fail：恢复失败
* `create_time` - DB实例创建时间，采用unix计时时间戳
* `expire_time` - DB实例修改时间，采用unix计时时间戳
* `modify_time` - DB实例过期时间，采用unix计时时间戳