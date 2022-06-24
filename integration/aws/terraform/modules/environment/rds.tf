locals {
  rds_instance_count         = 1
  rds_cluster_count          = 1
  rds_snapshot_count         = 1
  rds_cluster_snapshot_count = 1
}

resource "aws_db_instance" "test" {
  count = local.rds_instance_count

  identifier_prefix    = "test-${count.index}-"
  allocated_storage    = 10
  engine               = "postgres"
  engine_version       = "13.4"
  instance_class       = "db.t3.micro"
  db_subnet_group_name = aws_db_subnet_group.private.name
  username             = "postgres"
  password             = "password"
  parameter_group_name = "default.postgres13"
  skip_final_snapshot  = true

  tags = {
    test = "rds-instance-${count.index}"
  }
}

resource "aws_rds_cluster" "postgresql" {
  count = local.rds_cluster_count

  cluster_identifier_prefix = "test-${count.index}-"
  engine                    = "aurora-postgresql"
  db_subnet_group_name      = aws_db_subnet_group.private.name
  availability_zones        = module.vpc.private_subnet_azs
  master_username           = "postgres"
  master_password           = "password"
  skip_final_snapshot       = true

  tags = {
    test = "rds-cluster-${count.index}"
  }
}

resource "aws_db_subnet_group" "private" {
  name_prefix = "test-"
  subnet_ids  = module.vpc.private_subnet_ids
}

resource "random_string" "db_snapshot_suffix" {
  count = local.rds_snapshot_count

  length  = 8
  special = false
  upper   = false
}

resource "aws_db_snapshot" "test" {
  count = local.rds_snapshot_count

  db_instance_identifier = aws_db_instance.test[count.index].id
  db_snapshot_identifier = "test-${count.index}-${random_string.db_snapshot_suffix[count.index].id}"

  tags = {
    "test" : "rds-snapshot-${count.index}"
  }
}


resource "random_string" "db_cluster_snapshot_suffix" {
  count = local.rds_cluster_snapshot_count

  length  = 8
  special = false
  upper   = false
}

resource "aws_db_cluster_snapshot" "test" {
  count = local.rds_cluster_snapshot_count

  db_cluster_identifier          = aws_rds_cluster.postgresql[count.index].id
  db_cluster_snapshot_identifier = "test-cluster-${count.index}-${random_string.db_snapshot_suffix[count.index].id}"

  tags = {
    "test" : "rds-cluster-snapshot-${count.index}"
  }
}
