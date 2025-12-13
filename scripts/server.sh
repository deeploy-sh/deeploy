#!/usr/bin/env bash
set -euo pipefail

if [[ $(uname) != "Linux" ]]; then
    echo "Please run this script on Linux"
    exit 1
fi

if [[ $EUID -ne 0 ]]; then
    echo "Please run with sudo"
    exit 1
fi

VERSION=${1:-latest}
echo "Installing deeploy server version: $VERSION"

# Check for Docker
if command -v docker &>/dev/null; then
    echo "Docker already installed"
else
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com | sudo bash
fi

# Create install directory
INSTALL_DIR="/opt/deeploy"
mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

# Download docker-compose.yml
echo "Downloading docker-compose.yml..."
curl -fsSL "https://raw.githubusercontent.com/deeploy-sh/deeploy/main/docker-compose.yml" \
  -o docker-compose.yml

# Start services with specified version
echo "Starting deeploy..."
DEEPLOY_VERSION=$VERSION docker compose pull
DEEPLOY_VERSION=$VERSION docker compose up -d

IP=$(hostname -I | awk '{print $1}')
echo ""
echo "Deeploy is running!"
echo "  Dashboard: http://$IP:8090"
