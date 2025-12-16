# deeploy

Modern Deployment. Terminal First.

## Features

- Zero-downtime deployments
- Terminal-first UI
- Git-based deployments
- Automatic HTTPS via Traefik

## Install

### Server (VPS)

```bash
# Latest
curl -fsSL https://deeploy.sh/server.sh | sudo bash

# Specific version/branch
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s v1.0.0
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s dev
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
