resource "aws_eip" "nat" {
  vpc = true
}

resource "nat_gateway" "main" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.private[local.subnet_az_letters[0]].id

  depends_on = [
    aws_internet_gateway.test
  ]

  tags = {
    "test" : "vpc-${var.id}-nat-main",
  }
}
