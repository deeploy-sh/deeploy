#!/usr/bin/env bash

# Exit on any error, undefined vars, or pipe fails
set -euo pipefail

# Must run on Linux
if [[ $(uname) != "Linux" ]]; then
    echo "Please run this script on Linux"
    exit 1
fi

# Must run with root privileges
if [[ $EUID -ne 0 ]]; then
    echo "Please run with sudo"
    exit 1
fi

# Get version arg or use 'latest'
VERSION=${1:-latest}
echo "🚀 Installing deeploy version: $VERSION"

# Check for Docker
if command -v docker &>/dev/null; then
    echo "✅ Docker already installed"
else
    echo "🐋 Installing Docker..."
    curl -fsSL https://get.docker.com | sudo bash
fi

# Handle Docker volume
VOLUME_NAME="deeploy_data"
if ! docker volume inspect "$VOLUME_NAME" &>/dev/null; then
    echo "📂 Creating Docker volume: $VOLUME_NAME"
    docker volume create "$VOLUME_NAME"
else
    echo "📂 Docker volume $VOLUME_NAME already exists"
fi

echo "📦 Starting deeploy..."

# Pull image and remove existing container
docker pull ghcr.io/deeploy-sh/deeploy:"$VERSION"
docker rm -f deeploy &>/dev/null || true

# Start container
docker run -d \
    --name deeploy \
    -p 8090:8090 \
    -v "$VOLUME_NAME":/app/data \
    ghcr.io/deeploy-sh/deeploy:"$VERSION"

# Get IP for display
IP=$(hostname -I | awk '{print $1}')
echo "✨ Deeploy ($VERSION) is running on http://$IP:8090"
