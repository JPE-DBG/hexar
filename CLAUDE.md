# Hexar — Real-Time Hex Strategy Game

## Overview

Hexar is a fast-paced, real-time multiplayer hex strategy game. Players build a cycling deck of cards and spend a continuously-filling bar resource to play them. Units auto-march across a hex map toward the enemy capital. Target game length: ~15 minutes.

**Core Loop:** Bar fills → play top card of deck (spend bar) → card cycles back → buy new cards from shared shop → units march to enemy capital → first to reduce enemy capital to 0 HP wins

**Genre:** Real-time deck-building strategy | **Players:** 1v1 (2–4 planned) | **Win Condition:** Capital destruction

### Design Pillars

1. **Single Resource:** One currency (bar) pays for everything — playing cards, buying from shop, removing cards. No secondary resources.
2. **Deck as Strategy:** Deck composition is the primary strategic axis. Buying cards and removing cards are as important as playing them.
3. **Readable Combat:** Units are visible on the hex grid, move at human-readable speed, fight when they meet. No hidden rolls.
4. **Meaningful Map:** Distance = reaction time. Claimed hexes = building slots. Map is space and time, not income.

---

## Game Rules

### Bar

- **Fill rate:** 0.1 bar/sec (fixed — does not change with territory or buildings)
- **Scale:** 0–10 (accumulates; capped at 10)
- **Single currency:** bar pays for all actions — playing cards, buying shop cards, removing cards
- Bar never resets; it accumulates until spent or capped

### Deck

Each player has their own ordered deck of cards. The deck is a queue — cards cycle from top to bottom.

**Playing a card:**
- The top card is always visible
- Pay its bar cost → effect resolves immediately → card goes to the bottom of the deck
- All cards cycle back (including building cards — playing a building card again places another copy on a different hex)

**Aside slot:**
- Push the current top card sideways into the aside slot — no bar cost to push
- The aside slot is blocked until you pay the card's cost and play it
- Played aside card goes to the bottom of the deck
- You cannot push a new card while the aside slot is occupied
- Starting aside slots: **1**
- Additional aside slots: unlocked by placing a specific building card (from shop) on a free hex

**Deck lockout prevention:**
- If the top card cannot be played (insufficient bar, or no free hex for a building card), it auto-moves to the bottom after a fixed timer (TBD — single constant, tuned during playtesting)
- Ensures the deck keeps moving in all edge cases

**Deck size:**
- Starting deck: 7 cards (see Starting Deck below)
- No maximum — buying cards from the shop grows the deck
- Larger deck = each card appears less frequently (trade-off: power vs cycle speed)
- Card removal (via shop) permanently removes a card from the deck, thinning it

---

### Starting Deck (7 cards, both players identical)

| Card | Cost | Effect | Qty |
|------|------|--------|-----|
| Hex Claim | 1 bar | Claim one adjacent unclaimed hex (player chooses which) | ×2 |
| +1 Bar Burst | 0 bar | Instantly add +1 to current bar | ×2 |
| +15% Bar Speed | 1 bar | Bar fill rate ×1.15 for 15s (stacks additively with other active copies) | ×2 |
| Basic Soldier | 2 bar | Deploy a soldier unit toward the enemy capital | ×1 |

---

### Map

**Claimed hexes:**
- Each player starts with their capital hex claimed
- Hex Claim cards expand territory to adjacent unclaimed hexes (player chooses which)
- Claimed hexes are building slots — one building per hex
- No free claimed hex = building cards are unplayable (auto-cycle via lockout timer)
- Players keep hex ownership even if the building on it is destroyed
- A destroyed building frees the hex for a new building

**Unit movement:**
- Units march hex-by-hex toward the enemy capital via shortest path
- Units fight enemy soldiers and towers they encounter en route
- Units ignore non-combat buildings (MVP)
- Units do not capture hexes they pass through (MVP)

---

### Units

#### Basic Soldier

