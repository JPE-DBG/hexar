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
**Win Conditions:** Capital Capture (see Victory Conditions)

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
- **Capital hex:** Has innate Power 1 (no building needed). All other owned hexes start at Power 0 until a Power building is placed.
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
- If no choice made, the hex with the lowest income drops automatically; among income ties, the hex with the least building investment drops first
- On drop: player receives **50% refund of the building cost** on that hex (if any building present)
- Dropped hex becomes unclaimed instantly (enemy can grab it for 10 gold)
- Net income per hex +1/sec (before buildings/upgrades)

### Resources

- **Gold:** Main currency for buildings, upgrades, and attacks
- **Tech Points (TP):** Earned from Research buildings (0.2 TP/sec per Research level), spent on tech tree

**Tech bonus stacking:** Economy building formula: `(base + 0.6 × level) × 1.5`. Flat tech bonuses (e.g., Prosperity +1/sec) add to output; multiplier techs (e.g., Compound Growth ×1.25) multiply the building formula output.
```
hex_income (Economy, no techs)        = (2 + 0.6 × level) × 1.5       → L1: 3.9/sec, L5: 7.5/sec
hex_income (Economy + Prosperity)     = (2 + 0.6 × level) × 1.5 + 1.0 → L1: 4.9/sec
hex_income (Economy + Compound Growth)= (2 + 0.6 × level) × 1.5 × 1.25→ L1: 4.875/sec
hex_income (Economy + both)           = (2 + 0.6 × level) × 1.5 × 1.25 + 1.0 → L1: 5.875/sec
```

### Buildings (One per Hex)

Each hex can have **one building** of three types. There is no separate "build" action — upgrading an empty hex to L1 is the first upgrade step. Building can be demolished (refund 50%) and replaced.

#### Economy Building
- **Effect:** +50% resources/sec from that hex (stacks additively with tech bonuses)
- **Upgrade cost:** L1=60, L2=120, L3=240, L4=480 (formula: `BuildCost × 2^level`, where level is current level before upgrade)
- **Reward (Compounding):** +0.6 resources/sec per level, multiplied by the +50% bonus
- **Formula:** `(base + 0.6 × level) × 1.5` → Level 1: 3.9/sec, Level 5: (2 + 3.0) × 1.5 = 7.5/sec
- **Max level:** Unlimited (incremental)

#### Power Building
- **Effect:** +1 Power per level (Power = level, so L1=1, L2=2, etc.)
- **Upgrade cost:** L1=60, L2=120, L3=240, L4=480 (formula: `BuildCost × 2^level`)
- **Max level:** Unlimited (incremental)

#### Research Building
- **Effect:** +0.2 TP/sec per level (fuel for tech tree)
- **Upgrade cost:** L1=80, L2=160, L3=320, L4=640 (formula: `BuildCost × 2^level`, BuildCost=80)
- **Reward (Linear):** +0.2 TP/sec per level (constant gain)
- **Max level:** Unlimited (incremental)

---

## Combat System

### Attack Rules

- **Adjacent only:** Can only attack hexes touching your hex
- **Unclaimed hex cost:** 10 gold (instant takeover, no battle)
- **Enemy hex cost:** 100 gold (triggers battle)
- **Requirement vs enemy:** Attacker Power > Defender Power (strictly greater, ≤ fails with no cost)
- **Garrison threshold:** If defender has Garrison tech, the effective attack threshold is `Defender Power + Garrison bonus`. The UI blocks attacks that cannot win even at battle start — the hint shows "Need Pwr > N" where N includes the Garrison bonus. (Server validates only base power; the UI prevents committing 100g to a guaranteed-loss battle.)

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
- **Garrison interaction:** If Garrison tech is unlocked, adjacent owned hexes each add +1 Power automatically (cap: +2 from Garrison). Counter-spend and Garrison share a total cap of +3 — combined boost cannot exceed +3.
- **Siege Mastery interaction:** If the attacker has Siege Mastery, their battle timers are 40% shorter. The seconds-remaining cap compresses the defender's reaction window — a 4-second battle still allows up to +3, but requires an immediate response.
- Example: Battle at 7 seconds, you're losing 3 vs 5. Spend 150 gold over 3 sec → Power jumps to 6, you win
- Example: Only 1 second left in battle → max +1 Power boost (50 gold), even if you have gold to spare
- Example (Garrison): 2 adjacent owned hexes give +2 passive (Garrison cap reached). Counter-spend max is now +1 (total cap is +3, Garrison already contributes +2)
- Cap prevents defender from "buying" complete victory; time pressure makes decisions tense

