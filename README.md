# Basbas table reservations

## Development

    go run .

## Test

    go test ./...

## Deploy

    podman login rg.nl-ams.scw.cloud
    podman build -t rg.nl-ams.scw.cloud/basbas/basbas:latest .
    podman push rg.nl-ams.scw.cloud/basbas/basbas:latest
    cd terraform
    terraform init
    terraform apply
