module "eks" {
  source             = "../eks"
  id                 = "main"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
}