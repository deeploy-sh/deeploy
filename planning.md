# deeploy MVP 0.1 - Complete Implementation Plan

## Vision
Docker deployment platform with TUI-first approach. Alternative to Dokploy/Coolify.

## Architecture

### Components
1. **deeploy** (CLI/TUI Client)
   - Runs locally on developer machine (Mac/Linux/Windows)
   - Bubbletea TUI framework
   - Connects to deeployd server via HTTP/REST API
   - Displays projects, services, logs, status

2. **deeployd** (Server/Daemon)
   - Runs on VPS as Docker container
   - Go HTTP server (Port 8090)
   - SQLite database
   - Docker SDK for container management
   - Git installed in container for repo cloning
   - Mounts: /var/run/docker.sock (controls host Docker)

3. **Traefik** (Reverse Proxy)
   - Runs as Docker container on VPS
   - Exposes ports 80/443
   - Reads Docker labels from containers
   - Automatic SSL via Let's Encrypt
   - Routes traffic: Domain → Container

---

## Core Concepts

### Project vs Service Hierarchy

**Project** = Logical grouping (UI/DB wrapper)
- Has a name
- Groups multiple services
- All services share a Docker network

**Service** = 1 GitHub repo = 1 deployment unit
- Belongs to a project
- Has own Git repository
- Either Dockerfile OR docker-compose.yml
- Own domains, ports, ENV vars
- Can have multiple containers (with compose)

**Example:**
```
Project: "my-saas-app"
├── Service: "frontend"
│   ├── GitHub: github.com/user/frontend
│   ├── Method: docker-compose.yml
│   ├── Domains: [myapp.com, www.myapp.com]
│   ├── ENV: {API_URL=https://api.myapp.com}
│   └── Status: Running ✓
│
├── Service: "api"
│   ├── GitHub: github.com/user/backend-api
│   ├── Method: Dockerfile
│   ├── Port: 4000
│   ├── Domains: [api.myapp.com]
│   ├── ENV: {DB_HOST=db, DB_PASSWORD=secret}
│   └── Status: Running ✓
│
└── Service: "admin"
    ├── GitHub: github.com/user/admin-dashboard
    ├── Method: docker-compose.yml
    ├── Domains: [admin.myapp.com]
    └── Status: Stopped
```

**Networking:**
- All services in project "my-saas-app" in network "my-saas-app_network"
- Services can reach each other: `http://api:4000`
- Additionally all in "traefik" network for external access

---

## Database Schema

```go
type User struct {
    ID            string
    Email         string
    Password      string    // bcrypt hashed
    GitHubToken   string    // AES-256 encrypted, Personal Access Token
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type Project struct {
    ID        string
    UserID    string
    Name      string      // e.g. "my-saas-app"
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Service struct {
    ID           string
    ProjectID    string
    Name         string           // e.g. "frontend", "api"
    GitURL       string           // github.com/user/repo
    GitToken     string           // encrypted (or uses User.GitHubToken)
    BuildMethod  string           // "dockerfile" | "compose"
    Port         int              // only for Dockerfile, e.g. 3000
    Domains      []string         // ["myapp.com", "www.myapp.com"]
    EnvVars      map[string]string // encrypted, e.g. {API_KEY: "secret"}
    Status       string           // "deploying" | "running" | "stopped" | "failed"
    ContainerIDs []string         // can be multiple with compose
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

---

## API Endpoints (deeployd Server)

### Authentication
```
POST /auth/login
POST /auth/register
GET  /dashboard (JWT protected)
```

### Projects
```
POST   /projects             - Create new project
GET    /projects             - List all user projects
GET    /projects/:id         - Get project details
DELETE /projects/:id         - Delete project + all services
```

### Services
```
POST   /projects/:pid/services              - Add service to project
GET    /projects/:pid/services              - List services in project
GET    /services/:id                        - Get service details
PUT    /services/:id                        - Update service config
DELETE /services/:id                        - Delete service

