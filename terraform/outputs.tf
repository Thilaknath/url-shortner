output "cluster_name" {
  value = module.eks.cluster_id
}

output "kubeconfig" {
  value = module.eks.kubeconfig
  sensitive = true
}

output "redis_endpoint" {
  value = aws_elasticache_cluster.example.cache_nodes.0.address
}

output "redis_port" {
  value = aws_elasticache_cluster.example.port
}