| Stat | Value |
|------|-------|
| Bar cost | 2 |
| HP | 2 |
| Attack power | 1 |
| Attack speed | 1 hit/sec |

**Combat resolution (simultaneous attacks):**

*1v1:* Both attack once/sec. Equal soldiers die after 2 seconds.

*2 attackers vs 1 defender:*
- Defender takes 2 damage/sec → dies in 1 second
- Defender attacks one attacker: 1 damage before dying
- Result: one attacker at 1 HP, one at 2 HP — both continue to capital

*Soldier vs Tower:*
- Tower fires at soldiers within its range (stats vary by tower card)
- Soldier attacks tower on contact
- Soldier may die to tower fire before reaching it, or survive and destroy the tower

*Soldier vs Capital:*
- Soldier reaches capital hex → deals 1 HP damage to capital → soldier consumed
- Capital does not fight back (no built-in defense)

---

### Shop

- **Shared pool** visible to both players at all times
- **Dominion-style:** fixed set of available card types, each with limited quantity
- **Buy at any time:** pay bar cost → new card goes to the bottom of your deck
- **Competition:** both players draw from the same pool — first to buy gets the card
- **Card removal:** shop service (not a card) — permanently removes your current top card or aside card from your deck; mid-to-high bar cost (TBD)

**Buying trade-offs:**
- Growing the deck makes each card appear less often
- High-cost purchases drain bar, slowing play
- Removing weak starter cards thins the deck, cycling to power cards faster

#### Shop Card Pool (TBD — design session required before implementation)

Card types confirmed for the shop:

| Type | Category | Notes |
|------|----------|-------|
| Tower | Building / Military | Fights soldiers within range; HP and attack stats TBD |
| Bar Speed Building | Building / Economy | Permanent bar fill rate boost on its hex; bonus value TBD |
| Spawner Building | Building / Military | Produces a unit every X seconds off-deck; high cost; stats TBD |
| Aside Slot Building | Building / Utility | Permanently adds 1 aside slot; requires free hex |
| Additional unit cards | Unit | Faster, tankier, or special-behavior variants; TBD |
| Capital Heal | Misc | Restores capital HP; TBD |

Exact card costs, stats, and quantities are TBD — to be designed before shop implementation.

---

### Win Condition

- **Capital HP:** 20 (subject to playtesting)
- **No built-in defense:** capital is a pure HP counter
- **Permanent damage:** no HP regeneration (capital heal may be a shop card — TBD)
- When capital HP reaches 0, that player loses immediately

---

## Server Architecture

### Tech Stack

| Layer | Choice | Future |
|---|---|---|
| Server | Go | — |
| Client | TypeScript | — |
| Rendering | HTML Canvas + DOM overlays | — |
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
│   │   ├── state.go            # GameState, HexState, Player, Deck, Unit structs
│   │   ├── tick.go             # RunTick(state, dt) — phase function
│   │   ├── bar.go              # bar fill accumulation
│   │   ├── deck.go             # deck queue, aside slot, card play, auto-cycle
│   │   ├── units.go            # unit march, combat resolution
│   │   ├── shop.go             # shared shop pool, buy, card removal
│   │   ├── victory.go          # capital HP win condition
│   │   ├── action.go           # Action types, validate+apply dispatch
│   │   ├── hexmath.go          # axial coords, adjacency, distance
│   │   └── constants.go        # ALL numeric constants
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
│       ├── render/              # Canvas hex grid, unit tokens, buildings
│       ├── input/               # Pixel→hex detection, action dispatch
│       ├── ui/                  # DOM: HUD (bar meter), deck panel, shop panel
│       ├── hexmath.ts           # Axial math (mirrors server)
│       └── constants.ts         # Mirrors server constants
└── .claude/skills/
```

**Key rule:** `internal/game/` has zero imports outside stdlib. `RunTick(state, dt)` is a pure function — fully testable with `go test` alone.

### Tick Execution Order

6 phases, executed strictly in sequence every 100ms:

```
Phase 1: PROCESS ACTIONS — drain queued player actions (PlayCard, BuyCard, PushAside, ClaimHex)
Phase 2: BAR            — accrue bar for each player (×0.1 per tick)
Phase 3: UNITS          — advance unit positions, resolve combat, deal capital damage
Phase 4: DECK           — advance auto-cycle timers, move stuck cards to bottom
Phase 5: VICTORY        — check capital HP, trigger win
Phase 6: DELTA          — diff vs previous tick, broadcast to clients
```

### State Sync

- **Snapshot:** Full state on initial connect and reconnect (`prevSnapshot = nil` forces full send)
- **Delta:** Only changed fields per tick. Target: < 500 bytes/tick at idle
- Deck state: always sent in full (order matters; partial delta would corrupt queue)
- Bar: always sent (continuously changing; omitting would leave stale value)
- Units: always sent — reverts to empty when all units die; omitting would leave stale units on client

### Session & Connection

**Lobby flow:** Same as MVP1 — `POST /lobby/create`, `POST /lobby/join`, `GET /ws?code=&token=`

**Session persistence:** URL hash `/#CODE:TOKEN` — same reconnect mechanism as MVP1.

