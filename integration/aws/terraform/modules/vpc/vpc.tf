resource "aws_vpc" "vpc" {
  cidr_block = var.vpc_cidr
  tags = {
    "test" : "vpc-${var.id}"
    "Name" : var.id
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    "test" : "vpc-${var.id}-igw"
  }
}
