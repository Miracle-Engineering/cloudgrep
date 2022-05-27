locals {
  rds_instance_count = 1
}

resource "aws_db_instance" "test" {
  count = 1

  allocated_storage    = 10
  engine               = "postgres"
  engine_version       = "13.4"
  instance_class       = "db.t3.micro"
  username             = "postgres"
  password             = "password"
  parameter_group_name = "default.postgres13"
  skip_final_snapshot  = true

  tags = {
    test = "rds-instance-${count.index}"
  }
}
