# Specify the provider and access details
provider "ucloud" {
  region = "${var.primary_region}"
}

data "ucloud_zones" "default" {}

# Create the primary cluster nodes in different availability zones.
# It is data center level standby to ensure high availability of application.
module "app-primary" {
  source = "./web"
  region = "${var.primary_region}"

  instance_password = "${var.password}"
  instance_count    = "${var.primary_node_count}"

  vpc_network    = ["192.168.0.0/16"]
  subnet_network = "192.168.1.0/24"

  cross_zone = true
}

# Create the replication cluster nodes in different regions.
# It is city level standby to ensure high availability of application.
module "app-replication-cross-region" {
  source = "./web"
  region = "${var.replica_region}"

  instance_password = "${var.password}"
  instance_count    = "${var.replica_node_count}"

  vpc_network    = ["192.168.0.0/16"]
  subnet_network = "192.168.2.0/24"

  cross_zone = false
}

# Create UCloud Dedicated Private Network.
# You can use it to achieve high-speed, stable, secure, and dedicated communications between different data centers.
# The most frequent scenario is to create network connection of clusters across regions.
resource "ucloud_udpn_connection" "default" {
  bandwidth   = 2
  peer_region = "${var.replica_region}"
}

# Create a VPC Peering Connection for establish a connection to put multi vpc network into same network plane.
resource "ucloud_vpc_peering_connection" "default" {
  vpc_id      = "${module.app-primary.vpc_id}"
  peer_region = "${var.replica_region}"
  peer_vpc_id = "${module.app-replication-cross-region.vpc_id}"

  # UDPN connection is required by vpc peering connection across multi-region.
  depends_on = ["ucloud_udpn_connection.default"]
}
