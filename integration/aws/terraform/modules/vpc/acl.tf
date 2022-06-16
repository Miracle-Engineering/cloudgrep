locals {
  acl_count = 1
}

resource "aws_network_acl" "test" {
  count = local.acl_count

  vpc_id = aws_vpc.vpc.id

  egress {
    action    = "deny"
    rule_no   = 100
    protocol  = "-1"
    from_port = 0
    to_port   = 0
  }

  ingress {
    action    = "deny"
    rule_no   = 100
    protocol  = "-1"
    from_port = 0
    to_port   = 0
  }

  tags = {
    "test" : "vpc-${var.id}-acl-${count.index}"
  }
}
