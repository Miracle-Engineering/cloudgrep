resource "random_id" "key_suffix" {
  byte_length = 8
  keepers = {
    spot_instances     = var.spot_instances
    node_disk_size     = var.node_disk_size
    node_instance_type = var.node_instance_type
  }
}


resource "aws_eks_node_group" "node_group" {
  cluster_name    = aws_eks_cluster.cluster.name
  node_group_name = "${var.id}-default-${random_id.key_suffix.hex}"
  node_role_arn   = data.aws_iam_role.node_group.arn
  subnet_ids      = aws_eks_cluster.cluster.vpc_config[0].subnet_ids
  capacity_type   = var.spot_instances ? "SPOT" : "ON_DEMAND"


  ami_type = var.ami_type

  disk_size      = var.node_disk_size
  instance_types = [var.node_instance_type]
  labels         = { node_group_name = "${var.id}-default" }

  scaling_config {
    max_size     = var.max_nodes
    desired_size = max(var.min_nodes, 1)
    min_size     = var.min_nodes
  }

  # Optional: Allow external changes without Terraform plan difference
  lifecycle {
    ignore_changes        = [scaling_config[0].desired_size, node_group_name, subnet_ids]
    create_before_destroy = true
  }

  tags = {
    test = "eks-cluster-${var.id}-default-node-group"
  }
}