POST   /services/:id/deploy                 - Deploy/Redeploy service
POST   /services/:id/start                  - Start stopped service
POST   /services/:id/stop                   - Stop running service
GET    /services/:id/logs?follow=true       - Stream logs (SSE/WebSocket)
GET    /services/:id/status                 - Get current status
```

---

## Deployment Flow (what happens on deeployd?)

### Trigger: User clicks "Deploy" in TUI

**1. API Request**
```json
POST /services/:id/deploy
Authorization: Bearer <jwt-token>
```

**2. Server Logic (deeployd)**
```go
// Status Update
service.Status = "deploying"
db.Save(service)

// Git Operations
projectDir := "/var/deeploy/projects/" + service.ProjectID
serviceDir := projectDir + "/" + service.ID

if exists(serviceDir) {
    // Redeploy: git pull
    exec("git", "-C", serviceDir, "pull")
} else {
    // Initial deploy: git clone
    gitURL := fmt.Sprintf("https://%s@%s", token, service.GitURL)
    exec("git", "clone", gitURL, serviceDir)
}

// Detect Build Method
hasCompose := fileExists(serviceDir + "/docker-compose.yml")
hasDockerfile := fileExists(serviceDir + "/Dockerfile")

// Generate .env file
writeEnvFile(serviceDir + "/.env", service.EnvVars)

// Generate deeploy compose wrapper
if service.BuildMethod == "compose" {
    generateComposeWrapper(service)
} else {
    generateComposeFromDockerfile(service)
}

// Docker Operations
exec("docker-compose", "-f", "docker-compose.deeploy.yml",
     "-p", service.ProjectID, "up", "-d", "--build")

// Update Status
service.Status = "running"
service.ContainerIDs = getContainerIDs()
db.Save(service)
```

**3. Generated docker-compose.deeploy.yml (Dockerfile Example)**
```yaml
version: '3.8'

services:
  api:
    build: .
    container_name: my-saas-app_api
    env_file: .env
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api-${SERVICE_ID}.rule=Host(`api.myapp.com`)"
      - "traefik.http.services.api-${SERVICE_ID}.loadbalancer.server.port=4000"
      - "traefik.http.routers.api-${SERVICE_ID}.tls.certresolver=letsencrypt"
    networks:
      - my-saas-app_network
      - traefik

networks:
  my-saas-app_network:
    name: my-saas-app_network
  traefik:
    external: true
```

**4. Generated docker-compose.deeploy.yml (Compose Example)**
```yaml
version: '3.8'

# Include user's original compose
include:
  - docker-compose.yml

# Inject Traefik labels + networks
services:
  frontend:  # Service from user compose
    env_file: .env
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.frontend-${SERVICE_ID}.rule=Host(`myapp.com`, `www.myapp.com`)"
      - "traefik.http.routers.frontend-${SERVICE_ID}.tls.certresolver=letsencrypt"
    networks:
      - my-saas-app_network
      - traefik

  backend:   # If in user compose
    env_file: .env
    networks:
      - my-saas-app_network

networks:
  my-saas-app_network:
    name: my-saas-app_network
  traefik:
    external: true
