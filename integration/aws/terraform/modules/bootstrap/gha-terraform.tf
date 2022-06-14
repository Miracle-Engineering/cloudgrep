resource "aws_iam_role" "gha_terraform" {
  name               = "github-actions-terraform"
  assume_role_policy = data.aws_iam_policy_document.gha_terraform_assume_role.json

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/ReadOnlyAccess",
  ]

  inline_policy {
    name   = "terraform-state"
    policy = data.aws_iam_policy_document.gha_terraform.json
  }

  inline_policy {
    name   = "resources"
    policy = data.aws_iam_policy_document.gha_terraform_resources.json
  }
}

data "aws_iam_policy_document" "gha_terraform_assume_role" {
  statement {
    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github_actions.arn]
    }

    actions = ["sts:AssumeRoleWithWebIdentity"]
    condition {
      test     = "StringEquals"
      variable = "${local.gh_actions_oidc_provider_id}:aud"
      values   = ["sts.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "${local.gh_actions_oidc_provider_id}:sub"
      values   = ["repo:${var.trusted_repo}:*"]
    }

    # FIXME(patrick): This condition is not matching the actions workflow as expected, even though the claims appear correct
    # condition {
    #   test     = "StringEquals"
    #   variable = "${local.gh_actions_oidc_provider_id}:workflow"
    #   values   = ["Integration Test Setup"]
    # }
  }
}

data "aws_iam_policy_document" "gha_terraform" {
  statement {
    actions   = ["s3:ListBucket"]
    resources = [data.aws_s3_bucket.state_storage.arn]
  }

  statement {
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject"
    ]
    resources = ["${data.aws_s3_bucket.state_storage.arn}/environment/tfstate"]
  }

  statement {
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:DeleteItem"
    ]
    resources = [data.aws_dynamodb_table.state_storage.arn]
  }

  statement {
    effect = "Deny"
    actions = [
      "s3:DeleteBucket*",
      "s3:PutBucket*",
      "s3:*PublicAccessBlock",
    ]
    resources = [data.aws_s3_bucket.state_storage.arn]
  }

  statement {
    effect = "Deny"
    actions = [
      "dynamodb:*Table",
    ]
    resources = [data.aws_dynamodb_table.state_storage.arn]
  }
}

data "aws_iam_policy_document" "gha_terraform_resources" {
  statement {
    effect    = "Deny"
    resources = ["*"]
    actions = [
      "ec2:*CapacityReservation", // Block access to the reservation APIs
      "ec2:Purchase*",
      "rds:Purchase*",
      "cloudtrail:*",       // Prevent turning of any cloudtrail logs
      "cloudwatch:*Alarm*", // Protect any billing alarms
    ]
  }

  statement {
    actions = [
      "autoscaling:*",
      "ec2:*",
      "elasticloadbalancing:*",
      "iam:CreateServiceLinkedRole",
      "lambda:*",
      "rds:*",
      "s3:*",
    ]
    resources = ["*"]
  }
}
