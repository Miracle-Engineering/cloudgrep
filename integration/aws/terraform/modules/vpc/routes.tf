resource "aws_route_table" "private" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    "test" : "vpc-${var.id}-route-table-private"
  }

  // We don't actually want to route anything, to avoid EC2 instances from connecting to the internet
  route = []
}
