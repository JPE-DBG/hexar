# Hexar — Real-Time Hex Strategy Game

## Overview

Hexar is a fast-paced, real-time multiplayer hex strategy game inspired by Antiyoy. Players conquer hexes through economic growth and strategic upgrades on a small hexagonal map. Target game length: 25–30 minutes.

**Core Loop:** Earn gold → Upgrade buildings → Gain power → Conquer adjacent hexes → Repeat

**Genre:** RTS economic conquest | **Players:** 1v1 (2–4 planned) | **Win Condition:** Capital Capture

### Design Pillars

1. **Economic Competition:** Victory through resource management and smart upgrades, not reflexes
2. **Meaningful Decisions:** Each building choice and tech unlock trades off against alternatives
3. **Tense Endgame:** Border warfare and tight resource battles in final minutes
4. **No Snowballing:** Stepped maintenance costs keep powerful players vulnerable

---

## Game Rules

### Hexes & Ownership

- **Unclaimed:** No owner, Power 0, free to claim for 10 gold (instant if Attacker Power ≥ 1)
- **Owned:** Controlled by a player; generates 2 gold/sec base income; can hold one building
- **Capital:** Each player starts with one capital hex. Innate Power 1 (no building needed). Losing it ends the game.

**Maintenance costs (stepped):**
| Hex range | Cost/sec each |
|---|---|
| 1–10 | 1 |
| 11–20 | 2 |
| 21+ | 3 |

**Why stepped:** Flat 1/sec maintenance can never exceed 2/sec base income, so auto-drop would be unreachable. Stepped costs force Economy building investment as territory scales.

**Auto-drop rules:**
- When net income goes negative, a 10-second grace period begins (20s with Resilience tech)
- Player chooses which hex to drop; UI shows a prompt with the timer
- If no choice made: server drops the hex with lowest income; ties broken by lowest building investment; capital and battle-active hexes are protected
- On drop: 50% building refund (70% with Resilience); hex becomes unclaimed instantly

**Maintenance math example:**
```
30 hexes: maintenance 60/sec vs income 60/sec → net 0/sec (needs Economy buildings)
31 hexes: maintenance 63/sec vs income 62/sec → net -1/sec (auto-drop triggers)

With 10 Economy L1 buildings on 30 hexes:
  income = 10 × 3.9/sec + 20 × 2/sec = 79/sec → net +19/sec (sustainable)
```

---

### Resources

- **Gold:** Main currency — buildings, upgrades, attacks, counter-spend
- **Tech Points (TP):** Earned from Research buildings (0.2 TP/sec per level); spent on tech tree

---

### Buildings

One building per hex. No separate "build" action — upgrading an empty hex to L1 is the first step. Buildings can be demolished for a 50% refund and replaced.

**Upgrade cost formula:** `BuildCost × 2^currentLevel`

| Building | BuildCost | Effect | Income formula |
|---|---|---|---|
| Economy (Gold) | 60g | +income from this hex | `(2 + 0.6 × level) × 1.5` → L1: 3.9/sec |
| Power | 60g | +1 Power per level | — |
| Research | 80g | +0.2 TP/sec per level | — |

**Economy building costs:** L1=60, L2=120, L3=240, L4=480  
**Power building costs:** L1=60, L2=120, L3=240, L4=480  
**Research building costs:** L1=80, L2=160, L3=320, L4=640

**Tech bonus stacking on Economy:**
```
No techs:             (2 + 0.6 × level) × 1.5       → L1: 3.9/sec
+ Prosperity:         (2 + 0.6 × level) × 1.5 + 1.0 → L1: 4.9/sec
+ Compound Growth:    (2 + 0.6 × level) × 1.5 × 1.25→ L1: 4.875/sec
+ Both:               (2 + 0.6 × level) × 1.5 × 1.25 + 1.0 → L1: 5.875/sec
```

**Demolish refund:** `BuildCost × (2^level - 1) × 0.5`  
Examples: Gold L1 → 30g; Gold L2 → 90g; Gold L3 → 210g

---

### Combat

**Claiming unclaimed hex:** 10 gold, instant (attacker needs Power ≥ 1). No battle.

