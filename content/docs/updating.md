---
title: Updating
description: Update deeploy to the latest version
order: 7
---

## Update Server

SSH into your VPS and run the install script again:

```bash
curl -fsSL https://deeploy.sh/server.sh | sudo bash
```

This pulls the latest Docker image and restarts the server. Your data is preserved.

### Specific Version

```bash
curl -fsSL https://deeploy.sh/server.sh | sudo bash -s v0.2.0
```

## Update TUI

Run the install script again on your machine:

```bash
curl -fsSL https://deeploy.sh/tui.sh | bash
```

### Specific Version

```bash
curl -fsSL https://deeploy.sh/tui.sh | bash -s v0.2.0
```

## Check Version

The TUI shows your current version and the latest available version in the status bar. You can also check via command palette (`Alt+P`) → "Info".