**Pause on disconnect:** Same forfeit budget (120s cumulative). Loop ticks continue but `RunTick` skipped while paused.

**Duplicate tab:** Server closes new WS with code 4001.

---

## Client & UI

### Layout (TBD — graphical design session required before implementation)

The deck-building UI is fundamentally different from the MVP1 sidebar layout. A dedicated graphical design session must produce a layout spec before client code is written. Key decisions pending:

- Overall screen split (hex map area vs deck/shop panels)
- Deck queue visualization (top card prominent, aside slot, card backs visible behind)
- Bar fill visualization (strip, meter, or edge)
- Shop panel (always visible or pop-out)
- Card anatomy (name, cost, effect text)
- Unit representation on hex grid (token moving hex-by-hex)

Output of design session: `.claude/design/ui-layout-deck-building.md`

### Canvas + DOM Split

- **Canvas (`renderer.ts`):** Hex grid, unit tokens marching hex-by-hex, buildings on hexes, effects
- **DOM overlays:** HUD (bar meter), deck panel (top card + aside slot), shop panel, lobby/victory screens

### Known Issues (non-blocking, document before touching)

None yet — client rewrite not started.

---

## Testing

**Policy:** Update CLAUDE.md first (design spec), implement the feature, then add tests. Tests validate the spec, not the implementation.

### Layer 1: Game Logic (`internal/game/`)

Pure functions, stdlib only, no goroutines, table-driven. `package game`.

| File | Test cases |
|---|---|
| `bar_test.go` | Bar accrues at 0.1/tick; caps at 10; bar correctly deducted on PlayCard; bar correctly deducted on BuyCard |
| `deck_test.go` | Play from top cycles card to bottom; aside slot blocks on push; aside slot card plays and goes to bottom; cannot push when slot occupied; auto-cycle timer moves stuck card to bottom; starting deck initialized with correct 7 cards |
| `units_test.go` | Unit march advances toward capital; 1v1 combat both die after 2s; 2v1 correct HP after fight; unit consumed on capital contact; capital HP decrements |
| `shop_test.go` | BuyCard deducts bar; bought card goes to bottom of deck; shared pool quantity decrements; second player blocked when quantity 0; card removal permanently removes card from deck |
| `victory_test.go` | Capital at 0 HP sets WinReason; damage is permanent (no reset between ticks) |

Run: `go test ./internal/game/...`

### Layer 2: Room Integration (`internal/room/`)

Real goroutines, minimal `time.Sleep`. Same patterns and test structure as MVP1.

### Layer 3: Playwright E2E (`e2e/`)

Two `BrowserContext` objects, `webServer` config auto-starts both servers. Test specs written after client implementation.

