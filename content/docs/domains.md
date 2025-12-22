---
title: Domains
description: Custom domains and HTTPS setup
order: 4
---

## Server Domain

**Important:** Without a custom server domain, all communication between the TUI and server happens over HTTP. This means your password and data are sent in plaintext. Set up a domain immediately after installation.

### Setup

1. **Point your domain to your server**

   Add an A record in your DNS provider:
   ```
   Type: A
   Name: @ (or subdomain like "deploy")
   Value: YOUR_SERVER_IP
   ```

2. **Set the domain in deeploy**

   Open the command palette (`Alt+P`) → "Set Server Domain" → Enter your domain

3. **Done**

   Let's Encrypt automatically provisions an SSL certificate. Your TUI will now connect via HTTPS.

**Note:** Ports 80 and 443 must be open on your server for SSL to work.

## Pod Domains

Each pod needs at least one domain to be accessible.

### Auto-Generated Domains

The quickest way to get started. Deeploy generates a subdomain using [sslip.io](https://sslip.io):

```
pod-abc123.1.2.3.4.sslip.io
```

Works instantly, no DNS configuration needed.

### Custom Domains

For production apps, use your own domain:

1. **Add DNS record**
   ```
   Type: A
   Name: myapp (for myapp.example.com)
   Value: YOUR_SERVER_IP
   ```

2. **Add domain to pod**

   Pod → Domains → New → Enter `myapp.example.com`

3. **Deploy**

   SSL certificate is automatically provisioned.

### Multiple Domains

A single pod can have multiple domains. Useful for:
- `www.example.com` and `example.com`
- Different subdomains pointing to the same app
