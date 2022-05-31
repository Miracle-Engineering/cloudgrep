locals {
  rds_instance_count = 1
  rds_cluster_count = 1
}

resource "aws_db_instance" "test" {
  count = local.rds_instance_count

  identifier_prefix = "testing-${count.index}-"
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

resource "aws_rds_cluster" "postgresql" {
  count  = local.rds_cluster_count

  cluster_identifier_prefix = "testing-${count.index}-"
  engine = "aurora-postgresql"
  availability_zones = local.vpc_azs
  master_username = "postgres"
  master_password = "password"

  tags = {
    test = "rds-cluster-${count.index}"
  }
}
