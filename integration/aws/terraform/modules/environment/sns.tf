locals {
  sns_count = 2
}

resource "random_string" "sns_topic_name" {
  count = local.sns_count

  length  = 8
  special = false
  upper   = false
}

resource "aws_sns_topic" "topic" {
  count           = local.sns_count
  name            = "testing--${count.index}-${random_string.sns_topic_name[count.index].result}"
  delivery_policy = <<EOF
{
  "http": {
    "defaultHealthyRetryPolicy": {
      "minDelayTarget": 20,
      "maxDelayTarget": 20,
      "numRetries": 3,
      "numMaxDelayRetries": 0,
      "numNoDelayRetries": 0,
      "numMinDelayRetries": 0,
      "backoffFunction": "linear"
    },
    "disableSubscriptionOverrides": false,
    "defaultThrottlePolicy": {
      "maxReceivesPerSecond": 1
    }
  }
}
EOF
  lifecycle {
    ignore_changes = [name, name_prefix]
  }
  tags = {
    test : "sns-topic-${count.index}"
    IntegrationTest : "true"
  }
}

## SNS topic policy
resource "aws_sns_topic_policy" "default" {
  count = local.sns_count
  arn   = aws_sns_topic.topic[count.index].arn

  policy = data.aws_iam_policy_document.sns_topic_policy.json
}

data "aws_iam_policy_document" "sns_topic_policy" {
  count     = local.sns_count
  policy_id = "__default_policy_ID-${count.index}"

  statement {
    actions = [
      "SNS:Subscribe",
      "SNS:SetTopicAttributes",
      "SNS:RemovePermission",
      "SNS:Receive",
      "SNS:Publish",
      "SNS:ListSubscriptionsByTopic",
      "SNS:GetTopicAttributes",
      "SNS:DeleteTopic",
      "SNS:AddPermission",
    ]

    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }

    resources = [
      aws_sns_topic.topic[count.index].arn,
    ]

    sid = "__default_statement_ID-${count.index}"
  }
}
