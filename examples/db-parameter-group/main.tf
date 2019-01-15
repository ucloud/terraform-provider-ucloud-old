# Specify the provider and access details
provider "ucloud" {
  region = "${var.region}"
}

# Query availability zone
data "ucloud_zones" "default" {}

# Query parameter group
data "ucloud_db_parameter_groups" "default" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  multi_az       = "false"
  engine            = "mysql"
  engine_version    = "5.7"
}

# Create parameter group
resource "ucloud_db_parameter_group" "example" {
  availability_zone = "${data.ucloud_zones.default.zones.0.id}"
  name              = "tf-example-parameter-group"
  description       = "this is a test"
  src_group_id      = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
  engine            = "mysql"
  engine_version    = "5.7"

  parameter_input {
    key   = "max_connections"
    value = "3000"
  }

  parameter_input {
    key   = "slow_query_log"
    value = "1"
  }
}
