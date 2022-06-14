data "aws_region" "current" {}
resource "aws_default_vpc" "default" {
  tags = {
    "test": "vpc-default"
  }
}

resource "aws_default_security_group" "default" {
  vpc_id = aws_default_vpc.default.id
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [aws_default_vpc.default.id]
  }
}

locals {
  vpc_azs = [
    "${data.aws_region.current.name}a",
    "${data.aws_region.current.name}b",
    "${data.aws_region.current.name}c",
  ]
}
