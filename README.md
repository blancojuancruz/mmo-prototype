# Server Setup Guide

## Prerequisites

| Tool | Purpose |
|------|---------|
| [Docker Desktop](https://www.docker.com/products/docker-desktop/) | Runs PostgreSQL and Redis containers |
| [Go 1.22+](https://go.dev/dl/) | Runs the game server |
| [golangci-lint](https://golangci-lint.run/usage/install/) | Code linting (optional) |

---

## Quick Start

### 1. Configure environment
```bash
cp .env.example .env
cp docker-compose.example.yml docker-compose.yml
```

Edit both files and set your credentials. Make sure `DB_USER` and `DB_PASSWORD` match between `.env` and `docker-compose.yml`.

### 2. Start the database

Make sure Docker Desktop is running, then:
```bash
docker-compose up -d
```

Verify containers are healthy:
```bash
docker ps
```

You should see `mmorpg-postgres` and `mmorpg-redis` with status `Up`.

### 3. Start the server
```bash
cd server
go run cmd/server/main.go
```

Expected output:
```
🎮 MMORPG Server starting...
✅ PostgreSQL connected
✅ Migrations completed
Server listening on port :8080
```

---

## Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/ping` | Health check |
| POST | `/auth/register` | Register a new account |
| POST | `/auth/login` | Login and get character data |
| POST | `/auth/character` | Create a character |
| POST | `/game/save_position` | Save character position |
| WS | `/ws` | WebSocket connection for real-time gameplay |

---

## Troubleshooting

**Docker daemon not running**
Open Docker Desktop and wait for the whale icon in the taskbar before running any `docker` commands.

**Port conflict on 5432**
Another PostgreSQL instance may be running locally. Stop it or change the port in `docker-compose.yml` and `.env`.

**Server fails to connect to DB**
Always start the Docker containers before running the server. The database must be ready before the server boots.

**Missing dependencies**
```bash
cd server
go mod tidy
```

**Linting**
```bash
cd server
golangci-lint run
```
