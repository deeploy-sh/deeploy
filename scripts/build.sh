#!/bin/bash
set -euo pipefail

# Build TUI binaries for all platforms
# https://akrabat.com/building-go-binaries-for-different-platforms/

VERSION=${1:-$(git describe --tags HEAD 2>/dev/null || echo "dev")}
OUTPUT_DIR="release"

platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
)

mkdir -p "$OUTPUT_DIR"

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    output_name="deeploy-${GOOS}-${GOARCH}"

    echo "Building $output_name..."
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.Version=$VERSION" \
        -o "$OUTPUT_DIR/$output_name" \
        ./cmd/deeploy
done

echo "Done! Binaries in $OUTPUT_DIR/"