**Attacking enemy hex:**
- Cost: 100 gold (only attacker pays)
- Requirement: Attacker Power > Defender Power (strictly greater; ≤ 0 diff = attack fails, no gold spent)
- Power diff > 3 → instant takeover (no timer)
- Power diff 1–3 → standard battle

**Battle duration:** `5 + (AttackerPower + DefenderPower) / 2` seconds (min 6s, max ~15s)

Winner determined when timer reaches 0 — whoever has higher Power at that moment. Attacker cannot cancel after committing 100 gold.

**Counter-spend (defender active defense):**
- Cost: 50 gold → +1 Power for remainder of battle
- Cap: `min(+3, seconds remaining)` — time pressure limits how much defender can buy
- **Garrison tech:** Each adjacent owned hex adds +1 passive (cap +2). Garrison and counter-spend share the +3 total cap.
- **Siege Mastery** (attacker): Timers 40% shorter, compressing defender's reaction window

**Power diff > 3 → instant takeover** (enemy hexes only; unclaimed hexes always instant)

**Garrison threshold in UI:** Effective attack requirement shown as `Defender Power + Garrison bonus`. UI blocks guaranteed-loss attacks; server validates only base power.

---

### Tech Tree

Research buildings generate 0.2 TP/sec per level. Unlock any tech in any order — no prerequisites. **Total tree: 460 TP across 12 techs.** No single game unlocks everything.

| Tech | Cost | Effect | Archetype |
|------|------|--------|-----------|
| Blitz | 20 TP | Unclaimed hex claims cost 0g | Aggressor |
| Fortify | 20 TP | Spend 40g to prevent instant-takeover on one hex for 90s | Defender |
| Prosperity | 25 TP | Each Economy building +1/sec flat bonus | Builder |
| Reclamation | 25 TP | Recapturing a previously owned hex costs 50g | Territorial |
| Vanguard | 30 TP | After capturing an enemy hex, next attacks within 12s cost 50g (timer refreshes on each capture) | Aggressor |
| Garrison | 30 TP | During battle, each adjacent owned hex +1 defense Power (cap +2; total cap +3 shared with counter-spend) | Defender |
| Supply Lines | 40 TP | Maintenance: 0.9/sec (1–10), 1.8/sec (11–20), 2.7/sec (21+) | Builder |
| War Chest | 30 TP | Recover 30g when capturing an enemy hex | Territorial |
| Resilience | 45 TP | Auto-drop grace 10s→20s; drop refund 50%→70% | Defender |
| Iron Grip | 55 TP | All owned hexes permanently +1 Power | Aggressor |
| Compound Growth | 65 TP | Economy buildings ×1.25 output (applied before Prosperity flat bonus) | Builder |
| Siege Mastery | 75 TP | Your battle timers −40% (min 3s); equal-Power ties → attacker wins | Aggressor |

**Tech archetypes:**
- **Blitz Aggressor:** Blitz → Vanguard → Iron Grip — free land-grab, chain attacks, territory-wide Power
- **Economic Builder:** Prosperity → Supply Lines → Compound Growth — wide sustainable territory
- **Fortress Defender:** Garrison → Fortify → Iron Grip → Resilience — make attacking expensive
- **Siege Striker:** Iron Grip → Siege Mastery → Vanguard — fast decisive battles
- **Territorial:** Reclamation → Vanguard → War Chest — fluid borders, gold recovery

**Stacking rules:**
- **Reclamation + Vanguard both active:** Attack cost 100g − 50g − 50g = 0g (free reclaim during Vanguard window)
- **Compound Growth + Prosperity:** Multiplier applied first, then flat bonus: `× 1.25 + 1.0`
- **Iron Grip + Garrison:** Iron Grip always-on (+1 to all hexes, shown in hex label). Garrison defense-only, dynamic (+N def shown separately in UI, excluded from static hex label)

**Research investment guide:**
- 1 Research L1: ~180 TP in 15 min → ~5 cheap techs
- 2 Research L1: ~360 TP in 15 min → ~8 techs
- 3 Research L1: ~540 TP in 15 min → ~11 techs (full tree by min 13)

---

### Victory

