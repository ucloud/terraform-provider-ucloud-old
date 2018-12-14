---
layout: "ucloud"
page_title: "UCloud: ucloud_lb_listener"
sidebar_current: "docs-ucloud-resource-lb-listener"
description: |-
  Provides a Load Balancer Listener resource.
---

# ucloud_lb_listener

Provides a Load Balancer Listener resource.

## Example Usage

```hcl
resource "ucloud_lb" "web" {
    name = "tf-example-lb"
    tag  = "tf-example"
}

resource "ucloud_lb_listener" "example" {
    load_balancer_id = "${ucloud_lb.web.id}"
    protocol         = "HTTPS"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of load balancer instance.
* `protocol` - (Required) Listener protocol. Possible values: `"HTTP"`, `"HTTPS"` if `"ListenType"` is `"RequestProxy"`, `"TCP"` and `"UDP"` if `"ListenType"` is `"PacketsTransmit"`.
* `name` - (Optional) The name of the listener. (Default: `"Listener"`).
* `listen_type` - (Optional) The type of listener. Possible values are `"RequestProxy"` and `"PacketsTransmit"`. (Default: `"PacketsTransmit"`).
* `port` - (Optional) Port opened on the listeners to receive requests, range: 1-65535. (Default: `"80"`).
* `idle_timeout` - (Optional) Amount of time in seconds to wait for the response for in between two sessions if `"ListenType"` is `"RequestProxy"`, range: 0-86400. (Default: `"60"`). Amount of time in seconds to wait for one session if `"ListenType"` is `"PacketsTransmit"`, range: 60-900. The session will be closed as soon as no response if it is `"0"`.
* `method` - (Optional) The load balance method in which the listener is. Possible values are: `"Roundrobin"`, `"Source"`, `"ConsistentHash"`, `"SourcePort"` , `"ConsistentHashPort"`, `"WeightRoundrobin"` and `"Leastconn"`. The `"ConsistentHash"`, `"SourcePort"` , `"ConsistentHashPort"`, `"Roundrobin"`, `"Source"` and `"WeightRoundrobin"` are valid if `"listen_type"` is `"PacketsTransmit"`. The `"Roundrobin"`, `"Source"` and `"WeightRoundrobin"` and `"Leastconn"` are vaild if `"listen_type"` is `"RequestProxy"`. (Default: `"Roundrobin"`).
* `persistence` - (Optional) Indicate whether the persistence session is enabled, it is invaild if `"PersistenceType"` is `"None"`, an auto-generated string will be exported if `"PersistenceType"` is `"ServerInsert"`, a custom string will be exported if `"PersistenceType"` is `"UserDefined"`.
* `persistence_type` - (Optional) The type of session persistence of listener. Possible values are: `"None"` as disabled, `"ServerInsert"` as auto-generated string and `"UserDefined"` as cutom string. (Default: `"None"`).
* `health_check_type` - (Optional) Health check method. Possible values are `"Port"` as port checking and `"Path"` as http checking.
* `path` - (Optional) Health check path checking.
* `domain` - (Optional) Health check domain checking.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Listener status. Possible values are: `"allNormal"` for all resource functioning well, `"partNormal"` for partial resource functioning well and `"allException"` for all resource functioning exceptional.
`