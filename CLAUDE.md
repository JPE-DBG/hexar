# Hexar — Real-Time Hex Strategy Game

## Overview

Hexar is a fast-paced, real-time multiplayer hex strategy game inspired by Antiyoy. Players conquer hexes through economic growth and strategic upgrades, competing on small maps with 30-minute target game length.

**Core Loop:** Earn resources → Upgrade buildings → Gain power → Conquer adjacent hexes → Repeat

---

## Game Concept

**Genre:** Real-time strategy (RTS), economic conquest  
**Players:** 2-4 (starting with 1v1)  
**Game Duration:** ~25-30 minutes  
**Map:** Hexagonal grid, 60-80 total hexes per 1v1 map (players start with 1 hex, expand to 25-35)  
**Win Conditions:** Conquest, Tech Dominance, or Time Limit (see Victory Conditions)

### Design Pillars

1. **Economic Competition:** Victory through resource management and smart upgrades, not reflexes
2. **Meaningful Decisions:** Each building choice and tech unlock trades off against alternatives
3. **Tense Endgame:** Border warfare and tight resource battles in final minutes
4. **No Snowballing:** Maintenance costs keep powerful players vulnerable

---

## Core Mechanics

### Hexes & Ownership

- **Player hexes (owned):** Controlled by a player, generates income, has buildings
- **Unclaimed hexes (empty):** No owner, Power 0, no buildings
- **Starting position:** Each player starts with exactly 1 hex (their capital)
- **Capital hex:** Has innate Power 1 (no building needed). All other owned hexes start at Power 0 until a Defense building is placed.
- **Owned hex generates:** 2 resources/sec (base)
- **Maintenance cost (Stepped):**
  - Hexes 1-10: 1 resource/sec each
  - Hexes 11-20: 2 resources/sec each
  - Hexes 21+: 3 resources/sec each
- **Why stepped:** Flat maintenance (1/sec) can never exceed base income (2/sec), making auto-drop unreachable. Stepped costs force Economy building investment to sustain large territories.

**Maintenance math:**
```
10 hexes:  maintenance 10/sec vs income 20/sec  → net +10/sec (healthy)
20 hexes:  maintenance 30/sec vs income 40/sec  → net +10/sec (requires some Economy)
30 hexes:  maintenance 60/sec vs income 60/sec  → net  0/sec  (needs Economy buildings)
31 hexes:  maintenance 63/sec vs income 62/sec  → net  -1/sec (auto-drop triggers!)

With 10 Economy buildings (L1) on 30 hexes:
  income = 10 × 3.9/sec + 20 × 2/sec = 79/sec → net +19/sec (sustainable)
```

**Hex auto-drop rules:**
- When net income goes negative, the player must shed hexes until income is positive
- Player chooses which hex to drop (UI prompt, 10-second grace period)
- If no choice made, the hex with the lowest income drops automatically
- On drop: player receives **50% refund of the building cost** on that hex (if any building present)
- Dropped hex becomes unclaimed instantly (enemy can grab it for 10 gold)
- Net income per hex +1/sec (before buildings/upgrades)

### Resources

- **Gold:** Main currency for buildings, upgrades, and attacks
- **Tech Points (TP):** Earned from Research buildings (0.1 TP/sec per Research level), spent on tech tree

**Tech bonus stacking (additive):** All bonuses to resource generation stack additively inside one multiplier:
```
hex_income = base × (1 + economy_bonus + tech_bonuses)

Example: Economy building (+50%) + Production Boom (+30%):
  = 2 × (1 + 0.50 + 0.30) = 2 × 1.80 = 3.60/sec
```

### Buildings (One per Hex)

Each hex can have **one building** of three types. There is no separate "build" action — upgrading an empty hex to L1 is the first upgrade step. Building can be demolished (refund 50%) and replaced.

