# deeploy

Modern Deployment. Terminal First.

## Features

**Deployment**
- Zero-downtime deployments
- Git-based deployments with private repo support
- Real-time build logs

**Security**
- Zero-Config SSL - HTTPS automatic for all domains
- Auto-Provisioned Certificates via Let's Encrypt
- Secure by Default - HTTP redirects to HTTPS
- Auto-Renewal - Certificates renewed automatically

**Domains**
- Wildcard DNS via sslip.io - Instant domains without DNS config
- Custom domains with automatic SSL
- Multiple domains per pod

**Developer Experience**
- Terminal-first UI (TUI)
- Self-hosted - Full control over your infrastructure

## Install

### Server (VPS)

```bash
# Latest
curl -fsSL https://deeploy.sh/server.sh | sudo bash

# Specific version/branch
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s v1.0.0
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s dev

# Uninstall (dev/admin only - removes everything)
curl -fsSL https://raw.githubusercontent.com/deeploy-sh/deeploy/main/scripts/uninstall.sh | sudo bash
```

### TUI (Local)

```bash
curl -fsSL https://deeploy.sh/tui.sh | bash
```

## Development

```bash
# Server (VPS daemon)
task dev:server

# TUI
task dev:tui
```

## Contributing

Please read the [contributing guide](CONTRIBUTING.md).

## License

[Apache 2.0](LICENSE)
