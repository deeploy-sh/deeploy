---
title: Deploying
description: Deploy your first app
order: 3
---

Deploy your first application step by step.

## 1. Create a Project

Projects are containers that organize your pods. Open the command palette (`Alt+P`) and select "New Project".

## 2. Create a Pod

A pod is a single deployable application. One pod = one container.

**Required fields:**
- **Title** - A name for your pod
- **Repository URL** - Your Git repository (e.g., `https://github.com/user/repo`)
- **Branch** - The branch to deploy (default: `main`)
- **Dockerfile** - Path to your Dockerfile (default: `Dockerfile`)

For private repositories, you'll need to add a [Git Token](/docs/private-repos) first.

## 3. Add a Domain

Your pod needs a domain to be accessible. You have two options:

- **Auto-generated** - Instant subdomain like `pod-abc123.1.2.3.4.sslip.io`
- **Custom** - Your own domain like `myapp.example.com` (requires [DNS setup](/docs/domains))

## 4. Deploy

Hit "Deploy" and watch the build logs. The process:

1. Clone your repository
2. Build the Docker image
3. Start the container
4. Route traffic via Traefik

Once complete, your app is live at the domain URL.

## Managing Your Pod

**Stop** - Stops the container (keeps configuration)

**Restart** - Restarts the container

**Logs** - View container output

**Redeploy** - Pull latest code and rebuild

To update your app, just push to your repository and hit "Deploy" again.
