---
name: devops
description: "Use when: writing a Dockerfile for the Go+Vite build, setting up CI/CD with GitHub Actions, deploying to Fly.io, configuring environment variables for production, setting up health checks, or planning the migration from in-memory to persistent room state."
type: skill
---

# DevOps Skill — Hexar Deployment & Operations

Handles the full deployment lifecycle for Hexar: containerisation, CI/CD, Fly.io deployment, environment configuration, and the path from in-memory MVP to persistent production state.

**Prerequisites:** None. Can be invoked standalone.

---

## Context: Hexar's Deployment Shape

- **Single Go binary** serves both the REST API (`/lobby/*`) and static client files (Vite `dist/`)
- **WebSockets** require persistent connections — rules out serverless, requires sticky routing
- **In-memory room state** — server restart = all active games lost (known MVP limitation)
- **No database yet** — rooms and sessions live in Go maps
- **Current stack:** Go server on `:8080`, Vite client served from `dist/` directory

---

## Workflow

### 1. Dockerfile (Multi-Stage)

Standard Go + static assets pattern:

```dockerfile
# Stage 1: Build client
FROM node:20-alpine AS client-builder
WORKDIR /app/client
COPY client/package*.json ./
RUN npm ci
COPY client/ ./
RUN npm run build

# Stage 2: Build server
FROM golang:1.23-alpine AS server-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 3: Runtime
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=server-builder /app/server ./server
COPY --from=client-builder /app/client/dist ./client/dist
EXPOSE 8080
CMD ["./server", "-client-dir", "./client/dist"]
```

Key choices:
- **distroless runtime** — minimal attack surface, no shell, ~3MB image
- **CGO_ENABLED=0** — fully static binary, no glibc dependency
- **`npm ci`** not `npm install` — reproducible installs in CI

### 2. Fly.io Configuration

`fly.toml` essentials for Hexar:

```toml
app = "hexar"
primary_region = "ams"  # pick closest to your users

[build]
  dockerfile = "Dockerfile"

[env]
  PORT = "8080"

[http_service]
  internal_port = 8080
  force_https = true          # auto TLS, ws:// → wss://
  auto_stop_machines = false  # keep alive for WebSocket sessions
  auto_start_machines = true
  min_machines_running = 1

[[vm]]
  memory = "256mb"
  cpu_kind = "shared"
  cpus = 1
```

Critical: `auto_stop_machines = false` — Fly's default sleep kills active WebSocket connections.

### 3. Environment Variables

| Variable | Dev value | Prod value | Notes |
|---|---|---|---|
| `PORT` | `8080` | `8080` | Fly sets this automatically |
| `CLIENT_DIR` | `client/dist` | `./client/dist` | Path to Vite build output |
| `ENV` | `development` | `production` | Gate logging verbosity |

Set via `fly secrets set KEY=VALUE` for sensitive values. Non-sensitive config goes in `fly.toml [env]`.

Client-side: Vite env variable for WebSocket URL:
```typescript
// client/src/net/connection.ts
const protocol = import.meta.env.PROD ? 'wss' : 'ws';
const url = `${protocol}://${window.location.host}/ws?code=${code}&token=${token}`;
```

### 4. Health Check Endpoint

Add `GET /health` returning 200 + `{"status":"ok"}` to `internal/net/server.go`. Fly uses this to detect crashed instances.

```go
s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status":"ok"}`))
})
```

### 5. GitHub Actions CI/CD

`.github/workflows/deploy.yml`:

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: go test ./...

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

Gate deploy behind tests — never ship a broken build.

### 6. In-Memory → Persistent State (Future)

Current limitation: rooms live in a Go map. Server restart = all games lost.

Migration path when needed:
1. **SQLite via `modernc.org/sqlite`** — pure Go, no CGO, embeds in the binary; store room code, tokens, game state as JSON blob; load on startup, persist on state change
2. **Redis** — better for multi-instance scaling, but adds infrastructure dependency
3. **Graceful drain** — on SIGTERM, stop accepting new rooms, wait for active games to finish (max `disconnectGrace` = 2 min), then exit

For initial deployment, document the limitation in the lobby UI: "Server restarts end all active games."

---

## Pre-Deployment Checklist

### Build
- [ ] `docker build` succeeds locally
- [ ] Client WebSocket URL uses `wss://` in production builds
- [ ] `go build ./...` passes with no warnings
- [ ] `go test ./...` passes

### Fly.io
- [ ] `fly.toml` present with `auto_stop_machines = false`
- [ ] `FLY_API_TOKEN` set in GitHub secrets
- [ ] Health check endpoint exists and returns 200
- [ ] At least 1 machine always running (`min_machines_running = 1`)

### Post-Deploy Verification
- [ ] `https://your-app.fly.dev` loads the lobby
- [ ] Create game → room code appears
- [ ] Join from second browser → game starts
- [ ] WebSocket connects over `wss://` (check browser devtools Network tab)
- [ ] No console errors in production build

---

## Rules

- **Never deploy without tests passing.** Gate CI/CD on `go test ./...` at minimum.
- **`auto_stop_machines = false` is mandatory.** Fly's sleep kills WebSocket connections mid-game.
- **Secrets go in `fly secrets`, not `fly.toml`.** `fly.toml` is committed to git.
- **distroless > alpine for runtime.** Smaller image, no shell for attackers to exploit.
- **Document the in-memory limitation.** Don't silently lose player games — tell them in the UI.
