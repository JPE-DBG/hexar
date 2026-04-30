---
name: tech-stack
description: "Use when: choosing language, game framework, networking library, or rendering approach for Hexar before coding starts"
type: skill
---

# Tech Stack Skill — Hexar Technology Decisions

Evaluates technology choices against Hexar's specific requirements: real-time 1v1, tick-based economy, hex grid rendering, browser or desktop target, 30-minute sessions.

## Workflow

### 1. Read Requirements
- Load CLAUDE.md to extract hard constraints (tick rate, real-time, delta state networking, hex grid)
- Identify non-negotiables vs nice-to-haves
- Note implementation complexity hints (axial coordinates, server authority, battle timers)

### 2. Evaluate Options
For each layer (language, framework, networking, rendering):
- List 2-3 realistic options
- Score against Hexar-specific constraints
- Flag where a choice locks in other choices downstream

### 3. Identify Integration Risks
- Which combinations have known friction (e.g., language X + library Y)?
- Where does the stack become complex at multiplayer boundary?
- What would be hardest to change after 2 weeks of coding?

### 4. Recommend Stack
- State a concrete recommended stack with rationale
- List explicit trade-offs accepted
- Note what to prototype first to validate the choice

### 5. Document Decision
- Write decision to CLAUDE.md Implementation Notes section
- Record the alternatives considered and why they were rejected

---

## Key Decisions for Hexar

### Rendering
- **Browser (Canvas/WebGL):** Phaser 3, PixiJS, or raw Canvas API
- **Desktop:** Godot, Unity (overkill for hex grid), custom
- **Constraint:** Hex grid with Power labels, battle timers, resource counters — not graphically demanding

### Language / Runtime
- **TypeScript + Node.js:** Single language full-stack, fast iteration, large ecosystem
- **Rust (Bevy/custom):** Performance headroom, but slower iteration, smaller team fit
- **Go (server) + TS (client):** Strong concurrency model for game tick loop, familiar web client

### Networking
- **WebSocket:** Simple, browser-native, adequate for 1v1 hex game tick rates
- **WebRTC DataChannel:** Lower latency, P2P option, higher complexity
- **Constraint:** Delta state only (CLAUDE.md note) — protocol must support partial updates efficiently

### Game Loop Architecture
- **Server-authoritative tick:** Server runs economy ticks (every 100ms), sends deltas to clients
- **Client-side prediction:** Predicts local state, reconciles on server tick — adds complexity
- **Constraint:** Battle timers must be server-side to prevent cheating

---

## Output Format

```
## Stack Decision: [Layer]
**Chosen:** [Technology]
**Rejected alternatives:** [A (reason), B (reason)]
**Why:** [Specific fit for Hexar constraints]
**Risk:** [What could go wrong with this choice]
**Validate by:** [What to prototype in first week to confirm]
```

Final output: a complete recommended stack table + first prototype checklist.

---

## Notes

- Evaluate for 1v1 MVP first; 2v4 player scaling is secondary
- Browser-first lowers barrier for playtesting (no install)
- Avoid engines that abstract away networking (makes delta-state hard)
- Hex grid math (axial coordinates) works in any language — not a differentiator