#### Economy Building
- **Effect:** +50% resources/sec from that hex (stacks additively with tech bonuses)
- **Upgrade cost:** L1=60, L2=120, L3=240, L4=480 (formula: `BuildCost × 2^level`, where level is current level before upgrade)
- **Reward (Compounding):** +0.6 resources/sec per level, multiplied by the +50% bonus
- **Formula:** `(base + 0.6 × level) × 1.5` → Level 1: 3.9/sec, Level 5: (2 + 3.0) × 1.5 = 7.5/sec
- **Max level:** Unlimited (incremental)

#### Defense Building
- **Effect:** +1 Power per level (Power = level, so L1=1, L2=2, etc.)
- **Upgrade cost:** L1=60, L2=120, L3=240, L4=480 (formula: `BuildCost × 2^level`)
- **Max level:** Unlimited (incremental)

#### Research Building
- **Effect:** +0.1 TP/sec per level (fuel for tech tree)
- **Upgrade cost:** L1=80, L2=160, L3=320, L4=640 (formula: `BuildCost × 2^level`, BuildCost=80)
- **Reward (Linear):** +0.1 TP/sec per level (constant gain)
- **Max level:** Unlimited (incremental)

---

## Combat System

### Attack Rules

- **Adjacent only:** Can only attack hexes touching your hex
- **Unclaimed hex cost:** 10 gold (instant takeover, no battle)
- **Enemy hex cost:** 100 gold (triggers battle)
- **Requirement vs enemy:** Attacker Power > Defender Power (strictly greater, ≤ fails with no cost)

### Battle Resolution

**Attacking Empty Hex (Power 0):**
- Cost: 10 gold → Instant takeover, no battle
- Requirement: Attacker Power ≥ 1 (capital hex has innate Power 1, no building needed)
- Reason: Early land-grab is fast; gold bottleneck only matters for enemy combat

**Attacking Enemy Hex (Owned):**
- Cost: 100 gold (only attacker pays)
- **Requirement:** Attacker Power > Defender Power (strictly greater)
- If Power diff > 3: **Instant takeover** (no battle timer, no attrition)
- If Power diff = 1, 2, or 3: **Standard battle** (6-15 sec, see below)
- If Power diff ≤ 0: **Attack fails** (no gold spent, can retry)

**Standard Battle (1-3 power difference):**

```
Battle Duration = 5 + (Attacker Power + Defender Power) / 2 seconds

Min: 6 seconds (Power 1 vs 1)
Max: ~15 seconds (Power 10 vs 8)
```

- Battle timer visible to both players during countdown
- Winner is determined when timer reaches 0 — whoever has higher Power at that moment wins
- Attacker **cannot cancel** once attack is committed (100 gold is spent, no refund)
- Defender counter-spending mid-battle can flip the outcome, but neither side exits early
- Attacker takes hex if they win when timer expires; defender keeps hex if they win

### Defender Options (Active Defense)

**Counter-Spend (Capped)**

During the battle countdown, defender can spend resources to boost Defense Power:
- Cost: 50 gold/second during battle
- Effect: +1 Power per second spent
- **Cap:** min(+3 power, seconds remaining in battle) — cannot boost more than time allows
- **Garrison interaction:** If Garrison tech is unlocked, adjacent owned hexes each add +1 Power automatically (cap: +3 total). Counter-spend and Garrison share this cap — combined boost cannot exceed +3.
- Example: Battle at 9 seconds, you're losing 3 vs 5. Spend 150 gold over 3 sec → Power jumps to 6, you win
- Example: Only 1 second left in battle → max +1 Power boost (50 gold), even if you have gold to spare
- Example (Garrison): 2 adjacent owned hexes give +2 passive. Counter-spend max is now +1 (cap already at +3 with 2 spent)
- Cap prevents defender from "buying" complete victory; time pressure makes decisions tense

---

## Tech Tree

Research buildings generate Tech Points. Spend TP to unlock perks (global bonuses):