---

## Tech Tree

Research buildings generate Tech Points at **0.2 TP/sec per Research level**. Spend TP to unlock any tech in any order — no prerequisites. **Total tree: 460 TP across 12 techs. No single game unlocks everything**, so every game has a distinct tech build.

| Tech | Cost | Effect | Archetype |
|------|------|--------|-----------|
| Blitz | 20 TP | Unclaimed hex claims cost 0g (was 10g) | Aggressor |
| Fortify | 20 TP | New action: spend 40g to prevent instant-takeover on one hex for 90 seconds | Defender |
| Prosperity | 25 TP | Each Economy (Gold) building generates +1/sec additional income | Builder |
| Reclamation | 25 TP | Recapturing a hex you previously owned costs 50g instead of 100g | Territorial |
| Vanguard | 30 TP | After capturing an enemy hex, attacks within 12 seconds cost 50g (timer refreshes on each capture) | Aggressor |
| Garrison | 30 TP | During battle, each adjacent owned hex adds +1 Power to defense (cap: +2 from Garrison; total defensive cap remains +3, shared with counter-spend) | Defender |
| Supply Lines | 40 TP | Maintenance costs reduced: 0.9/sec (hexes 1-10), 1.8/sec (hexes 11-20), 2.7/sec (hexes 21+) | Builder |
| War Chest | 30 TP | When you capture an enemy hex, recover 30g | Territorial |
| Resilience | 45 TP | Auto-drop grace period doubled (10s → 20s); building refund on drop increased to 70% | Defender |
| Iron Grip | 55 TP | All owned hexes permanently +1 Power | Aggressor |
| Compound Growth | 65 TP | Economy (Gold) buildings output ×1.25 (applied before Prosperity's flat bonus) | Builder |
| Siege Mastery | 75 TP | Your attack battle timers reduced by 40% (minimum 3s); when timer expires at equal Power, attacker wins | Aggressor |

**Research investment guide:**
- 1 Research L1 (0.2 TP/sec): ~180 TP in 15 min → ~5 techs from the cheap end
- 2 Research L1 (0.4 TP/sec): ~360 TP in 15 min → ~8 techs (solid mix of cheap and mid-tier)
- 3 Research L1 (0.6 TP/sec): ~540 TP in 15 min → ~11 techs (Research specialist — viable path to full tree by minute 13)

**Archetypes enabled by tech combinations:**
- **Blitz Aggressor:** Blitz → Vanguard → Iron Grip — free land-grab, chain attacks, full-territory Power
- **Economic Builder:** Prosperity → Supply Lines → Compound Growth — sustain wide territory, win by income weight
- **Fortress Defender:** Garrison → Fortify → Iron Grip → Resilience — make attacking you too expensive
- **Siege Striker:** Iron Grip → Siege Mastery → Vanguard — fast decisive battles with Power baseline across all hexes
- **Territorial:** Reclamation → Vanguard → War Chest — fluid borders, sustain attack chains through gold recovery

**Tech Stacking Rules:**
- **Attack Cost Discounts Stack:** Reclamation (-50g) and Vanguard (-50g) stack additively when both conditions are met
  - Base attack cost: 100g
  - With Reclamation only: 50g (recapturing your own hex)
  - With Vanguard only: 50g (within 12s of previous capture)
  - **With both:** 0g (free attack when reclaiming your own hex during Vanguard window)
  - Example: You capture hex A, lose it to enemy, immediately reclaim it within 12s → 0g cost
  - This rewards aggressive territorial play and creates high-value moments during Vanguard windows
- **Prosperity + Compound Growth Stack:** Applied in sequence (base × 1.25 multiplier, then +1.0 flat bonus)
  - Economy L1 with both techs: (2 + 0.6) × 1.5 × 1.25 + 1.0 = 5.875 gold/sec
- **Iron Grip + Garrison Stack:** Iron Grip adds +1 to all owned hexes; Garrison adds up to +2 during defense battles
  - Defender with both: base power + Iron Grip +1 + Garrison +2 (max) = +3 total possible bonus
  - **Display distinction:** Iron Grip is static (always-on) → included in hex label and power total. Garrison is dynamic (defense-only, depends on adjacency) → shown as separate `"+N def"` note in build menu, excluded from the static power number on the hex.

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

T=2.5 min: Save 60 gold → Build Power on border hex (~5 sec at +11/sec)
           Spend 30 → Upgrade Power L1. Power = 2 on that hex.

T=3 min:   First enemy contact at borders
           Save 100 gold for first enemy hex attack (~9 sec at +11/sec)
           Power 2 vs enemy Power 1 → Standard battle (6.5 sec), you win

T=5 min+:  Border warfare begins in earnest
           ~300-400 gold saved for upgrades and attacks
           Power 2-3 on borders, Economy building on 2-3 hexes
```

---

## Victory Conditions

### Capital Capture
- **Win by capturing the enemy's capital hex.**
- All of the loser's hexes become unclaimed instantly.
- No timers, no percentages — one decisive target per player.

**Why capital capture:**
- Genre standard for hex conquest games (Antiyoy, Polytopia)
- Forces the full economic loop: economy → Power → territory → capital approach
- Clear target for both players with rich counterplay (defend capital with Power buildings, Garrison, counter-spend)
- Anti-stalemate is built-in: maintenance auto-drop gradually weakens over-extended players, exposing their capital; tech (Siege Mastery, Iron Grip) breaks late defensive deadlocks

**How to win:**
- Expand territory toward the enemy capital corridor
- Build Power advantage on hexes adjacent to the capital approach
- Execute the final capital attack when Power difference is decisive (diff > 3 = instant takeover)

**How to stop opponent:**
- Maintain high-Power hexes on your capital's adjacent hexes
- Use Garrison, Fortify, and counter-spend to defend capital battles
- Launch counter-offensives to force opponent to retreat and defend their own capital

---

## Balance Rules & Constraints

### Upgrade Costs (Unified formula: `BuildCost × 2^currentLevel`)
- **Gold (BuildCost=60):** L1=60, L2=120, L3=240, L4=480
- **Power (BuildCost=60):** L1=60, L2=120, L3=240, L4=480 (same curve as Gold)
- **Research (BuildCost=80):** L1=80, L2=160, L3=320, L4=640
- There is no separate "build" action — upgrading empty hex to L1 costs `BuildCost × 2^0 = BuildCost`
- **Gold building rewards:** +0.6/sec per level × 1.5 multiplier = +0.9/sec net gain per level
- **Power/Research rewards:** +1 Power per level; +0.2 TP/sec per level (constant)
- **Demolish refund:** 50% of total invested; for all buildings: TotalInvested = `BuildCost × (2^level - 1)`, refund = `BuildCost × (2^level - 1) × 0.5`
  - Gold L1: TotalInvested=60, refund=30; L2: TotalInvested=180, refund=90; L3: TotalInvested=420, refund=210

### Maintenance System (Prevents Extreme Expansion)
- Each hex costs 1 maintenance/sec to hold
- Forces quality over quantity
- At ~10+ hexes, income caps out without production upgrades
- If maintenance exceeds income, slowest hexes auto-drop (player chooses order)

### Tech Research Investment (Enhances Conquest)
- Research buildings generate 0.2 TP/sec per level — doubled from initial design to make single-building investment meaningful
- With 2-3 Research buildings, players unlock 6-10 techs in a 20-min game
- Tech investment competes with Gold/Power building investment — no "free" tech path
- All 12 techs support conquest; no separate tech-based win condition

### Early Game Parity
- All players start with 1 capital hex, Power 1, no buildings
- First minute: rapid hex expansion (unclaimed hexes taken instantly)
- 5-10 minute mark: players meet at borders
- Game stays competitive until ~15-minute mark if balanced play

### Capital Hex Rules

- **If capital is captured:** Player immediately loses the game. All their hexes become unclaimed instantly.
- **Captured capital:** Becomes a normal hex for the conqueror (no innate Power 1, no special rules)
- **Building on capital:** Owner can place Power buildings on their capital (Power stacks on top of innate Power 1, e.g., Power L2 = Power 3 total)
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
- **Power/Research cost:** Still use `BuildCost × 2^level` (L1=120/160g)
- **Effect:** Economy investment is accessible early game; exponential scaling caps extreme high-level chains. Power building requires deliberate gold commitment at every level.

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
- **T=3-5min:** First border skirmishes, Economy and Power buildings appear
- **T=5-15min:** Active border warfare, territory trades hands
- **T=15-25min:** Players contest capital approach corridors; tech investments (Siege Mastery, Iron Grip) enable decisive attacks
- **T=25-30min:** Final push on the enemy capital — high-Power adjacent hex + 100g attack

**Too small (30 hexes total):** Players meet at T=1min, constant warfare, no economy buildup, RNG-heavy
**Too large (120+ hexes):** Players farm 15+ minutes unopposed, snowball guaranteed, long game

### Playtesting Needs

- [ ] **Exponential cost feel:** Does progression curve feel right? Too fast/slow?
- [ ] **Counter-spend cap:** Is +3 power cap balanced? Create interesting battles?
- [ ] **Map size:** Does 70-hex map hit 30-min target? Adjust if needed.
- [ ] **Tech build diversity:** Which techs do players prioritize? Do all archetypes (Aggressor, Builder, Defender, Territorial) appear in practice?
- [ ] **Vanguard stacking:** Timer refreshes on each capture, allowing chain attacks at 50g each. Does this create unstoppable snowball, or is it balanced by Power requirements?

### Mechanical Unknowns

- [x] **Player starting position:** Opposite corners for 1v1 (decided)
- [x] **Building demolish refund:** 50% (implemented and playtested)
- [x] **Multiple battles:** One attacker at a time — second attack on same hex rejected (implemented)
- [ ] **Choke point deadlock:** Narrow maps allow a single high-Power hex to block all expansion indefinitely. Map generation must avoid single-hex corridors, or a flanking/bypass mechanic is needed.

### Future Mechanics (Post-MVP)

- [ ] **Tech Tree Benefit Display:** Show quantitative benefits in tech tree UI (e.g., "Prosperity: +3.0 gold/s total (3 Economy buildings)", "Supply Lines: -2.0 gold/s maintenance (current: 15 hexes)"). Helps players evaluate tech value before unlocking. Deferred as non-critical UX enhancement — current static descriptions are sufficient for MVP.
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

**Power display conventions (implemented in M5):**
- Hex canvas labels show **effective combat power** (base + Iron Grip), not raw building level
  - Power building L1 with Iron Grip → label "P2" (effective), not "P1" (building level)
  - Capital + Power L1 + Iron Grip → "P3" (1 capital innate + 1 level + 1 IG)
  - Non-Power owned hexes (economy, research, empty) show a small yellow power badge if power > 0 (capital innate power or Iron Grip)
- **Garrison excluded from hex label** — it's a defense-battle-only bonus, shown as `"+N def"` note in build menu
- **Fortify timer** rendered as shrinking lime border segments (clockwise from top), 90s → 0s

**Timer visualization conventions:**
- **Hex timers** (timer state lives on a specific hex): rendered as shrinking colored border segments, clockwise from the top point, full ring at max duration → empty at 0. Fortify is the reference implementation. All future hex-level timers must follow this pattern.
  - Color must be readable on both player hex colors (teal P1, red P2); lime `#c8ff70` is the established choice for Fortify
- **Player timers** (timer state lives on a player, not tied to a hex): shown in the HUD as a text indicator
  - Example: Vanguard (12s) → `"⚡Vanguard X.Xs"` in HUD — it is a player timer, not a hex timer, so it does not use border segments
- **Battle timers**: separate established pattern — amber pulsing ring + countdown text + boost dots (not changed to border segments)

**Timer stacking (multiple timers on same hex):** Battle timer always renders at full brightness (urgent, action-required). Passive hex timers (Fortify) dim to 40% opacity and lineWidth 1.5 while a battle is active on that hex.

| Timer | Type | Visual | When stacked with battle |
|---|---|---|---|
| Fortify (90s) | Hex timer | Lime `#c8ff70` border segments, clockwise | Dimmed 40% opacity, lineWidth 1.5 |
| Battle duration (6-15s) | Battle timer | Amber pulsing ring + countdown text + boost dots | Always full brightness (priority) |
| Vanguard (12s) | Player timer | HUD text `⚡Vanguard X.Xs` | n/a |

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

### Session & Connection Behaviour (implemented post-M6)

**Lobby flow:**
- `POST /lobby/create` → room code (4-char hex) + session token (32-char hex) + playerID
- `POST /lobby/join` → token + playerID (room full after 2 players)
- `GET /ws?code=…&token=…` → WebSocket connection; server validates token before upgrading

**Session persistence (URL hash):**
- On game start, client writes `/#CODE:TOKEN` into the URL (`history.replaceState`)
- On page load, client checks sessionStorage first (network-drop reconnects), then URL hash (tab-close reconnects)
- URL hash is per-tab (unlike localStorage), so two players on the same browser get different hashes and don't collide
- Hash is cleared on game end (victory, disconnect, "Already Connected")

**Waiting state:**
- Room created → `state.Waiting = true`; game loop goroutine NOT started
- Second player connects → server sets `state.Waiting = false`, starts the loop, broadcasts to player 1
- Player 1 sees "Waiting for opponent…" overlay with room code displayed

**Pause on disconnect:**
- Any player disconnects during an active game → `state.Paused = true`, `state.PauseTimeLeft` set to that player's remaining budget
- Loop ticks continue (so clients receive countdown updates) but `RunTick` is skipped — economy, battles, auto-drop all frozen
- Remaining player sees purple banner with live `M:SS` countdown
- Reconnecting player → `state.Paused = false`, budget saved to `remainingGrace[pid]`; game resumes from exact state
- Budget exhausted → loop clears pause, enqueues `ActionForfeit`; next tick processes it normally

**Forfeit budget:**
- Each player starts with 120 seconds total across the whole game
- Reconnecting saves the remaining time; a second disconnect resumes from where it left off, not from 120s again

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
Phase 5: VICTORY — capital capture (all loser hexes unclaimed instantly)
Phase 6: DELTA — diff vs previous tick, broadcast to clients
```

**Why:** Actions first = immediate effect. Economy before battles = counter-spend gold already deducted. Battles before victory = hex transfers count this tick. Victory last = reflects true final state.

### Edge Case Decisions

| Scenario | Decision |
|---|---|
| Disconnect mid-game | Game pauses immediately; opponent sees live countdown; player has 2-min cumulative budget to reconnect (not reset per disconnect); budget exhausted → forfeit |
| Two attacks on same hex | First-in-queue wins, second rejected (one battle per hex) |
| Reconnection | Full snapshot on connect (prevSnapshot = nil); delta sync resumes after that |
| Capital captured during outgoing attack | Immediate game over, all battles canceled |
| Counter-spend same tick as battle expires | Applied (Phase 1 runs before Phase 3) |
| Gold below 0 | Clamped to 0, actions rejected if insufficient |
| Player trapped (no Power≥1 border) | Must build Power to regain attack ability |
| Same player opens duplicate tab | Server closes new WS with code 4001; client shows "Already Connected" overlay, no retry |
| Both players not yet connected | Game loop does not start; `state.Waiting = true`; loop starts only when all player slots filled for first time |
| Game over by forfeit vs capital | `state.WinReason`: `"forfeit"` or `"capital"` — victory screen shows distinct subtitle |
| Demolish while hex selected | Hex is deselected after demolish completes — sidebar clears context so stale upgrade/attack buttons don't remain |
| Capture flash trigger | Compare incoming delta `hex.owner` against current `state.hexes` owner — **never** use `hex.previousOwner` from the wire. Server sets `PreviousOwner` on capture and never resets it; any future delta for that hex (e.g. fortifyTimer tick) would re-fire the flash spuriously every second |

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

### UI Architecture (M7 Design)

**Pattern:** Sidebar with sticky tools (left edge, 180px wide)

**Chosen over alternatives:**
- ~~Bottom menu~~ — 40% slower for batch building (10 actions vs 6 for placing 5 Economy buildings)
- ~~Radial menu~~ — occludes adjacent hexes during battles when players need to see Garrison bonuses
- ~~Mode-based~~ — creates "what mode am I in?" confusion; sidebar uses simple highlighted button state instead

**Sidebar sections:**
1. **CONTEXT** (dynamic, top): Upgrade/Attack/Demolish/Sell Hex/Fortify/Counter-spend buttons appear based on selected hex state
2. **BUILD** (always visible, bottom): Economy [Q], Power [W], Research [E]

**Sticky tool behavior:**
- Click tool button (or press hotkey) → button highlights → click hexes to place buildings
- Tool stays selected for batch operations (e.g., click Economy once, then click 5 hexes = 6 actions total)
- **Toggle to deselect:** Click selected button again (or press hotkey again, or press ESC) → unhighlights
- **Battle interaction:** BUILD tool is skipped when clicking a hex under active battle — click falls through to context selection showing counter-spend

**Keyboard shortcuts:**
- Q/W/E: Select Economy/Power/Research (toggle if already selected)
- D: Demolish selected hex (context action — fires on selected hex, not a sticky tool)
- X: Sell Hex selected hex (context action — fires on selected hex, not a sticky tool)
- Space: Upgrade selected hex (context action)
- A: Attack selected hex (context action)
- F: Fortify selected hex (context action, if tech unlocked)
- C: Counter-spend during battle (context action)
- ESC: Deselect any selected tool

**Battle restrictions:**
- No new building placement or upgrade allowed on hexes under active battle (server enforces `ErrBattleInProgress` in `ValidateUpgrade`)
- Demolish remains allowed during battle — defender may recover gold for counter-spend
- Upgrade button is hidden in context UI during battle; reappears automatically when battle resolves

**Why this is faster than bottom menu:**
- **Batch building (expansion phase):** 6 actions for 5 buildings vs. 10 actions with bottom menu
- **Mixed operations (warfare phase):** Equivalent speed (context actions work same way)
- **Always-visible affordances:** No hunting for buttons after clicking hex

**Responsive behavior:**
- Desktop (>1024px): Left sidebar, full labels, context area
- Mobile landscape (896×414): Bottom bar (60px), icon + label
- Mobile portrait (<768px): Bottom bar (50px), icons only

**Implementation files:**
- `client/src/ui/sidebar.ts` — replaces `buildmenu.ts`
- `client/src/style.css` — sidebar styles + responsive layouts
- `client/src/main.ts` — keyboard shortcuts and tool logic

### Milestone Status

**M1–M5:** Core game loop, hex grid, economy, combat, tech tree, UI polish, power display conventions.

**M6 (in-scope — implemented):**
- Simple lobby: `POST /lobby/create` and `POST /lobby/join` with room codes and session tokens
- Delta sync: incremental state updates (~300–500 bytes/tick vs 7KB full snapshot)
- WebSocket reconnect with exponential backoff (5 attempts, 1–16s)
- Disconnect forfeit via `ActionForfeit` enqueued through the action channel (not direct state mutation)

**Post-M6 additions (discovered during manual testing):**

| Feature | Trigger |
|---|---|
| Vite dev proxy for `/lobby/*` | "Create Game" always errored in dev mode — Vite had no proxy for the new HTTP endpoints |
| Duplicate tab prevention (WS close code 4001 + "Already Connected" overlay) | Opening the game in a duplicate tab silently evicted the original connection |
| URL hash session persistence (`/#CODE:TOKEN`) | Closing a tab lost sessionStorage; localStorage would collide when testing 2 players in same browser |
| Waiting state (`state.Waiting`) before both players connect | Game loop ran from room creation, giving player 1 a gold head-start before player 2 joined |
| Game pause on disconnect (`state.Paused`) | Active player had a free 2-minute window to expand unopposed while opponent was disconnected |
| Live forfeit countdown (`state.PauseTimeLeft`, `M:SS` banner) | Remaining player had no visibility into how long until forfeit |
| Cumulative reconnect budget (120s total, not reset per disconnect) | Per-disconnect timer allowed repeated short disconnects to accumulate unlimited free pause time |
| Win reason (`state.WinReason`: `"capital"` / `"forfeit"`) | Victory screen always showed "Capital captured" even when opponent forfeited |

**M7: Visual Polish**

Goal: Make the game look modern and appealing to retain real players. Every visual addition must reinforce player decision-making — decoration gets cut.

**Rendering approach (decided):** Stay on Canvas. rAF loop already exists in `renderer.ts`. No measured perf problem at 70 hexes. PixiJS migration only if rAF loop is measured > 12ms with effects enabled. No Svelte migration for the build menu — defer until live gold-threshold reactivity is actually needed.

**Effect priority (by decision-making value):**

| Effect | Value | Rationale |
|---|---|---|
| Capture flash on hex takeover | Must have | Confirms action resolved — player needs this to know attack succeeded |
| Battle ring pulse (enhance existing) | Must have | Urgency signal during live counter-spend window |
| **Sidebar with sticky tools** | Must have | Replace bottom menu; 40% faster batch building; always-visible tools with Q/W/E/D/X hotkeys and toggle-to-deselect |
| Lobby/overlay polish | Must have | First impression; affects player trust and willingness to return |
| Floating `+gold` on economy tick | Should have | Reinforces resource loop; helps players time purchases — keep small, float away from hex centre to avoid obscuring labels |
| **Sidebar CSS + enhanced SVG icons** | Should have | Color-coded icons (gold=#f9ca24, power=#e74c3c, research=#45b7d1); context actions; responsive layouts |
| HUD monospace numbers | Should have | Prevents layout shift on fast-updating gold/TP counters |
| Animated gold counter (smooth interpolation) | Nice to have | Low decision impact; add only after higher-priority items ship |
| Screen shake on capital fall | Nice to have | Pure feel — cut if time is tight |
| Idle hex shimmer | Skip | No decision value |

**Sequencing within M7:**
1. **Palette consolidation** — `constants.ts` as single source of truth for all colors; matching CSS file for DOM overlays; remove all inline `style.cssText` strings; no CSS variables that can drift from TS constants
2. **Hex renderer polish** — ~~radial gradient fill~~ **flat fill kept intentionally** (gradient tried and reverted; hex fill will be replaced with sprite-based 3D/material look in a future pass — `drawHexFill` is intentionally minimal until then); capital hex distinct icon; labels drawn in a separate Pass 3 (after rings and battles) so they render on top of fortify segments
3. **Effects layer** — `Effect[]` array in existing rAF loop; `addCaptureFlash(hex, color)` and `addFloater(hex, text)` called from `applyDelta`; capture flash compares incoming `hex.owner` against current `state.hexes` owner (not `previousOwner` from wire); effects are client-side only, never in game state
4. **Sidebar with sticky tools** — New `sidebar.ts` replaces `buildmenu.ts`; left-edge 180px; 3 sections (BUILD/MANAGE/CONTEXT); toggle-to-deselect behavior; Q/W/E/D/X/Space/A/F/C/ESC hotkeys; responsive layouts for mobile
5. **HUD polish** — Structured layout with monospace `tabular-nums` font; color-coded rates (green=positive, red=negative); visual hierarchy with labels/values/rates
6. **Overlay + lobby screens** — CSS for pause/victory/waiting/lobby; room code large and copyable on waiting screen; styled victory/forfeit distinction

**Technical constraints:**
- Canvas + DOM split preserved throughout
- Effects layer has zero game state access — purely visual, triggered by state transitions
- Render hierarchy must not be obscured: ownership → power → building → status

**M8: Testing**

Goal: Build confidence before exposing to real users; catch regressions as the visual layer grows. Start with zero mocks in `internal/game/` — it's pure, just call the functions.

**Layer 1: Game logic** (`internal/game/`) — stdlib only, no goroutines, table-driven:

| File | Test cases |
|---|---|
| `testhelpers_test.go` | `twoPlayerState()`, `runTicks(n, actions...)`, adjacent hex builder — shared by all game tests |
| `economy_test.go` | Gold accrues at 2/s per hex; stepped maintenance triggers at 10/20 hex boundaries; Prosperity + Compound Growth formula matches CLAUDE.md math |
| `combat_test.go` | Attack validation: insufficient power, insufficient gold, hex already in battle; battle resolution: winner at timer expiry, loser retains on tie (unless Siege Mastery); instant takeover when diff > 3 |
| `victory_test.go` | Capital capture sets `WinReason = "capital"`, all loser hexes unclaimed; `ForfeitPlayer` sets `WinReason = "forfeit"` |
| `tech_test.go` | Tech unlock deducts TP and sets bool; Reclamation + Vanguard both active + conditions met = 0g attack cost |
| `autodrop_test.go` | Negative net income triggers 10s grace; forced drop when grace expires; 50% refund on dropped hex |

**Layer 2: Room/lobby integration** (`internal/room/`) — real goroutines, minimal `time.Sleep` (1–2 tick durations):

| Test | Asserts |
|---|---|
| `TestWaitingState` | Loop does not start until 2nd player connects; P1 gold unchanged while waiting |
| `TestPauseOnDisconnect` | `state.Paused = true` immediately; remaining client receives snapshot |
| `TestUnpauseOnReconnect` | `state.Paused = false`; remaining grace saved to `remainingGrace[pid]` |
| `TestCumulativeGrace` | Second disconnect starts from remaining budget, not 120s |
| `TestForfeitEnqueued` | `PauseTimeLeft` hits 0 → `ActionForfeit` processed next tick |
| `TestDuplicateConnect` | `IsConnected` returns true; second connection replaces first cleanly |

**Layer 3: Playwright E2E** (`e2e/`) — runs locally via `make test-e2e`; no deployment needed. Playwright's `webServer` config auto-starts both the Go server (`:8080`) and the Vite dev server (`:5173`) before the suite runs. Multi-player scenarios use two `BrowserContext` objects in one test process.

| Spec | Tests |
|---|---|
| `lobby.spec.ts` | Create shows 4-char code; join with bad code shows error; waiting overlay appears; HUD activates when P2 joins |
| `gameplay.spec.ts` | Canvas + HUD visible after both connect; gold counter increases over time; sidebar visible |
| `connection.spec.ts` | Duplicate tab shows "Already Connected"; URL hash reconnects after tab close; pause banner appears on opponent disconnect |

**Setup:**
```bash
npm install          # installs @playwright/test from root package.json
npx playwright install chromium
make test-e2e        # starts both servers, runs e2e/, exits
```

**M9: Deployment**

Goal: Make the game accessible to real players outside localhost.

**Files to create (in implementation order):**

1. `internal/net/server.go` — add `GET /health` endpoint returning `{"status":"ok","version":"..."}` (Fly.io uses this for crash detection)
2. `client/src/net/connection.ts` — WebSocket URL: `const protocol = import.meta.env.PROD ? 'wss' : 'ws'`
3. `Dockerfile` — multi-stage: `node:20-alpine` (Vite build) → `golang:1.23-alpine` (server build, `CGO_ENABLED=0 GOOS=linux`) → `gcr.io/distroless/static-debian12` (runtime); distroless has no shell — minimal attack surface
4. `fly.toml` — `auto_stop_machines = false` (CRITICAL: Fly's default sleep kills active WebSocket connections), `min_machines_running = 1`, `force_https = true` (auto TLS upgrades `ws://` → `wss://`)
5. `.github/workflows/deploy.yml` — `test` job (`go test ./...`) + `deploy` job (`needs: test`, `flyctl deploy --remote-only`); deploy never runs if tests fail

**Deployment steps:**
1. Add `/health` endpoint + verify locally (`curl localhost:8080/health`)
2. Update WebSocket URL for `wss://` + verify Vite prod build connects correctly
3. Write `Dockerfile` → `docker build -t hexar .` locally → `docker run -p 8080:8080 hexar` smoke test
4. `fly launch` (generates initial `fly.toml`) → edit for Hexar constraints (`auto_stop_machines = false`)
5. Write GitHub Actions workflow; add `FLY_API_TOKEN` to GitHub secrets
6. `fly deploy` → post-deploy: load lobby, create game, join from second browser, verify `wss://` in DevTools Network tab

**Known limitation (document in lobby UI):** Rooms are in-memory — a server restart ends all active games. Future fix: SQLite via `modernc.org/sqlite` (pure Go, no CGO) persisting room state as JSON blob. Not needed for initial deployment.

**Post-M9 fixes (discovered after deployment):**

| Fix | Trigger |
|---|---|
| Mobile button layout — destructive actions (Sell Hex, Demolish) moved to top-left corner | On mobile, context action buttons at the bottom bar were obscured or mis-tapped; destructive icons repositioned away from the main build bar |
| Deselect on demolish | After demolishing a building the hex remained selected, leaving stale context buttons (upgrade, attack) visible |
| Fortify timer flickering fix | Hex would briefly fill with player color every 1s when fortified after capture — root cause: `previousOwner` is set on capture and never reset by server, so any delta for that hex (fortifyTimer tick) re-fired `addCaptureFlash`; fixed by comparing `hex.owner` vs current `state.hexes` owner in `onDelta`, not against `previousOwner` from wire |
| Delta quantization for FortifyTimer | `fortifyTimer` changed every 100ms tick, causing fortified hexes to appear in every delta; now only sent when timer crosses a 1-second boundary (`int(prev) != int(curr)`) |

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


