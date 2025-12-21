# deeploy

Modern Deployment. Terminal First. Open Source.

<video src="https://github.com/deeploy-sh/deeploy/raw/main/internal/shared/assets/video/hero.mp4" autoplay loop muted playsinline></video>

## Quick Start

### Server (VPS)

```bash
# Stable release (recommended)
curl -fsSL https://deeploy.sh/server.sh | sudo bash

# Bleeding edge
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s main

# Specific version, branch, or tag
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s v0.1.0
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s feature/cool-thing
```

### TUI (your machine)

```bash
# Stable release (recommended)
curl -fsSL https://deeploy.sh/tui.sh | bash

# Specific version or tag
curl -fsSL https://deeploy.sh/tui.sh | bash -s v0.1.0
```

## Features

- Zero-downtime deployments
- Auto SSL via Let's Encrypt
- Instant domains via sslip.io wildcard DNS
- Terminal-first UI (TUI)
- Self-hosted, you own your data

## Requirements

**Server (VPS)**
- Any Linux VPS (Ubuntu, Debian, etc.)
- Docker installed (script installs if missing)
- 1GB+ RAM recommended
- Ports 80, 443 open

**TUI**
- macOS or Linux (Windows via WSL)
- That's it

## Currently Supported

- GitHub repositories (public & private)
- Docker-based deployments
- Custom domains + SSL
- Single-user mode

## Coming in 0.2+

- GitLab & Bitbucket support
- Multi-user / Teams
- Webhooks & CI integration
- Rollback
- Resource limits
- And many more...

## Development

```bash
# Server (VPS daemon)
task dev:server

# TUI
task dev:tui

# Docs
task dev:docs
```

## License

[Apache 2.0](LICENSE)
