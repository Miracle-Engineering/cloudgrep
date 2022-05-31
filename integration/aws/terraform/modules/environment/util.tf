resource "aws_default_vpc" "default" {}

resource "aws_default_security_group" "default" {
    vpc_id = aws_default_vpc.default.id
}

data "aws_subnets" "default" {
    filter {
        name = "vpc-id"
        values = [aws_default_vpc.default.id]
    }
}
