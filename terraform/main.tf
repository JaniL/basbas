terraform {
  required_version = ">= 1.5"

  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "~> 2.0"
    }
  }
}

provider "scaleway" {
  region = "nl-ams"
  zone   = "nl-ams-1"
}

resource "scaleway_container_namespace" "basbas" {
  name        = "basbas"
  description = "BasBas table reservations"
}

resource "scaleway_registry_namespace" "basbas" {
  name        = "basbas"
  description = "Images for the BasBas table reservations app"
  is_public   = false
}

resource "scaleway_container" "basbas" {
  namespace_id = scaleway_container_namespace.basbas.id
  name         = "basbas"
  image        = var.container_image
  port         = 8080
  protocol     = "http1"
  privacy      = "public"

  https_connections_only = true

  min_scale          = 0
  max_scale          = var.max_scale
  memory_limit_bytes = 134217728
  cpu_limit          = 140
  timeout            = 30

  environment_variables = {
    PORT = "8080"
  }
}
