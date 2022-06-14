locals {
  lambda_function_count = 2
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/lambda.js"
  output_path = "${path.module}/lambda.zip"
}

resource "random_string" "lambda_test" {
  count = local.lambda_function_count

  length  = 8
  special = false
  upper   = false
}

resource "aws_lambda_function" "test" {
  count         = local.lambda_function_count
  function_name = "testing-${count.index}-${random_string.lambda_test[count.index].result}"
  role          = data.aws_iam_role.test_lambda.arn

  runtime          = "nodejs16.x"
  handler          = "lambda.handler"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  tags = {
    test = "lambda-function-${count.index}"
  }
}

data "aws_iam_role" "test_lambda" {
  name = "test-lambda-execution-role"
}
