module "eks" {
  source          = "terraform-aws-modules/eks/aws"
  cluster_name    = "url-shortner"
  cluster_version = "1.21"
  subnets         = module.vpc.private_subnets
  vpc_id          = module.vpc.vpc_id

  worker_groups = [
    {
      instance_type = "m5.large"
      asg_max_size  = 3
    }
  ]

  tags = {
    "Environment" = "dev"
  }
}
