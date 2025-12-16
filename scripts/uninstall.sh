#!/usr/bin/env bash
set -euo pipefail

# Uninstall script for deeploy server
# Removes all containers, images, volumes, and data for a clean slate

if [[ $(uname) != "Linux" ]]; then
    echo "Please run this script on Linux"
    exit 1
fi

if [[ $EUID -ne 0 ]]; then
    echo "Please run with sudo"
    exit 1
fi

echo "=== Deeploy Uninstall ==="
echo "This will remove ALL deeploy data including:"
echo "  - All deployed pods (deeploy-* containers)"
echo "  - Deeploy stack (postgres, traefik, app)"
echo "  - All images"
echo "  - Database data"
echo ""
read -p "Are you sure? (y/N) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi

echo ""
echo "Stopping and removing all deeploy containers..."
# Stop all containers with deeploy- prefix (deployed pods + stack)
docker ps -a --filter "name=deeploy-" --format "{{.ID}}" | xargs -r docker stop 2>/dev/null || true
docker ps -a --filter "name=deeploy-" --format "{{.ID}}" | xargs -r docker rm -f 2>/dev/null || true

echo "Removing deeploy images..."
# Remove all images with deeploy- prefix (built pod images)
docker images --filter "reference=deeploy-*" --format "{{.ID}}" | xargs -r docker rmi -f 2>/dev/null || true
# Remove stack images
docker rmi -f ghcr.io/deeploy-sh/deeploy 2>/dev/null || true
docker rmi -f postgres:16-alpine 2>/dev/null || true
docker rmi -f traefik:v3.2 2>/dev/null || true

echo "Removing docker volume..."
docker volume rm deeploy_postgres_data 2>/dev/null || true
docker volume rm postgres_data 2>/dev/null || true

echo "Removing docker network..."
docker network rm deeploy 2>/dev/null || true

echo "Removing install directory..."
rm -rf /opt/deeploy

echo "Pruning dangling images and build cache..."
docker image prune -f 2>/dev/null || true
docker builder prune -f 2>/dev/null || true

echo ""
echo "Done! VPS is clean."
echo "Run server.sh to reinstall fresh."