Run all tests: `go test ./... && make test-e2e`

---

## Deployment

### Infrastructure

- **Host:** Fly.io (`hexar.fly.dev`)
- **Dockerfile:** Multi-stage — `node:20-alpine` (Vite build) → `golang:1.23-alpine` (server build, `CGO_ENABLED=0 GOOS=linux`) → `gcr.io/distroless/static-debian12` (runtime)
- **`fly.toml`:** `auto_stop_machines = false` (CRITICAL — Fly default sleep kills active WebSocket connections), `min_machines_running = 1`, `force_https = true`
- **Health endpoint:** `GET /health` → `{"status":"ok","version":"..."}`

### CI/CD Pipeline (GitHub Actions)

`.github/workflows/deploy.yml` — deploys on push to `main` when code paths change. Add `[skip deploy]` to commit message to skip deploy while running tests.

---

## Current Status

**MVP1:** Complete on `release_mvp_1` branch (Power gate combat system). Deployed to Fly.io.

**Active branch (`mvp_fix_combat_2026_06_13`):** Redesigning combat to deck-building system.

### Implementation Progress

| Step | Status |
|------|--------|
| Design spec (deck-building concept) | ✅ `.claude/design/combat-deck-building-concept.md` |
| CLAUDE.md rewrite | ✅ Complete |
| Graphical design (UI layout spec) | ✅ `.claude/design/ui-layout-deck-building.md` |
| Delete old game logic / stub new state | ✅ Complete — new `internal/game/` compiles, all room tests pass |
| Core backend (bar, deck, units, win) | ⬜ Not started |
| Shop card pool design | ⬜ Not started |
| Shop backend implementation | ⬜ Not started |
| Client rewrite | ⬜ Not started |
| Playtest + iterate | ⬜ Not started |

### Open Questions (TBD)

- [ ] Shop card pool: which cards, costs, quantities
- [ ] Tower stats (HP, range, fire rate)
- [ ] Spawner unit type and spawn interval
- [ ] Bar speed building bonus value
- [ ] Auto-cycle lockout timer value
- [ ] Map size and capital distance
- [ ] Capital damage per soldier (likely 1)
- [ ] Bar scale cap (is 0–10 correct?)
- [ ] +15% bar stacking: additive confirmed, exact formula TBD

### Future Mechanics (Post-MVP)

- [ ] **4-player support** (2v2 with shared resources; MessagePack for bandwidth)
- [ ] **SQLite match history** (`modernc.org/sqlite`)
- [ ] **OAuth authentication** + ranked ladder
- [ ] **Special hex types** (terrain effects on unit speed)
- [ ] **Additional unit and building card expansion**

---

## Rejected Alternatives

### Combat System

- **Power gate (MVP1 — `release_mvp_1`):** Permanent stalemates when both players stack Power buildings at borders. Neither player can break through once power is matched.
- **Action Bar + Roguelike (Variant E):** Bar fills from hex count + economy buildings − tower maintenance. Research points as second resource. Roguelike stat-bonus picks. Too many parallel systems running simultaneously (bar + research points + territory income + pick management). Replaced by single-resource deck-building.

### Infrastructure

- **Node.js server:** Go developer; worse concurrency model for tick loops
- **Phaser/PixiJS:** Overkill for colored hexagons + text; adds framework weight
- **React/Vue/Svelte (for MVP):** Ceremony for "draw hexagons, update numbers" UI
- **WebRTC:** Massive complexity for P2P; WebSocket + 100ms ticks is fine for Hexar
- **Serverless (Lambda):** Can't maintain WebSocket + tick loop — wrong model for real-time games

---

## Feedback Log

**Purpose:** Track what changed during beta testing and why.

### Beta Cycle 1 (May 2026) — MVP1 Power Gate System

Full feedback log preserved on `release_mvp_1` branch CLAUDE.md. MVP1 testing revealed the stalemate problem that motivated the deck-building redesign.
