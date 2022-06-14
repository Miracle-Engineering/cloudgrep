locals {
  gh_actions_oidc_url        = "https://token.actions.githubusercontent.com"
  gh_actions_oidc_config_url = "${local.gh_actions_oidc_url}/.well-known/openid-configuration"
  gh_actions_oidc_config     = jsondecode(data.http.github_actions_oidc_config.body)
}

data "http" "github_actions_oidc_config" {
  url = local.gh_actions_oidc_config_url
}

data "tls_certificate" "github_actions" {
  url = local.gh_actions_oidc_config.jwks_uri
}

resource "aws_iam_openid_connect_provider" "github_actions" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.github_actions.certificates[0].sha1_fingerprint]
}

data "aws_arn" "github_actions_oidc_provider" {
  arn = aws_iam_openid_connect_provider.github_actions.arn
}

locals {
  gh_actions_oidc_provider_id = split("/", data.aws_arn.github_actions_oidc_provider.resource)[1]
}
