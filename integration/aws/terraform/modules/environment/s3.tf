locals {
  s3_bucket_count = 2
}

resource "aws_s3_bucket" "test" {
  count         = local.s3_bucket_count
  bucket_prefix = "cloudgrep-testing-${count.index}-"

  force_destroy = true

  tags = {
    test = "s3-bucket-${count.index}"
  }
}