| Tech | Cost | Effect |
|------|------|--------|
| Iron Grip | 50 TP | All hexes +1 Power |
| Production Boom | 40 TP | All hexes +30% resource generation |
| Efficient Conquest | 35 TP | Attack cost reduced to 75 gold |
| Garrison | 60 TP | During battle, each adjacent owned hex adds +1 Power to defense (cap: +3 total, shared with counter-spend) |

**Tech Level:** Total number of techs unlocked. Reaching Tech Level 4 (all techs unlocked) triggers Tech Dominance victory when combined with map control (see Victory Conditions).

**Rationale:** Costs reduced by ~40% from original to make tech tree achievable in 30-min games. First few techs are cheap to encourage early tech investment as a viable alternative to pure military.

---

## Economy Example (Early Game)

```
T=0s:      Control 1 hex (capital), Power 1 (innate), no buildings
           Income: 2/sec, Maintenance: 1/sec → Net +1/sec
           Gold: 0

T=0-10s:   Accumulate 10 gold → Claim adjacent unclaimed hex (10g, instant)
T=10-15s:  2 hexes, +2/sec net → Claim hex 3 (5 sec)
T=15-18s:  3 hexes, +3/sec net → Claim hex 4 (3 sec)
T=18-20s:  4 hexes, +4/sec net → Claim hex 5 (2 sec)
T=~2 min:  ~10 hexes claimed, land-grab phase slows as map fills

T=2 min:   Save 60 gold → Upgrade hex 1 to Economy L1 (~6 sec at +10/sec)
           Income: hex 1 = 3.9/sec, rest = 2/sec each → Total: 21.9/sec gross
           Maintenance: 10/sec → Net +11.9/sec

T=2.5 min: Save 60 gold → Build Defense on border hex (~5 sec at +11/sec)
           Spend 30 → Upgrade Defense L1. Power = 2 on that hex.

T=3 min:   First enemy contact at borders
           Save 100 gold for first enemy hex attack (~9 sec at +11/sec)
           Power 2 vs enemy Power 1 → Standard battle (6.5 sec), you win

T=5 min+:  Border warfare begins in earnest
           ~300-400 gold saved for upgrades and attacks
           Power 2-3 on borders, Economy building on 2-3 hexes
```

---

## Victory Conditions

### 1. Conquest Victory (Primary)
- Hold **60% of map for 10 consecutive seconds**
- Timer resets if you drop below 60%
- This is the most common win condition

**Why 60%:** In 1v1, 50% means each player holds half the map — a draw state, not a decisive lead. 60% requires genuinely dominating your opponent, not just tying.

**Strategic implications:** Aggressive players win by pushing early; defensive players must hold the line and counter-attack to reset timer.

**How to stop opponent:**
- Attack their border hexes aggressively
- Bring them below 60%, reset their timer
- Race to 60% yourself

### 2. Tech Dominance (Mid-Game Alternative)
- Reach **Tech Level 4 (all techs unlocked) AND hold 35% map simultaneously** for 10 consecutive seconds
- Rewards investing in research as a viable win condition

**Strategic implications:** Tech rush is faster than pure conquest (tech scaling matters). Early investment in Research building pays off. Opponent can counter by military pressure.

**How to stop opponent:**
- Rush military early while they tech
- Attack their Research hexes specifically
- Keep them pinned defending, prevent expansion
- Example: While opponent reaches Tech 2-3, you've conquered 40% through aggression, win by Conquest

### 3. Time Limit (Tiebreaker)
- At 30 minutes, highest hex count wins
- Rarely triggers in well-balanced games

---

## Balance Rules & Constraints

