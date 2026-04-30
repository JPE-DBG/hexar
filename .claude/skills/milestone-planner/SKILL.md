---
name: milestone-planner
description: "Use when: breaking Hexar design into buildable milestones, defining MVP scope, sequencing features, or deciding what to cut from first build"
type: skill
---

# Milestone Planner Skill — Hexar Build Sequence

Breaks the full CLAUDE.md design into ordered, testable milestones. Each milestone is shippable, validates a core assumption, and unlocks the next. Prevents building features whose dependencies aren't ready.

## Workflow

### 1. Inventory All Features
- Read CLAUDE.md and list every mechanic as a discrete feature
- Tag each: core loop / balance system / win condition / UI / networking

### 2. Find Dependencies
- Which features require others to be working first?
- What can be stubbed (e.g., static map before map generation)?
- What must be tested early to catch design flaws (economy ticks, combat resolution)?

### 3. Define Milestones
- Group features into milestones of 1-2 weeks each
- Each milestone ends with something playable or demonstrably testable
- Flag features that are post-MVP (nice-to-have, not needed to validate core loop)

### 4. Write Acceptance Criteria
- For each milestone: what does "done" look like?
- Prefer functional tests over unit tests at this stage ("can two players complete a game?")

### 5. Identify Risk Items
- Which milestone is most likely to surface a design flaw?
- Where might CLAUDE.md need revision based on implementation reality?

---

## Hexar Feature Inventory (from CLAUDE.md)

### Core Loop (must be in MVP)
- Hex grid rendering with ownership colors
- Economy tick (gold income, maintenance, net calculation)
- Unclaimed hex claiming (10 gold, instant)
- Defense building + Power system
- Enemy hex attack (100 gold, battle timer, counter-spend)
- Capital capture = game over
- Conquest victory (60% + 10 sec)

### Secondary Systems (MVP or early Phase 2)
- Economy building (income boost)
- Research building + TP generation
- Tech tree (4 techs)
- Tech Dominance victory condition
- Garrison passive defense
- Auto-drop on negative income

### Polish / Phase 3
- Map generation (avoid choke points)
- Multiplayer networking
- UI overlays (battle timer, resource counters, victory progress)
- Time limit tiebreaker
- Building demolish / refund

---

## Suggested Milestone Structure

```
M1 — Local Hex Grid (no game logic)
  Render hex grid, click to select, show coordinates
  Done when: grid displays correctly at target map size

M2 — Economy Loop (single player)
  Gold ticks, claim unclaimed hexes, maintenance, auto-drop
  Done when: one player can expand to 20 hexes and income math matches CLAUDE.md

M3 — Combat (single player vs dummy)
  Defense building, Power, attack button, battle timer, counter-spend
  Done when: all battle resolution cases (instant takeover, standard, fail) work correctly

M4 — Full 1v1 Local
  Two players on same machine, all buildings, tech tree, both victory conditions
  Done when: a complete game can be played and won

M5 — Networking
  Server-authoritative tick, delta state sync, two browser tabs as two players
  Done when: M4 game works over localhost WebSocket

M6 — Playtest Ready
  Map generation, UI polish, 30-min session target validated
  Done when: external playtesters can complete a game without guidance
```

---

## Output Format

```
## Milestone N — [Name]
**Goal:** [What assumption this validates]
**Features included:** [List]
**Explicitly deferred:** [What's NOT in this milestone]
**Done when:** [Acceptance criteria]
**Design risk:** [What CLAUDE.md assumption might break here]
```

---

## Notes

- Local 1v1 before networking — design flaws are cheaper to fix without network complexity
- Economy loop is the highest-risk milestone (stepped maintenance math must feel right in practice)
- Tech tree can be stubbed as manual toggles in M4 before full Research building flow
- Map generation is a post-MVP concern — use a hardcoded test map for M1-M5
