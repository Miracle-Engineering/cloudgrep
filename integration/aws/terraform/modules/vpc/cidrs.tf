locals {
  vpc_size = tonumber(split("/", var.vpc_cidr)[1])
}

module "subnet_addrs" {
  source  = "hashicorp/subnets/cidr"
  version = "1.0.0"

  base_cidr_block = var.vpc_cidr

  networks = [for s in local.subnet_az_letters : {
    name     = s
    new_bits = (32 - var.subnet_size - local.vpc_size)
  }]
}
