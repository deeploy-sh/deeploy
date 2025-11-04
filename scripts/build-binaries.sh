#!/bin/bash

# https://akrabat.com/building-go-binaries-for-different-platforms/
# https://akrabat.com/setting-the-version-of-a-go-application-when-building/

version=`git describe --tags HEAD`

platforms=(
"darwin/amd64"
"darwin/arm64"
"linux/amd64"
"linux/arm64"
)

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name="deeploy-${GOOS}-${GOARCH}"
    echo "Building release/$output_name..."
    env GOOS=$GOOS GOARCH=$GOARCH go build \
      -C app/cmd/tui \
      -ldflags "-X github.com/axadrn/deeploy/commands.Version=$version" \
      -o release/$output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting.'
        exit 1
    fi
done
