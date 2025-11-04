#!/bin/bash

echo 'Building local deeployd binary...'
GOOS=linux GOARCH=arm64 go build -o deeployd ./cmd/deeployd

echo 'Moving local deeploy binary to remote host...'
scp deeployd root@95.216.205.0:/usr/local/bin/deeployd

echo 'Finished!'
