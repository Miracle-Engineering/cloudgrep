# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_node_group
# https://docs.aws.amazon.com/eks/latest/userguide/create-node-role.html
data "aws_iam_policy_document" "node_group" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      identifiers = ["ec2.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_iam_role" "node_group" {
  name               = "eks-cluster-node-group-role"
  assume_role_policy = data.aws_iam_policy_document.node_group.json
  tags = {
    test = "eks-cluster-default-nodegroup-iam-role"
  }
}

resource "aws_iam_role_policy_attachment" "node_group_AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.node_group.name
}

resource "aws_iam_role_policy_attachment" "node_group_AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.node_group.name
}

resource "aws_iam_role_policy_attachment" "node_group_AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.node_group.name
}

# https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html
data "aws_iam_policy_document" "trust_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      identifiers = ["eks.amazonaws.com"]
      type        = "Service"
    }
  }
}

resource "aws_iam_role" "cluster_role" {
  name               = "eks-cluster-default-role"
  assume_role_policy = data.aws_iam_policy_document.trust_policy.json
  tags = {
    test = "eks-cluster-default-iam-role"
  }
}

resource "aws_iam_role_policy_attachment" "cluster_AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.cluster_role.name
}
