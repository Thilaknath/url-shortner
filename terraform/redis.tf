resource "aws_elasticache_subnet_group" "example" {
  name       = "example"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Name = "example"
  }
}

resource "aws_elasticache_cluster" "example" {
  cluster_id           = "example"
  engine               = "redis"
  node_type            = "cache.t2.micro"  # cost-effective option
  num_cache_nodes      = 1
  parameter_group_name = "default.redis6.x"
  subnet_group_name    = aws_elasticache_subnet_group.example.name
  security_group_ids   = [module.vpc.default_security_group_id]

  tags = {
    Name = "example"
  }
}

output "redis_endpoint" {
  value = aws_elasticache_cluster.example.cache_nodes.0.address
}

output "redis_port" {
  value = aws_elasticache_cluster.example.port
}