### Upgrade Costs (Unified formula: `BuildCost × 2^currentLevel`)
- **Economy (BuildCost=60):** L1=60, L2=120, L3=240, L4=480
- **Defense (BuildCost=60):** L1=60, L2=120, L3=240, L4=480 (same curve as Economy)
- **Research (BuildCost=80):** L1=80, L2=160, L3=320, L4=640
- There is no separate "build" action — upgrading empty hex to L1 costs `BuildCost × 2^0 = BuildCost`
- **Economy rewards:** +0.6/sec per level × 1.5 multiplier = +0.9/sec net gain per level
- **Defense/Research rewards:** +1 Power per level; +0.1 TP/sec per level (constant)
- **Demolish refund:** 50% of total invested; for all buildings: TotalInvested = `BuildCost × (2^level - 1)`, refund = `BuildCost × (2^level - 1) × 0.5`
  - Economy L1: TotalInvested=60, refund=30; L2: TotalInvested=180, refund=90; L3: TotalInvested=420, refund=210

### Maintenance System (Prevents Extreme Expansion)
- Each hex costs 1 maintenance/sec to hold
- Forces quality over quantity
- At ~10+ hexes, income caps out without production upgrades
- If maintenance exceeds income, slowest hexes auto-drop (player chooses order)

### Tech Scaling (Achievable in 30 min)
- Tech costs reduced for viability (all 4 techs = 185 TP total, achievable with 2-3 Research hexes by mid-game)
- Tech progression is a viable win condition (Tech Level 4 + 35% map achievable in ~15-18 min with Research investment)
- Opponents must balance military pressure with allowing tech growth

### Early Game Parity
- All players start with 1 capital hex, Power 1, no buildings
- First minute: rapid hex expansion (unclaimed hexes taken instantly)
- 5-10 minute mark: players meet at borders
- Game stays competitive until ~15-minute mark if balanced play

### Capital Hex Rules

- **If capital is captured:** Player immediately loses the game. All their hexes become unclaimed instantly.
- **Captured capital:** Becomes a normal hex for the conqueror (no innate Power 1, no special rules)
- **Building on capital:** Owner can place Defense buildings on their capital (Power stacks on top of innate Power 1, e.g., Defense L2 = Power 3 total)
- **Strategic implication:** Capital is a high-value target — losing it ends the game, so defending it is critical. Attacker targeting capital forces defender to split attention between borders and home base.

---

### No "Neutral" Hex State (Only Unclaimed)
- **Simplified model:** Hexes are either **owned by a player** or **unclaimed (empty)**
- **Dropped "neutral" concept:** Initial design had unclaimed hexes revert to "neutral" after conquest, but this adds complexity without benefit
- **Why:** 
  - When you attack an enemy hex and win, you capture it directly (you own it)
  - No intermediate "neutral" state where other players can immediately reclaim it
  - Simpler UI: only 2 states instead of 3
  - Matches Antiyoy's design (conquered = owned)

### Instant Takeover on Unclaimed Hexes Only
- Players can instantly claim any adjacent unclaimed hex if they have Power ≥ 1
- Conquering **enemy hexes** always requires battles (6-15 sec or instant if diff > 3)
- Effect: Early game is about **speed of expansion** into unclaimed territory. Midgame is about **tactical battles** over enemy hexes.

### Adjacent-Only Conquest
- Reinforces territorial control and front-line battles
- Prevents "teleport attacks" across map
- Creates natural borders and defensive positions

### Only Attacker Pays
- Defender can react (counter-spend, capped at +3) rather than commit upfront
- Asymmetric costs reward smart positioning
- Encourages aggressive expansion but punishes poor decisions

### Power Difference Auto-Takeover (diff > 3, Enemy Hexes Only)
- Instant conquests feel rewarding (massive advantage in hex battles)
- Close fights (1-3 diff) are tense and tactical
- Prevents "grinding" low-power battles between evenly matched players
- Limited to enemy hexes only; unclaimed hexes always taken instantly

### Exponential Upgrade Costs with Constant Rewards
- **Economy cost:** Cheap start (60g build + 60g L1→L2), doubles from there. L4 costs 480g.
- **Reward per level:** +0.9/sec net income gain (constant — each upgrade equally valuable)
- **Defense/Research cost:** Still use `BuildCost × 2^level` (L1=120/160g)
- **Effect:** Economy investment is accessible early game; exponential scaling caps extreme high-level chains. Defense requires deliberate gold commitment at every level.

