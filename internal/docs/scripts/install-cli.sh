#!/bin/bash

platforms=(
"darwin/amd64"
"darwin/arm64"
"linux/amd64"
"linux/arm64"
)

filename=""
uname=$(tr "[:upper:]" "[:lower:]" <<< $(uname))

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    OS=${platform_split[0]}
    ARCH=${platform_split[1]}
    if [[ "$uname" == "$OS" && $(uname -m) == "$ARCH" ]]; then
      filename="deeploy-${OS}-${ARCH}"
      break
    fi
done

downloadURL="https://github.com/deeploy-sh/deeploy/releases/latest/download/${filename}"
appname="deeploy"
curl -Lo "$appname" "$downloadURL"
chmod +x "$appname"
sudo mv "$appname" /usr/local/bin/deeploy
