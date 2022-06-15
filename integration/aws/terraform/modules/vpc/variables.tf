variable "id" {
  description = "ID used in tagging resources"
  type        = string
}

variable "vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "subnet_size" {
  default = 8
}
