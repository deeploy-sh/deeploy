# deeploy

[![GitHub Sponsors](https://img.shields.io/github/sponsors/axadrn?style=social&label=♥️%20Sponsor)](https://github.com/sponsors/axadrn)

Modern Deployment. Terminal First.

Deploy with a sleek TUI. For developers who live in the terminal.

## Features

- Modern Terminal UI AND CLI commands
- Docker-based deployments
- Open source and self-hosted
- Built with Go

## Why deeploy?

- First deployment platform with a proper Terminal UI (TUI first approach)
- Not just another web dashboard (but we'll have that too)
- Simple, fast, and dev-friendly
- Built with Go, not PHP or JS

## Quick Start

```bash
# Install server on your VPS (Hetzner, DigitalOcean, etc.)
curl -fsSL https://deeploy.sh/install.sh | sh

# Install CLI/TUI client on your local machine
curl -fsSL deeploy.sh/install-cli | sh

# Start deeploy
deeploy
```

## Usage

```bash
# Start the Terminal UI
deeploy
```

## Current Status

This is a pre-alpha release. The platform is under active development with upcoming features including:

- User Authentication
- Project Management
- Container Deployments
- Domain Management
- Templates
- And more!

## Requirements

- Server: Linux VPS (Hetzner, DigitalOcean, etc.) with Docker
- Client: Any machine for TUI client (macOS, Linux, Windows with WSL)

## Built With

- Go
- Bubbletea (TUI)
- SQLite
- Templ + templUI (Web UI coming soon)
- Docker

## Contributing

We welcome contributions from the community! Whether it's adding new features, improving existing ones, or enhancing documentation, your input is valuable. Please check our [contributing guidelines](CONTRIBUTING.md) for more information on how to get involved.

## License

Deeploy is open-source software licensed under the [MIT license](LICENSE).

## Support

For support, questions, or discussions, please [open an issue](https://github.com/deeploy-sh/deeploy/issues) on our GitHub repository or [visit our community (GitHub Discussions)](https://github.com/deeploy-sh/deeploy/discussions).

---

Built with ❤️ by the dev community, for the dev community.