**Capital Capture:** Capture the enemy capital → win immediately. All loser hexes become unclaimed instantly.

- Capital is captured like any enemy hex (battle or instant if diff > 3)
- Captured capital becomes a normal hex for the winner (loses innate Power 1)
- `state.WinReason`: `"capital"` or `"forfeit"` (disconnect budget exhausted)

---

## Server Architecture

### Tech Stack

| Layer | Choice | Future |
|---|---|---|
| Server | Go | — |
| Client | TypeScript | — |
| Rendering | HTML Canvas + DOM overlays | Svelte for complex UI |
| Bundler | Vite | — |
| Networking | WebSocket, JSON | MessagePack for 4-player |
| WebSocket lib | `github.com/coder/websocket` | — |
| Game loop | Goroutine + `time.Ticker` (100ms) | — |
| Database | None (in-memory) | SQLite for match history |
| Auth | None | OAuth when accounts added |

### Project Structure

```
hexar/
├── CLAUDE.md
├── go.mod
├── cmd/server/main.go           # HTTP+WS entry point, serves client
├── internal/
│   ├── game/                    # PURE logic — zero I/O, zero network, deterministic
│   │   ├── state.go            # GameState, Hex, Player, Battle structs
│   │   ├── tick.go             # RunTick(state, dt) — 6-phase function
│   │   ├── economy.go          # income, maintenance, auto-drop
│   │   ├── combat.go           # attack validation, battle resolution
│   │   ├── building.go         # build, upgrade, demolish
│   │   ├── tech.go             # unlock validation, bonus computation
│   │   ├── victory.go          # victory condition checks
│   │   ├── action.go           # Action types, validate+apply dispatch
│   │   ├── hexmath.go          # axial coords, adjacency, distance
│   │   └── constants.go        # ALL numeric constants from CLAUDE.md
│   ├── room/                    # Owns GameState + ticker goroutine
│   │   ├── room.go
│   │   └── loop.go
│   ├── net/                     # WebSocket, HTTP, message serialization
│   │   ├── server.go
│   │   ├── client.go
│   │   ├── messages.go
│   │   └── delta.go
│   └── mapgen/                  # Hardcoded map (MVP)
│       └── generator.go
├── client/
│   └── src/
│       ├── main.ts
│       ├── net/                 # WebSocket, reconnect, message types
│       ├── state/               # Game state mirror, apply delta/snapshot
│       ├── render/              # Canvas hex grid, buildings, battles
│       ├── input/               # Pixel→hex detection, action dispatch
│       ├── ui/                  # DOM: HUD, sidebar, tech tree
│       ├── hexmath.ts           # Axial math (mirrors server)
│       └── constants.ts         # Mirrors server constants
└── .claude/skills/
```

**Key rule:** `internal/game/` has zero imports outside stdlib. `RunTick(state, dt)` is a pure function — fully testable with `go test` alone.

### Tick Execution Order

6 phases, executed strictly in sequence every 100ms:

```
Phase 1: PROCESS ACTIONS — drain queued player intentions, validate, apply
Phase 2: ECONOMY      — accrue gold/TP (×0.1 per tick)
Phase 3: BATTLES      — decrement timers 0.1s, resolve expired
Phase 4: AUTO-DROP    — check negative income, manage grace, drop if expired
Phase 5: VICTORY      — check capital capture, unclaim loser hexes
Phase 6: DELTA        — diff vs previous tick, broadcast to clients
```

Actions first = immediate effect. Economy before battles = counter-spend gold deducted before battle resolves. Victory last = reflects true final state.

### State Sync

- **Snapshot:** Full state on initial connect and reconnect (`prevSnapshot = nil` forces full send)
- **Delta:** Only changed fields per tick. Target: < 500 bytes/tick at idle
  - `Tech []bool` omitted when unchanged (only goes false→true, never reverts)
  - `Battles` omitted when empty (client preserves last known battles via `?? state.battles`)
  - `AutoDropActive`, `AutoDropGrace`, `VanguardTimer` always sent — these fields revert to zero/false at runtime; omitting them would leave the client with stale nonzero values via delta merge spread
  - `FortifyTimer` quantized to 1-second boundaries to avoid appearing in every tick

