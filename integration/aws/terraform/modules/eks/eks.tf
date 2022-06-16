resource "aws_security_group" "eks" {
  name_prefix = "eks-cluster-${var.id}"
  description = "EKS cluster security group."
  vpc_id      = var.vpc_id

  tags = {
    test = "eks-cluster-${var.id}-sg"
  }

  ingress {
    description = "allowallfromself"
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
    self        = true
  }

  egress {
    description = "alloutbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    # To be fixed - for now user can create an SG manually to override.
    #tfsec:ignore:aws-vpc-no-public-egress-sgr
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_eks_cluster" "cluster" {
  name     = var.id
  role_arn = aws_iam_role.cluster_role.arn
  version  = var.k8s_version

  # To be fixed, although early opta users need this - Cluster allows access from a public CIDR: 0.0.0.0/0
  #tfsec:ignore:aws-eks-no-public-cluster-access-to-cidr
  vpc_config {
    subnet_ids              = var.private_subnet_ids
    security_group_ids      = [aws_security_group.eks.id]
    endpoint_private_access = false # TODO: make this true once we got VPN figured out
    #tfsec:ignore:aws-eks-no-public-cluster-access
    endpoint_public_access = true # TODO: make this false once we got VPN figured out
  }

  # https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  # Ensure that IAM Role permissions are created before and deleted after EKS Cluster handling.
  # Otherwise, EKS will not be able to properly delete EKS managed EC2 infrastructure such as Security Groups.
  depends_on = [
    aws_iam_role_policy_attachment.cluster_AmazonEKSClusterPolicy,
    aws_cloudwatch_log_group.cluster_logs,
  ]

  tags = {
    test = "eks-cluster-${var.id}"
  }
  lifecycle {
    ignore_changes = [vpc_config[0].security_group_ids]
  }
}

resource "aws_cloudwatch_log_group" "cluster_logs" {
  name              = "/aws/eks/${var.id}/cluster"
  retention_in_days = var.eks_log_retention
  tags = {
    test = "eks-cluster-${var.id}-log-group"
  }
  lifecycle {
    ignore_changes = [name]
  }
}
