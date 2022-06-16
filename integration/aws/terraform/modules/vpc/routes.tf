resource "aws_route_table" "private" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    "test" : "vpc-${var.id}-route-table-private"
  }

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main.id
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    "test" : "vpc-${var.id}-route-table-public"
  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.test.id
  }
}
