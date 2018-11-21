---
layout: "ucloud"
page_title: "UCloud: ucloud_db_instance"
sidebar_current: "docs-ucloud-resource-db-param-group"
description: |-
  Provides a Database instance resource.
---

# ucloud_db_instance

Provides a Database instance resource.

~> **使用限制** 注意，当创建从库时主库id必传，且从库的参数与主库的参数一致，对于已创建的从库，如果此参数置空，则将当前从库提升为主库，与原主库分离。
## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where database instances are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `backup_zone` - (Optional) 跨可用区高可用备库所在可用区
* `master_id` - (Optional) 主库实例的id，当创建从库时必传；对于已创建的从库，如果此参数置空，则将当前从库提升为主库，与原主库分离
* `is_lock` - (Optional) 当创建从库时是否锁主库，默认为true，代表锁主库
* `is_force` - (Optional) 是否强制(如果从库落后可能会禁止提升)，如果落后情况下，强制提升丢失数, 默认false 
* `password` - (Optional) 管理员密码.
* `engine` - (Required) Database type, possible values are: "mysql", "percona", "postgresql".
* `engine_version` - (Required) The database engine version, possible values are: "5.1", "5.5", "5.6", "5.7", "9.4", "9.6", "10.4".
* `name` - (Optional)  实例名称
* `instance_storage` - (Optional) 磁盘空间(GB), 暂时支持20G - 3000G
* `param_group_id` - (Optional) DB实例使用的配置参数组id.
* `instance_class` - (Required) 数据库机型.基本格式为"engine-type-memory",其中 engine 可以为"mysql","percona","postgresql"；type可以为"basic","ha",分别代表普通版和高可用版，高可用版实例采用双主热备架构，可以彻底解决因宕机或硬件故障造成的数据库不可用，mysql与percona只支持高可用版，postgresql现只支持普通版；memory可以为1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64,单位GB
* `port` - (Optional) 端口号，mysql与percona默认3306，postgresql默认5432
* `instance_charge_type` - (Optional) The charge type of database instance, possible values are: "Year", "Month" and "Dynamic" as pay by hour (specific permission required). the dafault is "Month".
* `instance_duration` - (Optional) The duration that you will buy the resource, the default value is "1". It is not required when "Dynamic" (pay by hour), the value is "0" when pay by month and the instance will be vaild till the last day of that month.
* `vpc_id` - (Optional) The ID of VPC linked to the database instances.
* `subnet_id` - (Optional) The ID of subnet.
* `backup_count` - (Optional) 备份策略，每周备份数量，默认7次
* `backup_duration` - (Optional) 备份策略，备份开始时间，单位小时计，默认1点
* `backup_begin_time` - (Optional) 备份策略，备份开始时间，单位小时计，默认1点
* `backup_date` - (Optional) 备份时期标记位。共7位，每一位为一周中一天的备份情况，0表示关闭当天备份，1表示打开当天备份。最右边的一位为星期天的备份开关，其余从右到左依次为星期一到星期六的备份配置开关，每周必须至少设置两天备份。例如：1100000表示打开星期六和星期五的备份功能.
* `backup_id` - (Optional) 备份id，如果指定，则表明从备份恢复实例
* `backup_black_list` - (Optional) 黑名单，规范示例,指定库mysql.%;test.%; 指定表city.address;

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - DB状态标记 Init：初始化中，Fail：安装失败，Starting：启动中，Running：运行，Shutdown：关闭中，Shutoff：已关闭，Delete：已删除，Upgrading：升级中，Promoting：提升为独库进行中，Recovering：恢复中，Recover fail：恢复失败
* `create_time` - DB实例创建时间，采用unix计时时间戳
* `expire_time` - DB实例修改时间，采用unix计时时间戳
* `modify_time` - DB实例过期时间，采用unix计时时间戳
* `role` - DB实例角色，区分master/slave