variable "container_image" {
  description = "Fully qualified image reference pushed to the Scaleway registry"
  type        = string
  default     = "rg.nl-ams.scw.cloud/basbas/basbas:latest"
}

variable "max_scale" {
  description = "Maximum number of concurrently running container instances"
  type        = number
  default     = 5
}
