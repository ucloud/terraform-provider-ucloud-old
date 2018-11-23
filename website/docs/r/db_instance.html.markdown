---
layout: "ucloud"
page_title: "UCloud: ucloud_db_instance"
sidebar_current: "docs-ucloud-resource-db-param-group"
description: |-
  Provides a Database instance resource.
---

# ucloud_db_instance

Provides a Database instance resource.

~> **使用限制** 注意，普通版数据库配置升降级时，大约需要关闭实例5分钟，请提前安排好您的业务！为避免数据丢失。重置密码时，请先确认是否有待提交事务，重置密码会立即生效，请谨慎操作。
## Example Usage

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) Availability zone where database instances are located. Such as: "cn-bj-01". You may refer to [list of availability zone](https://docs.ucloud.cn/api/summary/regionlist)
* `backup_zone` - (Optional) 跨可用区高可用备库所在可用区；单可用区高可用实例可以承受服务器和机架级别的故障；跨可用区高可用实例可以承受机房级别的故障；注意：因为多可用区之间存在一定的网络延迟，对于单个更新的响应时间会比单可用区高可用实例长
* `password` - (Optional) 管理员密码.The password for the database instance, should have between 8-30 characters.It must contain least 3 items of Capital letters, small letter, numbers and special characters. The special characters incloud <code>`()~!@#$%^&*-+=_|{}\[]:;'<>,.?/</code> When it is changed, the instance will reboot to make the change take effect.
* `engine` - (Required) Database type, possible values are: "mysql", "percona", "postgresql".
* `engine_version` - (Required) The database engine version, possible values are: "5.5", "5.6", "5.7", "9.4", "9.6".其中"mysql"和"percona"只支持 "5.5", "5.6", "5.7"，且"5.5"版本不支持创建从库，postgresql只支持"9.4", "9.6"版本
* `name` - (Optional)  实例名称，should have 1 - 63 characters and only support chinese, english, numbers, '-', '_', '.'.
* `instance_storage` - (Optional) 磁盘空间(GB), 暂时支持20G - 3000G；硬盘步长10G。SSD机型：内存8G及以下时硬盘容量上限为500G，内存12~24G时硬盘容量上限为1000G，内存32G时硬盘容量上限为2000G，内存48G及以上时硬盘容量上限为3000G。
* `parameter_group_id` - (Optional) DB实例使用的配置参数组id.注意：对于跨可用区高可用实例，需要传入跨可用区配置参数组id
* `instance_type` - (Required) 数据库机型.基本格式为"engine-type-memory",其中 engine 可以为"mysql","percona","postgresql"；type可以为"basic","ha",分别代表普通版和高可用版，高可用版实例采用双主热备架构，可以彻底解决因宕机或硬件故障造成的数据库不可用，mysql与percona只支持高可用版，postgresql现只支持普通版；memory可以为1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64,单位GB
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
