---
layout: "ucloud"
page_title: "UCloud: ucloud_eip"
sidebar_current: "docs-ucloud-resource-eip"
description: |-
  Provides an Elastic IP resource.
---

# ucloud_eip

Provides an Elastic IP resource.

## Example Usage

```hcl
resource "ucloud_eip" "example" {
    bandwidth            = 2
    internet_charge_mode = "Bandwidth"
    name                 = "tf-example-eip"
    tag                  = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). the ranges for bandwidth are: 1-200 for pay by traffic, 1-800 for pay by bandwith. (Default: `"1"`).
* `eip_duration` - (Optional) The duration that you will buy the resource. (Default: `"1"`). It is not required when `"Dynamic"` (pay by hour), the value is `"0"` when `"Month"`(pay by month) and the instance will be vaild till the last day of that month.
* `internet_charge_mode` -(Optional) Elastic IP charge mode. Possible values are: `"Traffic"` as pay by traffic, `"Bandwidth"` as pay by bandwidth, `"ShareBandwidth"` as pay by shared bandwidth. (Default: `"Bandwidth"`).
* `internet_charge_type` - (Optional) Elastic IP charge type. Possible values are: `"Year"` as pay by year, `"Month"` as pay by month, `"Dynamic"` as pay by hour (specific permission required). (Default: `"Month"`).
* `internet_type` - (Optional) Elastic IP routes. Possible values are: `"International"` as internaltional IP and `"Bgp"` as BGP IP. (Default: `"Bgp"`).
* `name` - (Optional) The name of the EIP, which contains 1-63 characters and only support Chinese, English, numbers, '-', '_', '.'. (Default: `"EIP"`).
* `remark` - (Optional) The remarks of the EIP. (Default: `""`).
* `tag` - (Optional) A mapping of tags to assign to the EIP, which contains 1-63 characters and only support Chinese, English, numbers, '-', '_', '.'. (Default: `"Default"`, means no tag assigned). 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for EIP, formatted in RFC3339 time string.
* `expire_time` - The expiration time for EIP, formatted in RFC3339 time string.
* `ip_set` - It is a nested type which documented below.
* `resource` - It is a nested type which documented below.
* `status` - EIP status. Possible values are: `"used"` as in use, `"free"` as available and `"freeze"` as associating.

The attribute (`ip_set`) support the following:

* `internet_type` - Elastic IP routes. Possible values are: `"International"` as internaltional IP and `"Bgp"` as BGP IP.
* `ip` - Elastic IP address.

The attribute (`resource`) support the following:

* `eip_id` - The ID of EIP.
* `resource_id` - The ID of the resource with EIP attached.
* `resource_type` - The type of resource with EIP attached. Possible values are `"instance"` as instance, `"vrouter"` as visual router, `"lb"` as load balancer.