```

---

## TUI Pages & Flows

### 1. Connect Page (existing)
- Enter server IP:Port
- Browser OAuth flow
- Save token

### 2. Dashboard/Projects List (new)
```
┌─────────────────────────────────────────────────┐
│ 🔥 deeploy                                      │
├─────────────────────────────────────────────────┤
│                                                 │
│  Projects                                       │
│  ┌───────────────────────────────────────────┐ │
│  │ 📁 my-saas-app                      [Open]│ │
│  │    3 services • 2 running • 1 stopped    │ │
│  ├───────────────────────────────────────────┤ │
│  │ 📁 test-project                     [Open]│ │
│  │    1 service • 1 running                 │ │
│  └───────────────────────────────────────────┘ │
│                                                 │
│  [+ New Project]                                │
│                                                 │
│  [Esc] Exit                                     │
└─────────────────────────────────────────────────┘
```

### 3. Project Detail Page (new)
```
┌─────────────────────────────────────────────────┐
│ 📁 my-saas-app                          [← Back]│
├─────────────────────────────────────────────────┤
│                                                 │
│  Services                                       │
│  ┌───────────────────────────────────────────┐ │
│  │ 🟢 frontend                               │ │
│  │    myapp.com, www.myapp.com               │ │
│  │    [Logs] [Stop] [Redeploy] [Delete]      │ │
│  ├───────────────────────────────────────────┤ │
│  │ 🟢 api                                    │ │
│  │    api.myapp.com                          │ │
│  │    [Logs] [Stop] [Redeploy] [Delete]      │ │
│  ├───────────────────────────────────────────┤ │
│  │ 🔴 admin                                  │ │
│  │    admin.myapp.com                        │ │
│  │    [Logs] [Start] [Redeploy] [Delete]     │ │
│  └───────────────────────────────────────────┘ │
│                                                 │
│  [+ Add Service]                                │
│                                                 │
└─────────────────────────────────────────────────┘
```

### 4. Add Service Form (new)
```
┌─────────────────────────────────────────────────┐
│ Add Service to: my-saas-app             [← Back]│
├─────────────────────────────────────────────────┤
│                                                 │
│  Service Name                                   │
│  ┌───────────────────────────────────────────┐ │
│  │ worker-queue                              │ │
│  └───────────────────────────────────────────┘ │
│                                                 │
│  GitHub URL                                     │
│  ┌───────────────────────────────────────────┐ │
│  │ github.com/user/worker                    │ │
│  └───────────────────────────────────────────┘ │
│                                                 │
│  Build Method                                   │
│  ○ Dockerfile    ● docker-compose.yml          │
│                                                 │
│  Domains (optional, comma separated)            │
│  ┌───────────────────────────────────────────┐ │
│  │ worker.myapp.com                          │ │
│  └───────────────────────────────────────────┘ │
│                                                 │
│  Environment Variables                          │
│  ┌─────────────────┬─────────────────────────┐ │
│  │ REDIS_URL       │ redis://localhost:6379  │ │
│  │ QUEUE_NAME      │ default                 │ │
│  │ [+ Add Variable]                          │ │
│  └───────────────────────────────────────────┘ │
│                                                 │
│  [Deploy Service]  [Cancel]                     │
│                                                 │
└─────────────────────────────────────────────────┘
```

### 5. Logs View (new)
```
┌─────────────────────────────────────────────────┐
│ Logs: my-saas-app / api                 [← Back]│
├─────────────────────────────────────────────────┤
│                                                 │
│  [2024-01-03 10:23:15] Server starting...       │
│  [2024-01-03 10:23:16] Connected to database    │
│  [2024-01-03 10:23:16] Listening on :4000       │
│  [2024-01-03 10:24:32] GET /api/users 200 45ms  │
│  [2024-01-03 10:24:33] POST /api/auth 201 120ms │
│  ...                                            │
│                                                 │
│  [↓ Auto-scroll: ON]  [Esc] Back                │
│                                                 │
└─────────────────────────────────────────────────┘
```

---

## Installation & Setup

### VPS Installation (install.sh)

```bash
#!/usr/bin/env bash
set -euo pipefail

# 1. System Check
if [[ $(uname) != "Linux" ]]; then
    echo "Error: Linux required"
    exit 1
fi

if [[ $EUID -ne 0 ]]; then
    echo "Error: Run with sudo"
    exit 1
fi

# 2. Install Docker
if ! command -v docker &>/dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com | bash
fi

# 3. Create Docker Networks
docker network create traefik 2>/dev/null || true

# 4. Start Traefik
echo "Starting Traefik..."
docker run -d \
  --name traefik \
  --restart unless-stopped \
  --network traefik \
  -p 80:80 -p 443:443 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /var/deeploy/traefik/letsencrypt:/letsencrypt \
  traefik:v2.10 \
  --api.dashboard=true \
  --providers.docker=true \
  --providers.docker.network=traefik \
  --entrypoints.web.address=:80 \
  --entrypoints.websecure.address=:443 \
  --entrypoints.web.http.redirections.entrypoint.to=websecure \
  --certificatesresolvers.letsencrypt.acme.email=admin@example.com \
  --certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json \
  --certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web

