data "aws_iam_policy_document" "lambda_iam_assume_role" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "test_lambda_execution_role" {
  name = "test-lambda-execution-role"

  assume_role_policy = data.aws_iam_policy_document.lambda_iam_assume_role.json
}
