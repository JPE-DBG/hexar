# hexar

## Hexar — Fly.io Deployment Guide

### Prerequisites

Install flyctl: download `flyctl_*_Windows_x86_64.zip` from [github.com/superfly/flyctl/releases](https://github.com/superfly/flyctl/releases/latest), extract `flyctl.exe`, add to PATH.

---

### First-time setup

```bash
# 1. Log in (opens browser)
fly auth login

# 2. Register the app (pick a globally unique name)
fly apps create <your-app-name>

# 3. Update fly.toml line 2 with your chosen name
#    app = '<your-app-name>'

# 4. Deploy (Fly builds remotely from your Dockerfile)
fly deploy

# 5. Verify
curl https://<your-app-name>.fly.dev/health
# → {"status":"ok","version":"<git-sha>"}
```

---

### GitHub Actions auto-deploy

Every push to `main` runs tests then deploys automatically.

**One-time secret setup:**

```bash
# Generate a long-lived deploy token
fly tokens create deploy -x 999999h
```

Copy the output, then in GitHub:
**Repo → Settings → Secrets and variables → Actions → New repository secret**
- Name: `FLY_API_TOKEN`
- Value: *(paste token)*

After that, push to `main` → tests pass → Fly deploys automatically.

---

### Redeploying manually

```bash
fly deploy          # from repo root
fly logs            # stream logs if something looks wrong
fly status          # check machine health
```

---

### Key config notes (`fly.toml`)

| Setting | Value | Why |
|---|---|---|
| `auto_stop_machines` | `false` | **Critical** — Fly sleeping a machine kills active WebSocket connections mid-game |
| `min_machines_running` | `1` | Always one instance ready, no cold-start delay for players |
| `force_https` | `true` | Auto-upgrades `ws://` → `wss://` at Fly's edge |
| Health check path | `/health` | Fly uses this to detect crashes and restart the machine |

---

### Known limitation

Rooms are in-memory — a server restart (deploy or crash) ends all active games. Players see a disconnect overlay and are prompted to return to the lobby. Future fix: persist room state to SQLite.