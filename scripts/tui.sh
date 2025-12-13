#!/usr/bin/env bash
set -euo pipefail

VERSION=${1:-latest}
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize arch
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

FILENAME="deeploy-${OS}-${ARCH}"

if [[ "$VERSION" == "latest" ]]; then
    URL="https://github.com/deeploy-sh/deeploy/releases/latest/download/${FILENAME}"
else
    URL="https://github.com/deeploy-sh/deeploy/releases/download/${VERSION}/${FILENAME}"
fi

echo "Downloading deeploy TUI ($VERSION)..."
curl -fsSL -o deeploy "$URL"
chmod +x deeploy

echo "Installing to /usr/local/bin..."
sudo mv deeploy /usr/local/bin/deeploy

echo "Installed! Run 'deeploy' to start."
