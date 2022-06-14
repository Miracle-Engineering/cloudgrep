data "aws_caller_identity" "current" {}

data "aws_s3_bucket" "state_storage" {
  bucket = "${data.aws_caller_identity.current.account_id}-terraform-state"
}

data "aws_dynamodb_table" "state_storage" {
  name = "terraform-locks"
}