---

## Open Questions / TODO

### Map Size & Pacing Analysis

**Recommended:** 60-80 total hexes per 1v1 map

**Rationale:**
- Players start with 1 hex each (2 total)
- ~30-40 unclaimed hexes available
- ~25-30 hexes per player after midgame (realistic distribution)

**Pacing with 70 hexes (example):**
- **T=0-2min:** Each player rapidly claims ~10 unclaimed hexes (10g each, income scales fast)
- **T=2-3min:** Land-grab slows, players meet at borders with ~10 hexes each
- **T=3-5min:** First border skirmishes, Economy and Defense buildings appear
- **T=5-15min:** Active border warfare, territory trades hands
- **T=15-25min:** One player pushes toward 60% (42 hexes) or secures Tech Level 3
- **T=25-30min:** Final race to victory condition

**Too small (30 hexes total):** Players meet at T=1min, constant warfare, no economy buildup, RNG-heavy
**Too large (120+ hexes):** Players farm 15+ minutes unopposed, snowball guaranteed, long game

### Playtesting Needs

- [ ] **Exponential cost feel:** Does progression curve feel right? Too fast/slow?
- [ ] **Counter-spend cap:** Is +3 power cap balanced? Create interesting battles?
- [ ] **Map size:** Does 70-hex map hit 30-min target? Adjust if needed.
- [ ] **Conquest threshold:** Does 60% + 10 sec create tense endgame?
- [ ] **Tech viability:** Do players build Research? Or pure military?

### Mechanical Unknowns

- [ ] **Player starting position:** Do all players start equidistant? Random corners? (Recommend: opposite corners for 1v1)
- [ ] **Building demolish refund:** Is 50% refund fair or should it be 100%?
- [ ] **Multiple battles:** Can hex be attacked by multiple enemies simultaneously? (Recommend: one attacker at a time, queue battles)
- [ ] **Choke point deadlock:** Narrow maps allow a single high-Power hex to block all expansion indefinitely. Map generation must avoid single-hex corridors, or a flanking/bypass mechanic is needed.

### Future Mechanics (Post-MVP)

- [ ] Alliances (2v2 mode with shared resources)
- [ ] Diplomacy (trade, temporary truces)
- [ ] Special hex types (mountains, water, resources)
- [ ] Hero units (unique units on hexes)
- [ ] Persistent progression (ranked ladder, cosmetics)

---

## Implementation Notes

### Scope
- **Start with 1v1** — Easier to balance before adding 3-4 player variants
- **UI priority:** Show hex Power prominently, battle timer clearly, resource flow transparent

### Tech Stack (Decided)

| Layer | Choice | Swap later to |
|---|---|---|
| Server | Go | — |
| Client | TypeScript | — |
| Rendering | HTML Canvas (hex grid) + DOM (UI overlays) | Svelte for complex UI |
| Bundler | Vite | — |
| Networking | WebSocket, JSON messages | MessagePack for 4-player |
| WebSocket lib | `github.com/coder/websocket` | — |
| Game loop | Goroutine + `time.Ticker` (100ms) | — |
| Deployment | localhost (MVP) | VPS or Fly.io |
| Database | None — in-memory state | SQLite for match history |
| Auth | None (MVP) | OAuth when accounts added |
| Testing | Go `testing` (headless) + manual browser | Playwright |

