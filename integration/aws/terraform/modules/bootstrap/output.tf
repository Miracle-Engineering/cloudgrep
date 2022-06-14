output "github_actions_iam_role" {
  value = aws_iam_role.gha_tests.arn
}
