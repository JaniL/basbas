output "container_url" {
  description = "Public URL of the deployed app"
  value       = replace(scaleway_container.basbas.public_endpoint, "/^([^\\/]+)$/", "https://$1")
}

output "registry_endpoint" {
  description = "Registry endpoint for podman login and push"
  value       = scaleway_registry_namespace.basbas.endpoint
}
