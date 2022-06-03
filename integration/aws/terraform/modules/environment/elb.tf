locals {
  alb_count = 2
}

resource "aws_lb" "alb" {
  count = local.alb_count

  name_prefix        = "test-"
  load_balancer_type = "application"
  internal = true
  security_groups    = [aws_default_security_group.default.id]
  subnets            = data.aws_subnets.default.ids

  tags = {
    test = "elb-alb-${count.index}"
  }
}
