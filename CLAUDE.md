# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go-based license management service supporting various software products (JetBrains, GitLab, FinalShell, MobaXterm, JRebel). Built with Gin framework, GORM for database access, and a React frontend embedded into the Go binary.

## Build & Run Commands

### Backend (Go)
```bash
# Install dependencies
go mod download

# Build
go build -o license-server

# Run
./license-server

# Run with version flag
./license-server -version
```

### Frontend (React)
```bash
cd web
npm install
npm run build    # Build for production (output to web/build/)
npm start        # Development server
npm test         # Run tests
```

### Docker
```bash
# Build image
docker build -t license-server .

# Run with docker-compose
docker-compose up -d
```

### Testing
```bash
# Run all Go tests
go test ./...

# Run benchmarks
go test -bench=. ./test/...

# Run specific test package
go test ./test/jetbrains/
go test ./test/jrebel/
```

## Architecture

### Entry Point & Routing
- `main.go` - Application entry, embeds frontend via `//go:embed web/build`, configures Gin router
- `router/router.go` - API route definitions, supports both `/api/*` prefix and root-level paths

### Service Modules (Controller-Service pattern)
Each license type has its own package:
- `jetbrains/` - JetBrains products (has v1/v2 API versions under `code/api/` and `code/service/`)
- `gitlab/` - GitLab license
- `finalshell/` - FinalShell license
- `mobaxterm/` - MobaXterm license
- `jrebel/` - JRebel license

### Core Infrastructure
- `config/` - Configuration loading from env vars, database setup
- `initialize/` - Component initialization (certificates, JetBrains, GitLab)
- `cron/` - Scheduled tasks
- `logger/` - Logging utilities
- `crypto/` - Encryption utilities (AES-CBC-PKCS7)

### Frontend
- `web/` - React 19 + TypeScript + Ant Design application
- Embedded into Go binary at build time

## Configuration

Environment variables (or `.env` file):
- `HTTP_HOST` - Server host (default: 0.0.0.0)
- `HTTP_PORT` - Server port (default: 5000)
- `DATA_DIR` - Data directory (default: /data)
- `DATABASE_DRIVER` - sqlite, mysql, or postgresql (default: sqlite)
- `DATABASE_DSN` - Connection string (for sqlite: filename, for mysql: user:pass@tcp(host:port)/dbname)

## API Patterns

APIs are registered in `router/router.go` under these route groups:
- `/server/` - Server status/version
- `/jetbrains/` - License generation, product/plugin management
- `/gitlab/` - License generation
- `/final-shell/` - License generation
- `/mobaxterm/` - License generation
- `/jrebel/` and `/agent/` - JRebel leases

All routes work both with and without `/api/` prefix.