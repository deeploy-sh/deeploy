---
title: Getting Started
description: Install deeploy in minutes
order: 2
---

Get up and running with deeploy in minutes.

## Requirements

**Server (VPS)**
- Any Linux VPS (Ubuntu, Debian, etc.)
- Docker installed (script installs if missing)
- 1GB+ RAM recommended
- Ports 80, 443 open

**TUI (your machine)**
- macOS or Linux (Windows via WSL)

## Install Server

SSH into your VPS and run:

```bash
curl -fsSL https://deeploy.sh/server.sh | sudo bash
```

### Options

```bash
# Bleeding edge
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s main

# Specific version
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s v0.1.0
```

## Install TUI

On your local machine:

```bash
curl -fsSL https://deeploy.sh/tui.sh | bash
```

### Options

```bash
# Specific version
curl -fsSL https://deeploy.sh/tui.sh | bash -s v0.1.0
```

## Next Steps

Run `deeploy` to launch the TUI and connect to your server.