# 5. Start deeployd
echo "Starting deeployd..."
docker pull ghcr.io/deeploy-sh/deeployd:latest
docker rm -f deeployd 2>/dev/null || true
docker run -d \
  --name deeployd \
  --restart unless-stopped \
  --network traefik \
  -p 8090:8090 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/deeploy/projects:/var/deeploy/projects \
  -v /var/deeploy/data:/app/data \
  -e ENCRYPTION_KEY=$(openssl rand -hex 32) \
  ghcr.io/deeploy-sh/deeployd:latest

IP=$(hostname -I | awk '{print $1}')
echo "✨ deeploy installed! Access at: http://$IP:8090"
```

### deeployd Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

# Install git (needed for cloning user repos)
RUN apk add --no-cache git

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o deeployd ./cmd/app

FROM alpine:latest

# Install git and docker CLI
RUN apk add --no-cache git docker-cli docker-compose

WORKDIR /app
COPY --from=builder /build/deeployd .

EXPOSE 8090
CMD ["./deeployd"]
```

---

## Local Development Setup

### Option 1: Docker Desktop (Mac)
```bash
# Terminal 1: Traefik
docker network create traefik
docker run -d --name traefik --network traefik \
  -p 80:80 -p 443:443 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  traefik:v2.10 \
  --api.dashboard=true \
  --providers.docker=true \
  --entrypoints.web.address=:80

# Terminal 2: deeployd (native for fast dev)
cd cmd/app
export ENCRYPTION_KEY=$(openssl rand -hex 32)
go run main.go

# Terminal 3: deeploy TUI
cd cmd/cli
go run main.go
```

### Option 2: Colima (Docker Alternative)
```bash
# Start Colima (does the same as Docker Desktop)
colima start --cpu 4 --memory 8

# Then same as Option 1
```

### Option 3: Local Domain Testing (Mac)

For testing with fake domains locally:

```bash
# 1. Edit /etc/hosts
sudo nano /etc/hosts

# Add test domains:
127.0.0.1 myapp.local
127.0.0.1 api.myapp.local
127.0.0.1 admin.myapp.local

# 2. Start Traefik (listens on port 80)
docker network create traefik
docker run -d --name traefik --network traefik \
  -p 80:80 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  traefik:v2.10 \
  --api.dashboard=true \
  --providers.docker=true \
  --entrypoints.web.address=:80

# 3. Deploy test service with domain
# In TUI: Add service with domain "myapp.local"
# Traefik routes: myapp.local -> container

# 4. Test
curl http://myapp.local
# Works! But no SSL (http only)
```

**Limitations:**
- No SSL/HTTPS (Let's Encrypt needs real domain)
- Must use `.local` or similar (not real TLDs)
- Need sudo to edit /etc/hosts

**Advantages:**
- Fast iteration
- No external dependencies
- Free

### Testing with Real Domains
```bash
# Hetzner VPS (CX11, ~3€/month)
# DNS Wildcard: *.dev.yourdomain.com → VPS IP
# Then real SSL testing possible with Let's Encrypt
```

---

## Security

### Token Encryption
```go
// AES-256-GCM for token/ENV storage
func Encrypt(plaintext string, key []byte) (string, error)
func Decrypt(ciphertext string, key []byte) (string, error)

// Key from environment: ENCRYPTION_KEY
// Generated during installation: openssl rand -hex 32
```

### Docker Socket Security
- deeployd runs as container
- Has access to host Docker socket
- → Can do EVERYTHING on host (privileged!)
- Important: deeployd server must be well secured (JWT auth, etc.)

---

## Out of Scope for 0.1

### Later (0.2+)
- GitHub OAuth instead of Personal Access Token
- Webhooks for auto-deploy on git push
- Branch selection (only main/master for 0.1)
- Rollback to older commits
- CPU/Memory monitoring per container
- Build cache optimization
- Docker registry support (private images)
- Database backups
- Service templates (e.g. "Next.js App")

---

## Research Topics Before Implementation

- [ ] Docker socket security best practices
- [ ] Traefik Let's Encrypt rate limits
- [ ] Multi-service docker-compose handling (which services get Traefik labels?)
- [ ] Token encryption: AES-256-GCM vs ChaCha20-Poly1305
- [ ] WebSocket vs SSE for log streaming
- [ ] Git credential handling in deeployd container
- [ ] How to detect which service in compose should be exposed (by port? by service name convention?)
