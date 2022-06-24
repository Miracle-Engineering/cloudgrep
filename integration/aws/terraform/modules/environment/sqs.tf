locals {
  sqs_queue_count = 1
}

resource "aws_sqs_queue" "sqs_queue" {
  count       = local.sqs_queue_count
  name_prefix = "testing-${count.index}"
  tags = {
    test : "sqs-queue-${count.index}"
  }
}