**Delta merge pattern (client):**
```typescript
players.set(String(p.id), { ...(existing ?? {}), ...p, tech: p.tech ?? existing?.tech ?? [] });
```

### Session & Connection

**Lobby flow:**
- `POST /lobby/create` → room code (4-char hex) + session token (32-char hex) + playerID
- `POST /lobby/join` → token + playerID (room full after 2 players)
- `GET /ws?code=…&token=…` → WebSocket connection; server validates token before upgrading

**Session persistence (URL hash):**
- On game start, client writes `/#CODE:TOKEN` into the URL (`history.replaceState`)
- On page load: check sessionStorage first (network-drop), then URL hash (tab-close)
- URL hash is per-tab — two players in the same browser get different hashes and don't collide
- Hash cleared on game end

**Waiting state:**
- Room created → `state.Waiting = true`; game loop goroutine NOT started
- Second player connects → server sets `state.Waiting = false`, starts loop, broadcasts to P1

**Pause on disconnect:**
- Any player disconnects → `state.Paused = true`, `PauseTimeLeft` set to remaining budget
- Loop ticks continue (clients get countdown) but `RunTick` is skipped — economy/battles/auto-drop frozen
- Reconnect → `state.Paused = false`, budget saved to `remainingGrace[pid]`
- Budget exhausted → loop enqueues `ActionForfeit`; next tick processes it

**Forfeit budget:** 120 seconds cumulative per player across the whole game (not reset per disconnect)

**Duplicate tab:** Server closes new WS with code 4001; client shows "Already Connected" overlay. Second `OnConnect` for same player during active game replaces the old client (evicts from activeClients, sends snapshot to new client).

### Edge Cases

| Scenario | Decision |
|---|---|
| Disconnect mid-game | Pause immediately; 120s cumulative budget; exhausted → forfeit |
| Two attacks on same hex | First-in-queue wins; second rejected (one battle per hex) |
| Reconnection | Full snapshot sent (`prevSnapshot = nil`); delta resumes after |
| Capital captured during outgoing attack | Immediate game over, all battles canceled |
| Counter-spend same tick as battle expires | Applied (Phase 1 runs before Phase 3) |
| Gold below 0 | Clamped to 0; actions rejected if insufficient |
| Player trapped (no Power ≥ 1 border) | Must build Power to regain attack ability |
| Both players not yet connected | Loop not started; `state.Waiting = true` |
| Capture flash trigger | Compare `hex.owner` in incoming delta against current `state.hexes` owner. **Never** use `hex.previousOwner` — server sets it on capture and never resets it, so any subsequent delta for that hex (e.g. fortifyTimer tick) would re-fire the flash spuriously |

---

## Client & UI

### Canvas + DOM Split

- **Canvas (`renderer.ts`):** Hex grid, buildings, battles, effects, timers
- **DOM overlays:** HUD, sidebar, lobby/victory screens, tech tree
- Effects are client-side only — never in game state. `addCaptureFlash(hex, color)` and `addFloater(hex, text)` triggered by state transitions in `applyDelta`.

### Power Display Conventions

- Hex canvas labels show **effective combat power** (base + Iron Grip), not raw building level
  - Power L1 + Iron Grip → label "P2"; Capital + Power L1 + Iron Grip → "P3"
  - Non-Power owned hexes show a small yellow power badge if power > 0 (capital innate or Iron Grip)
- **Garrison excluded from hex label** — defense-only, dynamic; shown as `"+N def"` note in sidebar, excluded from static hex label
- `drawHexFill` is intentionally flat — sprite-based 3D/material look planned for a future pass; do not add gradients

### Timer Visualization Conventions

| Timer | Type | Visual |
|---|---|---|
| Fortify (90s) | Hex timer | Lime `#c8ff70` border segments, clockwise, shrink to 0 |
| Battle (6–15s) | Battle timer | Amber pulsing ring + countdown text + boost dots |
| Vanguard (12s) | Player timer | HUD text `⚡Vanguard X.Xs` |

**Hex timers** (state on a specific hex): clockwise shrinking border segments, full ring at max → empty at 0. Fortify is the reference implementation. All future hex-level timers must follow this pattern. Color must be readable on both player colors (teal P1, red P2); lime `#c8ff70` is established for Fortify.

