# deeploy MVP 0.1 - Implementation TODOs

## Overview
35 tasks organized in 9 phases to implement the complete MVP 0.1.

---

## PHASE 1: Database Schema (4 tasks)

### 1.1 Create Project Model
- [ ] Create `internal/app/models/project.go`
- [ ] Define `ProjectDB` struct with fields: ID, UserID, Name, CreatedAt, UpdatedAt
- [ ] Define `ProjectApp` struct for API responses
- [ ] Add `ToProjectApp()` converter method

### 1.2 Create Service Model
- [ ] Create `internal/app/models/service.go`
- [ ] Define `ServiceDB` struct with fields:
  - ID, ProjectID, Name, GitURL, GitToken (encrypted)
  - BuildMethod, Port, Domains (JSON), EnvVars (JSON, encrypted)
  - Status, ContainerIDs (JSON), CreatedAt, UpdatedAt
- [ ] Define `ServiceApp` struct for API responses
- [ ] Add `ToServiceApp()` converter method

### 1.3 Update User Model
- [ ] Add `GitHubToken string` field to `UserDB` (will be encrypted)
- [ ] Update database schema/migration

### 1.4 Create Migration Script
- [ ] Create SQL migration for projects table
- [ ] Create SQL migration for services table
- [ ] Update users table with github_token column
- [ ] Add to `internal/app/db/db.go` init

---

## PHASE 2: Encryption (3 tasks)

### 2.1 Create Crypto Package
- [ ] Create `internal/app/crypto/crypto.go`
- [ ] Implement AES-256-GCM encryption
- [ ] Get encryption key from `ENCRYPTION_KEY` env var
- [ ] Handle key generation/validation

### 2.2 Add Encrypt/Decrypt Functions
- [ ] `Encrypt(plaintext string) (string, error)` - returns base64
- [ ] `Decrypt(ciphertext string) (string, error)` - from base64
- [ ] Add error handling for invalid keys/data

### 2.3 Update Repositories for Encryption
- [ ] Update user repo to encrypt/decrypt GitHubToken
- [ ] Service repo will encrypt GitToken and EnvVars (in Phase 4)

---

## PHASE 3: Projects API (4 tasks)

### 3.1 Create Project Repository
- [ ] Create `internal/app/repos/project.go`
- [ ] Implement `Create(project *models.ProjectDB) error`
- [ ] Implement `GetByID(id string) (*models.ProjectDB, error)`
- [ ] Implement `GetByUserID(userID string) ([]*models.ProjectDB, error)`
- [ ] Implement `Update(project *models.ProjectDB) error`
- [ ] Implement `Delete(id string) error`

### 3.2 Create Project Service Layer
- [ ] Create `internal/app/services/project.go`
- [ ] Add business logic for project creation
- [ ] Add validation (name length, user ownership, etc.)
- [ ] Handle project deletion (cascade delete services)

### 3.3 Create Project Handlers
- [ ] Create `internal/app/handler/project.go`
- [ ] `POST /projects` - Create project
- [ ] `GET /projects` - List user's projects
- [ ] `GET /projects/:id` - Get project details
- [ ] `DELETE /projects/:id` - Delete project + services
- [ ] Add JWT auth middleware

### 3.4 Wire Up Routes
- [ ] Create `internal/app/routes/project.go`
- [ ] Register all project routes in main.go
- [ ] Test with curl/Postman

---

## PHASE 4: Services API (4 tasks)

### 4.1 Create Service Repository
- [ ] Create `internal/app/repos/service.go`
- [ ] Implement `Create(service *models.ServiceDB) error`
- [ ] Implement `GetByID(id string) (*models.ServiceDB, error)`
- [ ] Implement `GetByProjectID(projectID string) ([]*models.ServiceDB, error)`
- [ ] Implement `Update(service *models.ServiceDB) error`
- [ ] Implement `Delete(id string) error`
- [ ] Encrypt/decrypt GitToken and EnvVars fields

### 4.2 Create Service Service Layer
- [ ] Create `internal/app/services/service.go`
- [ ] Add validation (GitHub URL, domains, etc.)
- [ ] Parse and validate domain names
- [ ] Handle build method selection

