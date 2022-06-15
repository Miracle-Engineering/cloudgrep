resource "aws_route_table" "private" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    "test" : "vpc-${var.id}-route-table-private"
  }

  route = []
}