**Player timers** (state on a player, not a hex): HUD text indicator only — no border segments.

**Timer stacking:** Battle timer always renders at full brightness (priority). Passive hex timers (Fortify) dim to 40% opacity, lineWidth 1.5 while a battle is active on that hex.

### Sidebar & Controls

**Pattern:** Left sidebar (180px wide), sticky build tools, dynamic context area

**Sections:**
1. **CONTEXT** (top, dynamic): Upgrade/Attack/Demolish/Sell Hex/Fortify/Counter-spend based on selected hex
2. **BUILD** (bottom, always visible): Economy [Q], Power [W], Research [E], Tech Tree [T]

**Sticky tool behavior:**
- Click tool button or press hotkey → highlights → click hexes to place buildings
- Tool stays selected for batch building
- Toggle to deselect: click again, press hotkey again, or press ESC
- Clicking a hex under active battle falls through to context selection (shows counter-spend)

**Keyboard shortcuts:**
| Key | Action |
|---|---|
| Q / W / E | Select Economy / Power / Research (toggle) |
| T | Open Tech Tree |
| Space | Upgrade selected hex |
| A | Attack selected hex |
| F | Fortify selected hex (requires Fortify tech) |
| C | Counter-spend during battle |
| D | Demolish selected hex → deselects hex |
| X | Sell Hex → deselects hex |
| ESC | Deselect tool |

**Battle restrictions:** No building placement or upgrade on hexes under active battle (server enforces `ErrBattleInProgress`). Demolish allowed — defender may recover gold for counter-spend.

### Responsive Behavior

| Viewport | Layout |
|---|---|
| Desktop (>1024px) | Left sidebar, full labels, context area |
| Mobile landscape (896×414) | Bottom bar 60px, icon + label |
| Mobile portrait (<768px) | Bottom bar 50px, icons only; destructive actions (Demolish/Sell) in separate top-left panel |

**Destructive panel:** `#destructive-panel` lives outside `#sidebar` so `position: fixed` is viewport-relative (backdrop-filter on `#sidebar` would make fixed children relative to it). CSS controls which copy is shown — buttons exist twice in DOM.

### Known Issues (non-blocking, document before touching)

| Issue | Location | Risk |
|---|---|---|
| Destructive DOM duplication | `sidebar.ts` — `this.destructiveEl.innerHTML` mirrors button HTML | Demolish/Sell buttons exist twice; CSS controls visibility. Querying `data-action` finds two nodes. Refactor before adding more destructive actions. |
| Mobile CSS transparency | `style.css` — portrait mode sets sidebar background to `transparent` | No backdrop-filter, no border; relies on canvas fill for contrast. Light hexes may cause readability issues. |

---

## Testing

**Policy:** Update CLAUDE.md first (design spec), implement the feature, then add tests. This ensures tests validate the spec, not the implementation.

### Layer 1: Game Logic (`internal/game/`)

Pure functions, stdlib only, no goroutines, table-driven. `package game` — tests can call unexported functions directly.

| File | Test cases |
|---|---|
| `action_test.go` | Counter-spend cap/time-cap/flip; tech unlock via action deducts TP; TP accumulation from Research buildings; voluntary drop action clears auto-drop flag |
| `building_test.go` | Building upgrade mechanics, costs, demolish refunds |
| `economy_test.go` | Gold accrues at 2/s per hex; stepped maintenance triggers at 10/20 hex boundaries; Prosperity + Compound Growth formula matches CLAUDE.md math |
| `combat_test.go` | Attack validation: insufficient power, insufficient gold, hex already in battle; battle resolution: winner at timer expiry, loser retains on tie (unless Siege Mastery); instant takeover when diff > 3 |
| `tech_test.go` | Tech unlock deducts TP and sets bool; Reclamation + Vanguard both active + conditions met = 0g attack cost |
| `victory_test.go` | Capital capture sets `WinReason = "capital"`, all loser hexes unclaimed; `ForfeitPlayer` sets `WinReason = "forfeit"` |
| `autodrop_test.go` | Negative net income triggers 10s grace; forced drop when grace expires; 50% refund on dropped hex; Resilience extends grace to 20s and refund to 70%; selection by lowest income then lowest invested; capital and battle-active hexes protected |

