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
~> **Note** Please try to keep the same settings for both master and slave DBs, otherwise it's likely to have an issue of data loss when making synchronization. The slave DB creation is not supported for mysql and percona in 5.5 version.
## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Availability zone where database slaves are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `master_id` - (Required) 主库实例的id，当创建从库时必传；对于已创建的从库，如果此参数置空，则将当前从库提升为主库，与原主库分离
* `master_id` - (Required) The ID of master DB instance, it is mandatory required to request when creating slave DB; the current slave DB will be promoted to new master one and separate from the origin master DB if it is set to "null" for the existing slave DB.
* `is_lock` - (Optional) 当创建从库时是否锁主库，默认为true，代表锁主库
* `is_lock` - (Optional) Specifies whether need to set master DB to read only when creating slave DB, possible values are "true" and "False", it is "true" by default.
* `password` - (Optional) 管理员密码.
* `password` - (Optional) The password for the database instance which should have 8-30 characters. It must contain at least 3 items of Capital letters, small letter, numbers and special characters. The special characters include <code>`()~!@#$%^&*-+=_|{}\[]:;'<>,.?/</code> When it is changed, the instance will reboot to make the change take effect.
* `name` - (Optional)  实例名称,The name of the DB instance, should have 1 - 63 characters and only support chinese, english, numbers, '-', '_', '.'.
* `instance_storage` - (Optional) 磁盘空间(GB), 暂时支持20G - 3000G
* `instance_storage` - (Optional) Specifies the allocated storage size in gigabytes (GB), range from 20 to 3000GB. The volume adjustment must be a multiple of 10 GB. When it is changed, the instance will reboot to make the change take effect.
* `parameter_group_id` - (Optional) DB实例使用的配置参数组id.
* `parameter_group_id` - (Optional) The ID of DB param group. Note: the "parameter_group_id" should be included in the request for the multiple zone DB instance. 
* `instance_type` - (Required) 数据库机型.基本格式为"engine-type-memory",其中 engine 可以为"mysql","percona","postgresql"；type可以为"basic","ha",分别代表普通版和高可用版，高可用版实例采用双主热备架构，可以彻底解决因宕机或硬件故障造成的数据库不可用，mysql与percona只支持高可用版，postgresql现只支持普通版；memory可以为1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64,单位GB
* `instance_type` - (Required) Specifies the type of DB instance with format "engine-type-memory", Possible values are:
  "mysql","percona" and "postgresql" for engine;
  "basic" as normal verison and  "ha" as high availability version for type, the dual mian hot standby structure which can thoroughly solved the issue of unsysnchronized DB caused by the system downtime or DB unavailable, the "ha" version only supports "mysql" and "percona" engine, the standard version only supports the "postgrsql" engine.
Possible values for memory are: 1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64GB.  
* `port` - (Optional) 端口号，mysql与percona默认3306，postgresql默认5432
* `port` - (Optional) The port on which the DB accepts connections, the default port is 3306 for mysql and percona and 5432 for postgresql.
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - DB状态标记 Init：初始化中，Fail：安装失败，Starting：启动中，Running：运行，Shutdown：关闭中，Shutoff：已关闭，Delete：已删除，Upgrading：升级中，Promoting：提升为独库进行中，Recovering：恢复中，Recover fail：恢复失败
* `status` - Specifies the status of DB, possible values are: "Init","Fail", "Starting", "Running", "Shutdown", "shutoff", "Delete", "Upgrading", "Promoting", "Recovering" and "Recover fail".
* `create_time` - DB实例创建时间，采用unix计时时间戳
* `create_time` - The creation time of DB, format in Unix timestamp.
* `expire_time` - DB实例修改时间，采用unix计时时间戳
* `expire_time` - The expiration time of DB, format in Unix timestamp.
* `modify_time` - DB实例过期时间，采用unix计时时间戳
* `modify_time` - The modification time of DB, format in Unix timestamp.
