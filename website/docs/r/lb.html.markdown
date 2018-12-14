---
layout: "ucloud"
page_title: "UCloud: ucloud_lb"
sidebar_current: "docs-ucloud-resource-lb"
description: |-
  Provides a Load Balancer resource.
---

# ucloud_lb

Provides a Load Balancer resource.

## Example Usage

```hcl
resource "ucloud_lb" "web" {
    name = "tf-example-lb"
    tag  = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `internal` - (Optional) Indicate whether the load balancer is intranet.
* `internet_charge_type` - (Optional) Charge type of load balancer. Possible values are: `"Year"` as pay by year, `"Month"` as pay by month, `"Dynamic"` as pay by hour (specific permission required). (Default: `"Month"`).
* `name` - (Optional) The name of the load balancer. (Default: `"LB"`).
* `remark` - (Optional) The remarks of the load balancer. (Default: is `""`).
* `subnet_id` - (Optional) The ID of subnet that intrant load balancer belongs to. This argumnet is not required if default subnet.
* `tag` - (Optional) A mapping of tags to assign to the load balancer, which contains 1-63 characters and only support Chinese, English, numbers, '-', '_', '.'. (Default: `"Default"`, means no tag assigned). 
* `vpc_id` - (Optional) The ID of the VPC linked to the Load Balancers, This argumnet is not required if default VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The time of creation for load balancer, formatted in RFC3339 time string.
* `expire_time` - The expiration time for load balancer, formatted in RFC3339 time string.
* `ip_set` - It is a nested type which documented below.
* `private_ip` - The IP address of intranet IP. It is `""` if `internal` is `"false"`.

The attribute (`ip_set`) support the following:

* `eip_id` - The ID of EIP.
* `internet_type` - Elastic IP routes. Possible values are: `"International"` as internaltional IP and `"Bgp"` as BGP IP.
* `ip` - Elastic IP address.
