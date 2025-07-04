# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Deeploy is a Go-based deployment/project management tool with both web and terminal user interfaces (TUI). It uses SQLite for data storage and JWT for authentication.

## Development Commands

### Web Development (Full Stack)
```bash
make dev              # Runs templ, server, and tailwind in parallel
make server           # Run web server only with hot reload (port 8090)
make templ            # Watch and generate templ templates
make tailwind         # Watch and compile Tailwind CSS
```

### TUI Development
```bash
go run cmd/tui/main.go        # Run TUI application
DEBUG=1 go run cmd/tui/main.go # Run with debug logging
make cli-debug                # Debug with Delve (port 43000)
make cli-log                  # Tail debug.log file
```

### Building
```bash
./scripts/build-binaries.sh   # Build for darwin/linux amd64/arm64
```

## Architecture

### Directory Structure
- `/cmd/web/` - Web server entry point
- `/cmd/tui/` - Terminal UI entry point
- `/internal/` - Core business logic:
  - `handlers/api/` - JSON API handlers
  - `handlers/web/` - HTML template handlers
  - `services/` - Business logic layer
  - `data/` - Repository pattern for data access
  - `ui/` - Templ templates for web UI
  - `tui/` - Bubble Tea components for terminal UI

### Tech Stack
- **Backend**: Go 1.23.3
- **Database**: SQLite with golang-migrate
- **Web UI**: Templ templates + Tailwind CSS
- **TUI**: Bubble Tea + Lipgloss
- **Auth**: JWT tokens (Bearer for API, Cookie for web)

### Key Patterns

#### Authentication Flow
- Dual token support: Bearer tokens for CLI/API, cookies for web
- Three middleware types:
  - `Auth()` - Validates token, adds user to context
  - `RequireAuth()` - Redirects to login if not authenticated
  - `RequireGuest()` - Redirects authenticated users away

#### Route Structure
- **API Routes**: `/api/*` prefix, JSON responses
- **Web Routes**: Direct paths, HTML responses
- All routes defined in `/internal/routes/`

#### Data Model
- **Users**: Authentication and ownership
- **Projects**: Top-level organizational unit
- **Pods**: Services within projects (renamed from "services")

### Environment Variables
Create `.env` file with:
```
GO_ENV=dev
JWT_SECRET=your-secret-key
```

## Common Tasks

### Adding a New API Endpoint
1. Create handler in `/internal/handlers/api/`
2. Add route in `/internal/routes/api.go`
3. Use JSON request/response pattern
4. Apply `auth.Auth()` middleware if authentication required

### Adding a New Web Page
1. Create templ template in `/internal/ui/`
2. Create handler in `/internal/handlers/web/`
3. Add route in `/internal/routes/web.go`
4. Run `make dev` to see changes with hot reload

### Working with Templates
- Templ files use `.templ` extension
- Components can be imported and composed
- Run `make templ` to watch for changes

### Database Changes
1. Create migration file in `/internal/db/migrations/`
2. Follow naming convention: `YYYYMMDDHHMMSS_description.up.sql`
3. Migrations run automatically on startup

## Important Notes

- Always use the repository pattern in `/internal/data/` for database access
- Keep business logic in `/internal/services/`
- Use context for passing user information from middleware
- TUI and web share the same authentication system
- No test files exist yet - consider adding tests when implementing new features