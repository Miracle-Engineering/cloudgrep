data "aws_region" "current" {}

variable "id" {
  description = "ID used in naming and tagging resources"
  type        = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "vpc_id" {
  type = string
}

variable "max_nodes" {
  type    = number
  default = 3
}

variable "min_nodes" {
  type    = number
  default = 1
}

variable "node_disk_size" {
  type    = number
  default = 20
}

variable "node_instance_type" {
  type    = string
  default = "t3.medium"
}

variable "k8s_version" {
  type    = string
  default = "1.21"
}

variable "spot_instances" {
  type    = bool
  default = true
}

variable "eks_log_retention" {
  type    = number
  default = 7
}

variable "ami_type" {
  type    = string
  default = "AL2_x86_64"
}
