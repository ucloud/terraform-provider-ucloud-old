---
layout: "ucloud"
page_title: "UCloud: ucloud_udpn_connection"
sidebar_current: "docs-ucloud-resource-udpn-connection"
description: |-
  Provides an UDPN Connection resource.
---

# ucloud_udpn_connection

UDPN (UCloud Dedicated Private Network)，提供各个数据中心之间的，低延迟、高质量的内网数据传输服务。常用于跨地域的集群内网互通，
UDPN (UCloud Dedicated Private Network)，you can use Dedicated Private Network to achieve high-speed, stable, secure, and dedicated communications between different data centers. The most frequent scenario is to create network connection of clusters across regions.
~> **与 VPC 产品的关联** 跨地域的 VPC 对等连接，必须先在两个地域之间创建 DPN，以保证专有网络之间的内网联通性。
~> **VPC interconnection** The cross-region Dedicated Private Network must be established if the two VPCs of different regions are expected to be connected.
~> **注意事项** UDPN高速通道使用隧道封装技术，会产生数据包头的额外开销，且计算进数据包的总长度中。由于隧道的包头字节数固定， 数据报文越大，隧道所占的开销占比越小。
~> **Note** The addtional packet head will be added and included in the overall length of packet due to the tunneling UDPN adopted. Since the number of the bytes of packet head is fixed, the bigger data packet is, the less usage will be taken for the packet head.
## Example Usage

```hcl
resource "ucloud_udpn_connection" "example" {
    bandwidth = 2
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). 取值范围 0 - 1000M. The default value is "1".
* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). range from 0 - 1000M. The default value is "1".
* `duration` - (Optional) The duration that you will buy the resource, the default value is "1". It is not required when "Dynamic" (pay by hour), the value is "0" when pay by month and the instance will be vaild till the last day of that month.
* `charge_type` - (Optional) Charge type. Possible values are: "Year" as pay by year, "Month" as pay by month, "Dynamic" as pay by hour. The default value is "Month".
* `peer_region` - (Optional) 专线对端地域，请 [参考地域与可用区列表](https://docs.ucloud.cn/api/summary/regionlist) 和 [UDPN 线路价格列表](https://docs.ucloud.cn/network/udpn/udpn_price)
* `peer_region` - (Optional) The correspondent region of dedicated connection, please refer to the region and availability zone list (https://docs.ucloud.cn/api/summary/regionlist) and UDPN price list (https://docs.ucloud.cn/network/udpn/udpn_price)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for DPN.
* `expire_time` - The expiration time for DPN.
