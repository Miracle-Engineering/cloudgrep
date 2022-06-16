locals {
  iam_policy_count = 1
  iam_role_count   = 1
  iam_user_count   = 1

  permission_boundary_policy = "github-actions-terraform-permissions-boundary"
}

data "aws_iam_policy" "permissions_boundary" {
  name = local.permission_boundary_policy
}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "assume_role_self" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.id}:root"
      ]
    }
  }
}

resource "aws_iam_role" "test" {
  count = local.iam_role_count

  name_prefix           = "test-${count.index}-"
  path                  = "/test/"
  assume_role_policy    = data.aws_iam_policy_document.assume_role_self.json
  permissions_boundary  = data.aws_iam_policy.permissions_boundary.arn
  force_detach_policies = true

  tags = {
    "test" : "iam-role-${count.index}"
  }
}

resource "aws_iam_user" "test" {
  count = local.iam_user_count

  name                 = "test-${count.index}"
  path                 = "/test/"
  permissions_boundary = data.aws_iam_policy.permissions_boundary.arn
  force_destroy        = true

  tags = {
    "test" : "iam-user-${count.index}"
  }
}

resource "aws_iam_policy" "test" {
  count = local.iam_policy_count

  name_prefix = "test-${count.index}-"
  path        = "/test/"
  policy      = data.aws_iam_policy_document.test_policy.json
}

data "aws_iam_policy_document" "test_policy" {
  statement {
    actions   = ["sts:GetCallerIdentity"]
    resources = ["*"]
  }
}