### 4.3 Create Service Handlers (Basic CRUD)
- [ ] Create `internal/app/handler/service.go`
- [ ] `POST /projects/:pid/services` - Add service
- [ ] `GET /projects/:pid/services` - List services
- [ ] `GET /services/:id` - Get service details
- [ ] `PUT /services/:id` - Update service config
- [ ] `DELETE /services/:id` - Delete service
- [ ] Add JWT auth middleware

### 4.4 Wire Up Routes
- [ ] Create `internal/app/routes/service.go`
- [ ] Register all service routes in main.go
- [ ] Test basic CRUD with curl

---

## PHASE 5: Docker Integration (3 tasks)

### 5.1 Create Docker Client Wrapper
- [ ] Create `internal/app/docker/client.go`
- [ ] Initialize Docker client with socket
- [ ] Add connection validation
- [ ] Handle Docker API errors

### 5.2 Implement Container Operations
- [ ] Create `internal/app/docker/container.go`
- [ ] `StartContainer(serviceID string) error`
- [ ] `StopContainer(containerID string) error`
- [ ] `DeleteContainer(containerID string) error`
- [ ] `GetContainerStatus(containerID string) (string, error)`
- [ ] `ListContainersByLabel(label string) ([]string, error)`

### 5.3 Implement Logs Streaming
- [ ] Create `internal/app/docker/logs.go`
- [ ] `StreamLogs(containerID string, follow bool) (io.ReadCloser, error)`
- [ ] Create SSE endpoint handler for logs
- [ ] `GET /services/:id/logs?follow=true` endpoint
- [ ] Handle connection cleanup

---

## PHASE 6: Deployment Engine (6 tasks)

### 6.1 Create Git Operations Module
- [ ] Create `internal/app/git/git.go`
- [ ] `Clone(url, token, destPath string) error` with exec
- [ ] `Pull(repoPath string) error` with exec
- [ ] Handle authentication with token in URL
- [ ] Add error handling for git failures

### 6.2 Create Compose Generator (Dockerfile Mode)
- [ ] Create `internal/app/deployer/compose.go`
- [ ] `GenerateComposeForDockerfile(service *models.ServiceDB) (string, error)`
- [ ] Template for single service from Dockerfile
- [ ] Add Traefik labels for domains
- [ ] Add networks (project network + traefik)
- [ ] Set port mapping

### 6.3 Create Compose Generator (Compose Mode)
- [ ] `GenerateComposeWrapper(service *models.ServiceDB) (string, error)`
- [ ] Use `include:` to reference user's compose
- [ ] Inject Traefik labels to services
- [ ] Add networks to all services
- [ ] Handle multi-service compose files

### 6.4 Create .env File Writer
- [ ] Create `internal/app/deployer/env.go`
- [ ] `WriteEnvFile(path string, envVars map[string]string) error`
- [ ] Format as KEY=VALUE lines
- [ ] Handle special characters/escaping

### 6.5 Implement Deploy Handler
- [ ] Create `internal/app/deployer/deployer.go`
- [ ] Main `Deploy(service *models.ServiceDB)` function
- [ ] Orchestrate: git clone/pull → env write → compose generate → docker up
- [ ] Update service status to "deploying" → "running" / "failed"
- [ ] Capture and store container IDs
- [ ] Handle errors and rollback

### 6.6 Add Deploy Endpoints
- [ ] `POST /services/:id/deploy` - Deploy/redeploy service
- [ ] `POST /services/:id/start` - Start stopped service
- [ ] `POST /services/:id/stop` - Stop running service
- [ ] `GET /services/:id/status` - Get current status
- [ ] Wire up in routes

---

## PHASE 7: TUI (6 tasks)

### 7.1 Update Dashboard for Projects List
- [ ] Update `internal/cli/ui/pages/dashboard.go`
- [ ] Fetch projects from API `GET /projects`
- [ ] Display project cards with service count
- [ ] Add "Open" button to view project details
- [ ] Add "+ New Project" button

### 7.2 Create Project Detail Page
- [ ] Create `internal/cli/ui/pages/project_detail.go`
- [ ] Fetch project + services from API
- [ ] Display services list with status indicators
- [ ] Show domains for each service
- [ ] Add action buttons: Logs, Stop/Start, Redeploy, Delete
- [ ] Add "+ Add Service" button

### 7.3 Create New Project Form
- [ ] Create `internal/cli/ui/pages/project_form.go`
- [ ] Input: Project name
- [ ] Validation
- [ ] Submit → POST /projects
- [ ] Navigate to project detail on success

