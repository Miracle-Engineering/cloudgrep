locals {
  cf_count = 1
}

# This is optional, see Opta docs for this here: https://docs.opta.dev/reference/aws/modules/cloudfront-distribution/
#tfsec:ignore:aws-cloudfront-enable-waf
resource "aws_cloudfront_distribution" "distribution" {
  count = local.cf_count

  comment         = "Test cloudfront distribution ${count.index}"
  enabled         = true
  is_ipv6_enabled = true
  price_class     = "PriceClass_200"
  aliases         = []

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
    ssl_support_method             = "sni-only"
  }

  origin {
    domain_name = aws_lb.alb[0].dns_name
    origin_id   = "DefaultLbOriginId"
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2", "TLSv1.1", "TLSv1"]
    }
  }

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD", "OPTIONS"]
    target_origin_id       = "DefaultLbOriginId"
    viewer_protocol_policy = "redirect-to-https"

    forwarded_values {
      query_string = true
      headers      = ["Origin", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Host"]
      cookies {
        forward = "all"
      }
    }
  }

  tags = {
    test = "cloudfront-distribution-${count.index}"
  }
}
