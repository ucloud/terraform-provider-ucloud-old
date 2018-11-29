variable "region" {
  description = "The region that will create resources in"
  default     = "cn-sh2"
}

variable "password" {
  default = "wA123456"
}

variable "slave_count" {
  default = 2
}

variable "instance_type" {
  default = "redis-master-1"
}
