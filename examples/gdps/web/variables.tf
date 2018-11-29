variable "region" {
  description = "The region to create resources in"
  default     = "cn-sh2"
}

variable "instance_password" {
  default = "wA123456"
}

variable "instance_count" {
  default = "1"
}

variable "count_format" {
  default = "%02d"
}

variable "vpc_network" {
  default = ["192.168.0.0/16"]
}

variable "subnet_network" {
  default = "192.168.1.0/24"
}

variable "cross_zone" {
  default = true
}
