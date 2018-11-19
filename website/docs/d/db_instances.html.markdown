---
layout: "ucloud"
page_title: "UCloud: ucloud_db_instances"
sidebar_current: "docs-ucloud-datasource-db-instances"
description: |-
  Provides a list of DB instance resources in the current region.
---

# ucloud_db_instances

This data source providers a list of DB instance resources according to their availability zone, DB instance ID, class type.

~> **使用限制** 注意，如果传入id则查询当前数据库id对应的资源信息，id可以是主库id或从库id；不传id时必须要传class_type，对于主库，则可以返回当前region下所有主库及其从库的信息，如果还传入availability_zone则只返回当前当前可用区下所有主库和从库的信息

## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where db instances are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `ids` - (Optional) The group of IDs of db instances that require to be retrieved, 当id为空时，将根据class type和Availability zone来检索，其中class type为必填.
* `class_type` - (Optional) DB种类，分为sql和postgresql,其中，sql代表mysql和percona,如果没有指定ids则此参数必填
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `tags` - (Optional) A mapping of tags to assign to instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_instances` - db instances is a nested type. instances documented below.
* `total_count` - Total number of instance that satisfy the condition.

The attribute (`db_instances`) support the following:

* `id` - The ID of db instance.
* `zone` - Availability zone where db instances are located.
* `backup_zone` - 跨可用区高可用备库所在可用区
* `name` - 实例名称
* `engine` - Database type.
* `engine_version` - The database engine version.
* `param_group_id` - DB实例使用的配置参数组id
* `src_db_id` - 对slave而言是master的DBId
* `instance_charge_type` - The charge type of db instance, possible values are: "Year", "Month" and "Dynamic" as pay by hour.
* `instance_storage` - 磁盘空间(GB)
* `memory` - 内存限制(GB)
* `role` - DB实例角色，区分master/slave.
* `disk_used_size` - DB实例磁盘已使用空间，单位GB
* `data_file_size` - DB实例数据文件大小，单位GB
* `system_file_size` - DB实例系统文件大小，单位GB
* `log_file_size` - DB实例日志文件大小，单位GB
* `backup_count` - 备份策略，不可修改，备份文件保留的数量，默认7次
* `backup_begin_time` - 备份策略，不可修改，开始时间，单位小时计，默认3点
* `backup_duration` - 备份策略，一天内备份时间间隔，单位小时，默认24小时
* `backup_blacklist` - 备份策略，备份黑名单，mongodb则不适用
* `backup_date` - 备份日期标记位。共7位,每一位为一周中一天的备份情况 0表示关闭当天备份,1表示打开当天备份。最右边的一位 为星期天的备份开关，其余从右到左依次为星期一到星期 六的备份配置开关，每周必须至少设置两天备份。 例如：1100000 表示打开星期六和星期五的自动备份功能
* `instance_mode` - UDB实例模式类型, 可选值如下: “Normal”： 普通版UDB实例 “HA”: 高可用版UDB实例
* `status` - DB状态标记 Init：初始化中，Fail：安装失败，Starting：启动中，Running：运行，Shutdown：关闭中，Shutoff：已关闭，Delete：已删除，Upgrading：升级中，Promoting：提升为独库进行中，Recovering：恢复中，Recover fail：恢复失败
* `create_time` - DB实例创建时间，采用UTC计时时间戳
* `expire_time` - DB实例修改时间，采用UTC计时时间戳
* `modify_time` - DB实例过期时间，采用UTC计时时间戳
* `port` - 端口号.
* `ip_set` - ip_set is a nested type. ip_set documented below.
* `slave_instances` - db instances is a nested type. instances documented below.

The attribute (`slave_instances`) support the following:

* `id` - The ID of db instance.
* `zone` - Availability zone where db instances are located.
* `name` - 实例名称
* `engine` - Database type.
* `engine_version` - The database engine version.
* `param_group_id` - DB实例使用的配置参数组id
* `src_db_id` - 对slave而言是master的DBId
* `instance_charge_type` - The charge type of db instance, possible values are: "Year", "Month" and "Dynamic" as pay by hour.
* `instance_storage` - 磁盘空间(GB)
* `memory` - 内存限制(GB)
* `role` - DB实例角色，区分master/slave.
* `disk_used_size` - DB实例磁盘已使用空间，单位GB
* `data_file_size` - DB实例数据文件大小，单位GB
* `system_file_size` - DB实例系统文件大小，单位GB
* `log_file_size` - DB实例日志文件大小，单位GB
* `backup_count` - 备份策略，不可修改，备份文件保留的数量，默认7次
* `backup_begin_time` - 备份策略，不可修改，开始时间，单位小时计，默认3点
* `backup_duration` - 备份策略，一天内备份时间间隔，单位小时，默认24小时
* `backup_blacklist` - 备份策略，备份黑名单，mongodb则不适用
* `backup_date` - 备份日期标记位。共7位,每一位为一周中一天的备份情况 0表示关闭当天备份,1表示打开当天备份。最右边的一位 为星期天的备份开关，其余从右到左依次为星期一到星期 六的备份配置开关，每周必须至少设置两天备份。 例如：1100000 表示打开星期六和星期五的自动备份功能
* `instance_mode` - UDB实例模式类型, 可选值如下: “Normal”： 普通版UDB实例 “HA”: 高可用版UDB实例
* `status` - DB状态标记 Init：初始化中，Fail：安装失败，Starting：启动中，Running：运行，Shutdown：关闭中，Shutoff：已关闭，Delete：已删除，Upgrading：升级中，Promoting：提升为独库进行中，Recovering：恢复中，Recover fail：恢复失败
* `create_time` - DB实例创建时间，采用UTC计时时间戳
* `expire_time` - DB实例修改时间，采用UTC计时时间戳
* `modify_time` - DB实例过期时间，采用UTC计时时间戳
* `port` - 端口号.
* `ip_set` - ip_set is a nested type. ip_set documented below.

The attribute (`ip_set`) support the following:

* `ip` - IP address.