Run: `go test ./internal/game/...`

### Layer 2: Room Integration (`internal/room/`)

Real goroutines, minimal `time.Sleep` (1–2 tick durations only). `package room` — tests access unexported fields directly (`room.activeClients`, `room.clientPlayer`, `room.remainingGrace`).

| Test | Asserts |
|---|---|
| `TestWaitingState` | `Waiting=true` until 2nd player connects; loop doesn't start with 1 player; `Waiting=false` when 2nd connects |
| `TestWaitingStateNoGoldAccrual` | Player 1 gold unchanged while waiting (loop not running) |
| `TestPauseOnDisconnect` | `state.Paused=true` immediately; `PauseTimeLeft` ≈ 120s; remaining client gets pause delta |
| `TestUnpauseOnReconnect` | `state.Paused=false` on reconnect; grace saved to `remainingGrace` |
| `TestCumulativeGrace` | 2nd disconnect uses saved grace (doesn't reset to 120s); decrements cumulatively |
| `TestForfeitEnqueued` | `PauseTimeLeft`→0 enqueues `ActionForfeit`; next tick → `WinReason="forfeit"` |
| `TestDuplicateConnectReplaces` | 2nd `OnConnect` for same player replaces old client; new client is in `activeClients`; old evicted |
| `TestDuplicateConnectDuringPause` | Reconnect during pause succeeds (not treated as duplicate) |

Run: `go test ./internal/room/...`

### Layer 3: Playwright E2E (`e2e/`)

Multi-player scenarios use two `BrowserContext` objects in one test process. Playwright's `webServer` config auto-starts both Go server (`:8080`) and Vite dev server (`:5173`).

| Spec | Tests |
|---|---|
| `lobby.spec.ts` | Create shows 4-char code; join with bad code shows error; waiting overlay appears; HUD activates when P2 joins |
| `gameplay.spec.ts` | Canvas + HUD visible after both connect; gold counter increases over time; sidebar visible |
| `connection.spec.ts` | Duplicate tab shows "Already Connected"; URL hash reconnects after tab close; pause banner on opponent disconnect |
| `responsive.spec.ts` | Desktop: sidebar left, destructive panel hidden; Mobile portrait: sidebar bottom, destructive panel top-left; Mobile landscape: sidebar + HUD visible |

**Setup:**
```bash
npm install                   # installs @playwright/test from root package.json
npx playwright install chromium
make test-e2e                 # starts both servers, runs e2e/, exits
```

Run all tests: `go test ./... && make test-e2e`

---

## Deployment

### Infrastructure

- **Host:** Fly.io (`hexar.fly.dev`)
- **Dockerfile:** Multi-stage — `node:20-alpine` (Vite build) → `golang:1.23-alpine` (server build, `CGO_ENABLED=0 GOOS=linux`) → `gcr.io/distroless/static-debian12` (runtime; no shell, minimal attack surface)
- **`fly.toml`:** `auto_stop_machines = false` (CRITICAL — Fly default sleep kills active WebSocket connections), `min_machines_running = 1`, `force_https = true` (auto TLS upgrades `ws://` → `wss://`)
- **Health endpoint:** `GET /health` → `{"status":"ok","version":"..."}` — Fly.io uses this for crash detection

### CI/CD Pipeline (GitHub Actions)

`.github/workflows/deploy.yml`:

```yaml
on:
  push:
    branches: [main]
    paths:                 # skip pipeline for doc-only changes
      - 'cmd/**'
      - 'internal/**'
      - 'client/**'
      - 'go.mod'
      - 'go.sum'
      - 'Dockerfile'
      - 'fly.toml'
jobs:
  test: ...               # go test ./...
  deploy:
    needs: test
    if: "!contains(github.event.head_commit.message, '[skip deploy]')"
```

- Add `[skip deploy]` to commit message to skip the deploy job while still running tests
- `FLY_API_TOKEN` stored in GitHub secrets

### Known Limitations

- Rooms are in-memory — a server restart ends all active games. Future fix: SQLite via `modernc.org/sqlite` (pure Go, no CGO) persisting room state as JSON blob.
- WebSocket URL switches automatically: `wss://` in production, `ws://` in dev (`import.meta.env.PROD`)

---

## Current Status & Future Work

### Status: M1–M9 Complete

All milestones shipped and deployed to Fly.io. M1–M6: core game loop, lobby, delta sync. M7: visual polish + sidebar. M8: testing (Go unit + room integration + Playwright E2E). M9: Fly.io deployment + GitHub Actions CI/CD.

### Post-M9 Fixes

| Fix | Trigger |
|---|---|
| Mobile destructive actions moved to top-left panel | On mobile, context actions at bottom bar were obscured or mis-tapped |
| Deselect on demolish | After demolish, hex stayed selected — stale upgrade/attack buttons remained |
| Fortify timer capture flash | Hex re-flashed every 1s when fortified — `previousOwner` never reset by server; fixed by comparing `hex.owner` vs `state.hexes` owner, not wire `previousOwner` |
| Delta quantization for FortifyTimer | Timer changed every 100ms tick → fortified hexes in every delta; now only sent when timer crosses a 1-second boundary |
| Keyboard shortcut deselect (D/X) | After D (Demolish) or X (Sell Hex), hex stayed selected; fixed by adding `selectedHex = null; renderer.setSelected(null); sidebar.hide()` to both keyboard handlers |
| Delta packet size optimization | Idle game sent 526 bytes/tick (>500 limit); fixed by diffing `Tech []bool` in `buildDelta` (omit when unchanged) and making `Battles` omitempty |
| omitempty on reverting PlayerDTO fields | `AutoDropActive/Grace/VanguardTimer` with omitempty caused client to preserve stale nonzero values via delta spread when fields went to zero; removed omitempty, always send these fields |
| CI/CD path filters | Pipeline ran on every commit including doc-only changes; added `paths:` filter + `[skip deploy]` convention |

### Open Questions (Playtesting)

- [ ] Does the exponential cost curve feel right? Too fast/slow?
- [ ] Is the +3 counter-spend cap balanced? Does it create interesting battles?
- [ ] Does the 70-hex map hit the 30-min game length target?
- [ ] Do all tech archetypes appear in practice, or do 2–3 dominate?
- [ ] Choke point deadlock: narrow maps allow a single high-Power hex to block all expansion. Map generation must avoid single-hex corridors, or add a flanking mechanic.
- [ ] Vanguard stacking: timer refreshes on each capture. Does this create an unstoppable snowball?

### Future Mechanics (Post-MVP)

- [ ] **Tech Tree Benefit Display:** Show quantitative benefits (e.g., "Prosperity: +3.0 gold/s total with 3 Economy buildings"). Deferred — current static descriptions sufficient.
- [ ] **4-player support** (2v2 with shared resources; MessagePack for bandwidth)
- [ ] **SQLite match history** (`modernc.org/sqlite` for room state persistence)
- [ ] **OAuth authentication** + ranked ladder + cosmetics
- [ ] **Special hex types** (mountains, water, resource nodes)
- [ ] **Sprite-based hex rendering** (replace flat fill with 3D/material look; `drawHexFill` intentionally minimal until then)

---

## Rejected Alternatives

- **Node.js server:** Go developer; worse concurrency model for tick loops
- **Phaser/PixiJS:** Overkill for colored hexagons + text; adds framework weight
- **React/Vue/Svelte (for MVP):** Ceremony for "draw hexagons, update numbers" UI
- **WebRTC:** Massive complexity for P2P; WebSocket + 100ms ticks is fine for Hexar
- **Serverless (Lambda):** Can't maintain WebSocket + tick loop — wrong model for real-time games
- **Player-hosted P2P:** Trust issues, NAT traversal, cheating risk
- **Radial context menu:** Occludes adjacent hexes during battles when Garrison bonuses are visible
- **Bottom build menu:** 40% slower for batch building (10 actions vs 6 for placing 5 buildings)
- **Neutral hex state (between unclaimed/owned):** Adds UI complexity with no gameplay benefit; matches Antiyoy's simpler 2-state model
