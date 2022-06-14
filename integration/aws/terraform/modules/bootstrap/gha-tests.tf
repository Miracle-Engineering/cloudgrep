resource "aws_iam_role" "gha_tests" {
  name               = "github-actions-tests"
  assume_role_policy = data.aws_iam_policy_document.gha_tests_assume_role.json

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/ReadOnlyAccess",
  ]
}

data "aws_iam_policy_document" "gha_tests_assume_role" {
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
  }
}
