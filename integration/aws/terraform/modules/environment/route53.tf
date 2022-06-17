locals {
  r53_health_check_count = 1
  r53_hosted_zone_count  = 1
}

resource "aws_route53_health_check" "test" {
  count = local.r53_health_check_count

  fqdn              = "example.com"
  port              = 80
  type              = "HTTP"
  resource_path     = "/"
  failure_threshold = "5"
  request_interval  = "30"

  tags = {
    "test" : "route53-health-check-${count.index}"
  }
}

resource "aws_route53_zone" "test" {
  count = local.r53_hosted_zone_count

  name          = "${count.index}.example.com"
  force_destroy = true

  tags = {
    "test" : "route53-hosted-zone-${count.index}"
  }
}