### 7.4 Create Add Service Form
- [ ] Create `internal/cli/ui/pages/service_form.go`
- [ ] Inputs: Name, GitHub URL, Build Method (radio)
- [ ] Port input (if Dockerfile)
- [ ] Domains input (comma-separated)
- [ ] ENV vars (key-value pairs, add/remove)
- [ ] Validation
- [ ] Submit → POST /projects/:pid/services
- [ ] Show deployment progress

### 7.5 Create Logs Viewer
- [ ] Create `internal/cli/ui/pages/logs.go`
- [ ] Connect to SSE endpoint `GET /services/:id/logs?follow=true`
- [ ] Display logs with auto-scroll
- [ ] Toggle auto-scroll on/off
- [ ] Handle connection errors
- [ ] Back button

### 7.6 Wire Up API Calls
- [ ] Create `internal/cli/api/client.go` HTTP client
- [ ] Add JWT token to all requests
- [ ] Handle auth errors (redirect to connect page)
- [ ] Add loading states in TUI
- [ ] Add error messages in TUI

---

## PHASE 8: Infrastructure (3 tasks)

### 8.1 Update install.sh
- [ ] Update `internal/app/install/install.sh`
- [ ] Add Traefik container setup with Let's Encrypt
- [ ] Configure Traefik networks
- [ ] Update deeployd container with proper volumes:
  - `/var/run/docker.sock:/var/run/docker.sock`
  - `/var/deeploy/projects:/var/deeploy/projects`
  - `/var/deeploy/data:/app/data`
- [ ] Generate and pass ENCRYPTION_KEY env var
- [ ] Test on fresh Ubuntu VPS

### 8.2 Create Dockerfile.deeployd
- [ ] Create `Dockerfile.deeployd` in project root
- [ ] Multi-stage build (builder + runtime)
- [ ] Install git in final image
- [ ] Install docker-cli and docker-compose in final image
- [ ] Copy built binary
- [ ] Expose port 8090
- [ ] Set CMD

### 8.3 Update Makefile
- [ ] Add `build-deeployd` target
- [ ] Build Docker image: `docker build -f Dockerfile.deeployd -t deeployd:latest .`
- [ ] Add `push-deeployd` target for ghcr.io
- [ ] Add `run-deeployd-local` for local testing

---

## PHASE 9: Testing (3 tasks)

### 9.1 Test Local Dev Setup
- [ ] Start Traefik container locally
- [ ] Run deeployd with `go run cmd/app/main.go`
- [ ] Run deeploy TUI with `go run cmd/cli/main.go`
- [ ] Connect TUI to localhost:8090
- [ ] Create test project
- [ ] Verify database entries

### 9.2 Test Domain Routing
- [ ] Add test domains to `/etc/hosts`:
  - `127.0.0.1 test.local`
  - `127.0.0.1 api.test.local`
- [ ] Deploy a test service with domain `test.local`
- [ ] Verify Traefik routes traffic correctly
- [ ] Test with `curl http://test.local`

### 9.3 End-to-End Deployment Test
- [ ] Create test GitHub repo with Dockerfile
- [ ] Add service via TUI with test repo
- [ ] Deploy and verify:
  - Git clone works
  - Docker image builds
  - Container starts
  - Traefik routes work
  - Logs streaming works
- [ ] Test stop/start/redeploy
- [ ] Test delete

---

## Progress Tracking

- **PHASE 1:** 0/4 tasks completed
- **PHASE 2:** 0/3 tasks completed
- **PHASE 3:** 0/4 tasks completed
- **PHASE 4:** 0/4 tasks completed
- **PHASE 5:** 0/3 tasks completed
- **PHASE 6:** 0/6 tasks completed
- **PHASE 7:** 0/6 tasks completed
- **PHASE 8:** 0/3 tasks completed
- **PHASE 9:** 0/3 tasks completed

**Total:** 0/35 tasks completed

---

## Notes

- All phases can be worked on somewhat independently after Phase 1-2
- Phase 6 (Deployment Engine) is the most complex
- Phase 7 (TUI) depends on Phase 3-4 APIs being ready
- Phase 9 (Testing) validates everything works together

## Next Steps

Start with Phase 1 to establish the data model foundation, then Phase 2 for security, then build out the API layers (Phase 3-4) before tackling the deployment engine (Phase 5-6).