### Architecture Constraints
- **Server-authoritative:** Client never modifies game state — only sends intentions ("claim hex X", "attack hex Y"), server validates and responds
- **Tick-based economy:** Every 100ms tick: process queued player actions → update gold/TP → advance battle timers → check victory → emit delta to clients
- **Delta state sync:** Send only what changed per tick (hex ownership, resource updates, battle progress). Full state dump only on initial connect/reconnect.
- **Hex grid math:** Axial coordinates. Pixel ↔ hex conversion for click detection.
- **One goroutine per game room.** Player actions arrive via WebSocket, queued into a channel, processed at next tick boundary.
- **Message types (JSON):** Define shared message schema early — client and server must agree on shape. Keep message types in a shared doc or generate from a single source.

### Tick Execution Order (6 phases, strict sequence)

```
Phase 1: PROCESS ACTIONS — drain queued player intentions, validate, apply
Phase 2: ECONOMY — calculate income/maintenance, accrue gold/TP (×0.1 per tick)
Phase 3: BATTLES — decrement timers by 0.1s, resolve expired (transfer or retain hex)
Phase 4: AUTO-DROP — check negative income, manage 10s grace, drop if expired
Phase 5: VICTORY — conquest (60%/10s), tech dom (4 techs+35%/10s), capital loss, time limit
Phase 6: DELTA — diff vs previous tick, broadcast to clients
```

**Why:** Actions first = immediate effect. Economy before battles = counter-spend gold already deducted. Battles before victory = hex transfers count this tick. Victory last = reflects true final state.

### Edge Case Decisions

| Scenario | Decision |
|---|---|
| Disconnect mid-battle | Game continues 30s, then forfeit |
| Two attacks on same hex | First-in-queue wins, second rejected (one battle per hex) |
| Reconnection | Full snapshot (~3KB), no replay log |
| Capital captured during outgoing attack | Immediate game over, all battles canceled |
| Counter-spend same tick as battle expires | Applied (Phase 1 runs before Phase 3) |
| Gold below 0 | Clamped to 0, actions rejected if insufficient |
| Player trapped (no Power≥1 border) | Must build Defense to regain attack ability |

### Project Structure

```
hexar/
├── CLAUDE.md
├── go.mod
├── cmd/server/main.go           # entry point, HTTP+WS, serves client
├── internal/
│   ├── game/                    # PURE logic (zero I/O, zero network, deterministic)
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
│   └── mapgen/                  # Hardcoded test map (MVP)
│       └── generator.go
├── client/
│   ├── index.html
│   ├── vite.config.ts
│   ├── package.json
│   └── src/
│       ├── main.ts
│       ├── net/                 # WebSocket, reconnect, message types
│       ├── state/              # Game state mirror, apply delta/snapshot
│       ├── render/             # Canvas hex grid, buildings, battles
│       ├── input/              # Pixel→hex detection, action dispatch
│       ├── ui/                 # DOM: HUD, build menu, tech tree
│       ├── hexmath.ts          # Axial math (mirrors server)
│       └── constants.ts        # Mirrors server constants
└── .claude/skills/
```

**Key rule:** `internal/game/` has zero imports outside stdlib. `RunTick(state, dt)` is a pure function — fully testable with `go test` alone.

### Rejected Alternatives
- **Node.js server:** Go developer, worse concurrency model for tick loops
- **Phaser/PixiJS:** Overkill for colored hexagons + text; adds framework weight
- **React/Vue/Svelte (for MVP):** Adds ceremony for simple "draw hexagons, update numbers" UI
- **WebRTC:** Massive complexity for P2P; WebSocket latency (100ms ticks) is fine for Hexar
- **Serverless (Lambda):** Can't maintain WebSocket + tick loop — wrong model for real-time games
- **Player-hosted P2P:** Trust issues, NAT traversal, cheating risk

### First-Week Prototype Tasks (ordered by risk)
1. Go tick loop → WebSocket → browser Canvas hex grid (prove full pipeline: 100ms ticks render in real time)
2. Click-to-hex detection (axial math: click pixel → identify hex → send action to server)
3. Two-tab 1v1 (two browser tabs, same game, both see same state within 100ms)
4. Economy tick test — headless Go test asserting gold matches CLAUDE.md math after N seconds


