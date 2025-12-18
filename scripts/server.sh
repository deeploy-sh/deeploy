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

# User parameter: branch, tag, or "latest" (default)
VERSION=${1:-latest}
echo "Installing deeploy server version: $VERSION"

# Map "latest" to "main" branch (latest is only a Docker tag alias, not a git branch)
BRANCH=${VERSION}
[[ "$VERSION" == "latest" ]] && BRANCH="main"

# Check for Docker
if command -v docker &>/dev/null; then
    echo "Docker already installed"
else
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com | sudo bash
fi

# Create install directory and subdirectories
INSTALL_DIR="/opt/deeploy"
mkdir -p "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR/traefik"
cd "$INSTALL_DIR"

# Download docker-compose.yml from same branch/tag as VERSION
echo "Downloading docker-compose.yml..."
curl -fsSL "https://raw.githubusercontent.com/deeploy-sh/deeploy/${BRANCH}/docker-compose.yml" \
  -o docker-compose.yml

# Generate secrets (only on first install)
if [[ ! -f .env ]]; then
    echo "Generating secrets..."
    JWT_SECRET=$(openssl rand -base64 32)
    ENCRYPTION_KEY=$(openssl rand -hex 16)  # 16 bytes = 32 hex chars
    cat > .env <<EOF
JWT_SECRET=$JWT_SECRET
ENCRYPTION_KEY=$ENCRYPTION_KEY
EOF
    chmod 600 .env
fi

# Start services (DEEPLOY_VERSION sets the image tag in docker-compose.yml)
echo "Starting deeploy..."
DEEPLOY_VERSION=$VERSION docker compose pull
DEEPLOY_VERSION=$VERSION docker compose up -d --force-recreate

IP=$(hostname -I | awk '{print $1}')
echo ""
echo "Deeploy is running!"
echo "  Dashboard: http://$IP:8090"
