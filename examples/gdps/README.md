# GDPS Example

The GDPS(Geographically Dispersed Parallel Sysplex) example will create instances in 2 region and 3 data center to provide high availability.

To run, configure your UCloud provider as described in https://www.terraform.io/docs/providers/ucloud/index.html

## Features

* create a module for reusing, include vpc, security group and instance
* instance for web application
    * 1 primary
    * 1 cross availability zone replication
    * 1 cross region replication
* create cross region dedicated private network
* put them into same VPC

## Setup Environment

```sh
export UCLOUD_PUBLIC_KEY="your public key"
export UCLOUD_PRIVATE_KEY="your private key"
export UCLOUD_PROJECT_ID="your project id"
```

## Running the example

run `terraform apply`
