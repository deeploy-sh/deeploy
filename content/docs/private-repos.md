---
title: Private Repositories
description: Deploy from private Git repositories
order: 5
---

To deploy from private repositories, you need to add a Git token.

## Create a GitHub Token

1. Go to [GitHub Settings → Developer Settings → Personal Access Tokens](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Select scope: **repo** (full access to private repositories)
4. Generate and copy the token

## Add Token to Deeploy

1. Open the command palette (`Alt+P`)
2. Select "Git Tokens" → "New"
3. Enter a name (e.g., "GitHub Personal")
4. Paste your token

Your token is encrypted before being stored.

## Use with a Pod

When creating or editing a pod:

1. Enter your private repository URL
2. Select the Git Token you created
3. Deploy

The token is used to clone your repository during the build process.

## Token Security

- Tokens are encrypted with AES-256 before storage
- Only used during git clone operations
- Never exposed in logs or UI
- Delete tokens anytime from the Git Tokens menu
