resource "aws_default_security_group" "default" {
  vpc_id = aws_vpc.vpc.id

  ingress = []
  egress  = []
}
