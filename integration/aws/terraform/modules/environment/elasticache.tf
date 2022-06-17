locals {
  redis_cluster_count = 1
}

resource "random_password" "redis_auth" {
  length  = 20
  special = false
}

resource "random_string" "redis_name_hash" {
  length  = 4
  special = false
  upper   = false
}


resource "aws_security_group" "elasticache" {
  name        = "test-elasticache-sg"
  description = "For usage by elasticache to give access to resources in the vpc"
  vpc_id      = module.vpc.id

  # https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-ug/elasticache-vpc-accessing.html
  ingress {
    description = "redis"
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr]
  }

}

resource "aws_elasticache_replication_group" "redis_cluster" {
  count                         = local.redis_cluster_count
  automatic_failover_enabled    = true
  auto_minor_version_upgrade    = true
  security_group_ids            = [aws_security_group.elasticache.id]
  subnet_group_name             = "test-${count.index}"
  replication_group_id          = "test-${count.index}-${random_string.redis_name_hash.result}"
  replication_group_description = "Elasticache test test-${count.index}-${random_string.redis_name_hash.result}"
  node_type                     = "cache.m4.large"
  engine_version                = "6.x"
  number_cache_clusters         = 2
  port                          = 6379
  apply_immediately             = true
  multi_az_enabled              = true
  auth_token                    = random_password.redis_auth.result
  transit_encryption_enabled    = true
  at_rest_encryption_enabled    = true
  snapshot_window               = "04:00-05:00"
  snapshot_retention_limit      = 0
  lifecycle {
    ignore_changes = [engine_version, replication_group_id, replication_group_description]
  }
}