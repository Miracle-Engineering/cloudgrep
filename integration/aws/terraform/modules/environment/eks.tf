module "eks" {
  source             = "../eks"
  id                 = "main"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  cluster_role       = "eks-cluster-default-role"
  node_group_role    = "eks-cluster-node-group-role"
}