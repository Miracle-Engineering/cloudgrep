locals {
  alb_count = 2
}

resource "aws_lb" "alb" {
  count = local.alb_count

  name_prefix        = "test2-"
  load_balancer_type = "application"
  internal           = true
  security_groups    = [module.vpc.default_sg_id]
  subnets            = module.vpc.private_subnet_ids

  tags = {
    test = "elb-alb-${count.index}"
  }
}
