---
name: architect
description: "Use when: designing server/client split, game loop, state sync model, data model, or system boundaries for Hexar before coding starts"
type: skill
---

# Architect Skill — Hexar System Design

Designs the technical architecture for Hexar: server authority model, game loop tick, state synchronization, data model, and module boundaries. Getting this right before coding avoids expensive refactors.

## Workflow

### 1. Extract Technical Constraints
- Read CLAUDE.md for timing-sensitive mechanics (battle timers, economy ticks, auto-drop grace period)
- Identify what must be server-authoritative (economy, combat resolution, victory check)
- Identify what can be client-side (animations, UI state, predicted hex claims)

### 2. Design Game Loop
- Define tick rate and what happens each tick
- Map each CLAUDE.md mechanic to a tick operation
- Identify which operations are event-driven vs tick-driven

### 3. Design State Model
- Define the canonical game state (what the server owns)
- Define the client view (subset of server state + local UI state)
- Design delta format: what changes per tick, how it's serialized

### 4. Design Module Boundaries
- Server modules: game loop, economy engine, combat resolver, victory checker, connection manager
- Client modules: renderer, input handler, state reconciler, UI overlay
- Shared: hex math, constants, message types

### 5. Document Architecture
- Write decisions to CLAUDE.md Implementation Notes
- Produce a module diagram description
- Flag open questions for the tech-stack skill to resolve

---

## Key Architecture Questions for Hexar

### Server Authority
- Economy ticks run server-side (prevents cheating income)
- Battle timer runs server-side (client shows countdown, server resolves)
- Hex claims: server validates Power ≥ 1 and adjacency before confirming
- Auto-drop: server triggers, client shows 10-second UI prompt

### Tick Rate Design
```
Every 100ms tick:
  - Increment gold for each player (hex_count × rate)
  - Decrement maintenance costs
  - Check auto-drop condition
  - Progress active battle timers
  - Check victory conditions (60% hold timer, Tech Level 4 + 35%)
  - Emit delta to all clients
```

### Delta State Format
Only send what changed:
- Hex ownership changes (hex_id, new_owner)
- Resource updates (player_id, gold_delta, tp_delta)
- Battle state (hex_id, timer, attacker_power, defender_power)
- Victory progress (player_id, conquest_timer, tech_level)

### Data Model Sketch
```
GameState:
  hexes: Map<hex_id, { owner, building_type, building_level, power }>
  players: Map<player_id, { gold, tp, tech_unlocked[], hex_count }>
  battles: Map<hex_id, { attacker_id, timer, committed_gold }>
  victory: { conquest_timer, conquest_leader }
```

---

## Output Format

```
## Module: [Name]
**Responsibility:** [What it owns]
**Inputs:** [What it receives]
**Outputs:** [What it emits]
**Hexar-specific concern:** [Why this module is non-trivial for this game]
```

Final output: module list, tick loop pseudocode, state model schema, list of decisions that depend on tech-stack choice.

---

## Notes

- Server must be authoritative for all gold/TP transactions (no client-side economy)
- Battle counter-spend cap (min(+3, seconds_remaining)) must be enforced server-side
- Capital capture detection (game over) must be immediate, not deferred to next tick
- Hex adjacency validation reused for both expansion and attack — centralize early
