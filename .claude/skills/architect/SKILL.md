---
name: architect
description: "Use when: designing server/client split, game loop, state sync model, data model, or system boundaries for Hexar before coding starts"
type: skill
---

# Architect Skill — Hexar System Design

Designs the technical architecture for Hexar by asking the right questions, then deriving answers from CLAUDE.md mechanics and tech-stack decisions.

**Prerequisites:** Run `/tech-stack` first — architecture depends on chosen language, framework, and networking layer.

## Workflow

### 1. Load Constraints from CLAUDE.md (MANDATORY)
- Read CLAUDE.md fresh. Extract all timing-sensitive mechanics (tick intervals, battle durations, grace periods).
- Note what must be tamper-proof (gold, Power, victory checks).
- Identify real-time requirements (what players see updating live).

### 2. Determine Authority Model
Answer for each mechanic in CLAUDE.md:
- **Must be server-authoritative?** (Can cheating here ruin the game?)
- **Can be client-predicted?** (Would latency make this feel bad without prediction?)
- **Pure client-side?** (Only affects local display, no gameplay impact?)

### 3. Design Game Loop
- What happens each tick? (Derive from CLAUDE.md economy/combat/victory rules)
- What is event-driven vs tick-driven? (e.g., "claim hex" = event; "income accrual" = tick)
- What's the tick rate? (Trade-off: precision vs bandwidth vs CPU)

### 4. Design State Model
- What is the minimal canonical state the server must hold?
- What does each client need to render the game? (Subset of server state + local UI)
- What delta format minimizes bandwidth while keeping clients in sync?

### 5. Define Module Boundaries
For each module, define: responsibility, inputs, outputs, and what Hexar-specific concern makes it non-trivial.

Key questions:
- Where does hex adjacency validation live? (Used by expansion AND attack — shared)
- Where does battle state live? (Server owns timer; client renders countdown)
- How are simultaneous actions handled? (Two attacks on same hex? Race conditions?)
- How does reconnection work? (Full state dump or replay from last known?)

### 6. Document Decisions
- Write architecture to CLAUDE.md Implementation Notes (or a dedicated ARCHITECTURE.md)
- Record decisions as ADRs (Architecture Decision Records): decision + context + consequences

---

## Architecture Questions Checklist

### Server Authority
- [ ] Which game actions require server validation before being applied?
- [ ] What's the latency budget? (How stale can client state be before it feels wrong?)
- [ ] How is clock synchronization handled? (Battle timers must agree between players)

### Game Loop
- [ ] What tick rate? (Consider: economy updates, battle timer resolution, bandwidth)
- [ ] What operations run every tick vs only on events?
- [ ] What order do operations execute within a tick? (Economy before victory check? Matters.)

### State Sync
- [ ] Full state on connect, deltas after — or replay-based?
- [ ] What's in a delta message? (Minimum fields to reconstruct state change)
- [ ] How does the client handle out-of-order or missed deltas?

### Data Model
- [ ] What entities exist? (Hexes, players, buildings, battles, tech unlocks)
- [ ] What are the relationships? (Hex → owner, hex → building, player → techs)
- [ ] What indexes are needed for fast queries? (Hex neighbors, player hex list)

### Module Boundaries
- [ ] What shared logic do server and client both need? (Hex math, constants, validation)
- [ ] What's the testing seam? (Can you run game logic without rendering? Without network?)
- [ ] What's the deployment unit? (Monorepo? Separate server/client packages?)

### Edge Cases
- [ ] Player disconnects mid-battle — what happens?
- [ ] Two players attack same hex simultaneously — who wins?
- [ ] Player's capital is attacked while they're attacking elsewhere — priority?

---

## Output Format

```
## Architecture Decision: [Topic]
**Question:** [What needed deciding]
**Context:** [Relevant CLAUDE.md constraints and tech-stack choices]
**Decision:** [What was chosen]
**Consequences:** [What this enables and what it prevents]
**Alternative rejected:** [What else was considered and why not]

## Module: [Name]
**Responsibility:** [What it owns — one sentence]
**Inputs:** [What it receives]
**Outputs:** [What it emits]
**Hexar-specific concern:** [Why this is non-trivial for THIS game]
```

---

## Quick Reference (verify against CLAUDE.md — these may be outdated)

```
Timing-sensitive mechanics:
  Economy tick: every 100ms (server-dependent)
  Battle duration: 6-15 seconds (formula: 5 + (Attacker Power + Defender Power) / 2)
  Counter-spend window: real-time during battle countdown (capped at +3 or seconds remaining)
  Auto-drop grace period: 10 seconds (20s with Resilience tech)
  Instant takeover: when Power diff > 3 (enemy hexes only)

State that must be server-authoritative:
  Gold/TP balances, hex ownership, building levels, Power values,
  battle timers, capital hex location, tech unlocks, maintenance calculations

Events (not ticks):
  Claim hex, build/upgrade/demolish, initiate attack, unlock tech,
  counter-spend activation, auto-drop choice, capital capture

Entities:
  Hex: {id, owner, building_type, building_level, power, is_capital}
  Player: {id, gold, tp, techs_unlocked[], hex_count, capital_hex_id}
  Battle: {hex_id, attacker_id, timer_remaining, attacker_power, defender_power, counter_spend_used}
  Victory: Capital capture → immediate game over, all loser hexes become unclaimed
```

---

## Rules

- **Derive, don't assume.** Every architectural choice must trace back to a CLAUDE.md requirement or a tech-stack decision.
- **No premature optimization.** Design for 1v1 first. "But what about 100 concurrent games" is not an MVP question.
- **Testability first.** If you can't run the game loop in a unit test without a browser, the architecture is wrong.
- **Name the unknowns.** If a decision can't be made without prototyping, say so — don't guess.
- **Quick reference is orientation only.** If it conflicts with CLAUDE.md, CLAUDE.md wins.
