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

# Resolve "latest" to the actual latest release tag from GitHub
# For bleeding-edge, use: ./server.sh main
BRANCH=${VERSION}
if [[ "$VERSION" == "latest" ]]; then
    echo "Fetching latest stable release..."
    BRANCH=$(curl -sL https://api.github.com/repos/deeploy-sh/deeploy/releases/latest 2>/dev/null | grep '"tag_name"' | cut -d'"' -f4 || echo "")
    if [[ -z "$BRANCH" ]]; then
        echo "No releases found, using main branch"
        BRANCH="main"
    fi
fi

echo "Installing deeploy server version: $BRANCH"

# Docker tags can't contain slashes - convert feat/example → feat-example
# (same logic as .github/workflows/ci.yml)
TAG=${BRANCH//\//-}

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
DEEPLOY_VERSION=$TAG docker compose pull
DEEPLOY_VERSION=$TAG docker compose up -d --force-recreate

IP=$(hostname -I | awk '{print $1}')
echo ""
echo "Deeploy is running!"
echo "  Dashboard: http://$IP:8090"
