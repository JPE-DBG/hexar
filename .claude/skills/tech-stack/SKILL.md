---
name: tech-stack
description: "Use when: choosing language, game framework, networking library, or rendering approach for Hexar before coding starts"
type: skill
---

# Tech Stack Skill — Hexar Technology Decisions

Evaluates technology choices against Hexar's specific requirements. Decisions made here feed into `/architect` and `/milestone-planner`.

**Prerequisites:** None. This is the first skill to run.

## Workflow

### 1. Load Requirements from CLAUDE.md (MANDATORY)
- Read CLAUDE.md to extract hard constraints (real-time ticks, delta state sync, hex grid, battle timers, 1v1 multiplayer).
- Separate non-negotiables from nice-to-haves.
- Note the CLAUDE.md "Implementation Notes" section for any existing decisions.

### 2. Evaluate Options Per Layer
For each technology layer:
- List 2-3 realistic options (not exhaustive surveys)
- Score against Hexar-specific constraints (not general "pros/cons")
- Flag where a choice locks in downstream decisions (e.g., "choosing X means you must use Y for networking")

### 3. Identify Integration Risks
- Which combinations have known friction?
- Where does complexity spike at the multiplayer boundary?
- What would be hardest to change after 2 weeks of coding?

### 4. Recommend Stack
- State concrete recommendations with rationale tied to CLAUDE.md requirements.
- List trade-offs accepted.
- Define a "first-week prototype" task that validates the riskiest choice.

### 5. Document Decisions
- Write to CLAUDE.md Implementation Notes section.
- Record rejected alternatives and why.

---

## Decision Areas

### 1. Language / Runtime
Questions to answer:
- Full-stack single language (faster iteration) vs specialized per layer (better fit)?
- What's the developer's existing expertise? (Don't fight the team)
- Server concurrency needs? (Tick loop + multiple connections)

### 2. Rendering / Client Framework
Questions to answer:
- Browser-first (no install, easy playtesting) or native (better performance)?
- 2D hex grid: how graphically demanding is it? (Hexar = simple: colored hexes + text labels)
- Does the framework fight you on custom real-time UI (battle timers, resource counters)?

### 3. Networking
Questions to answer:
- WebSocket (simple, browser-native) vs WebRTC (lower latency, P2P)?
- How much bandwidth does delta-state at chosen tick rate require?
- Does the solution handle reconnection gracefully?
- Does it support the "server-authoritative with client display" model?

### 4. Game Loop Host
Questions to answer:
- Where does the authoritative game loop run? (Dedicated server process? Serverless? One player hosts?)
- What's the latency budget for Hexar? (Economy ticks = tolerant; battle counter-spend = sensitive)
- How many concurrent games must one server instance support for MVP?

### 5. Deployment Target
Questions to answer:
- Where is the server hosted? (VPS, container service, serverless, edge?)
- What's the cost model? (Per-game? Always-on? Scale-to-zero?)
- What regions matter for latency? (1v1 = both players close to same server)
- How do players find and connect to a game? (Lobby? Matchmaking? Direct link?)

### 6. Persistence / Database
Questions to answer:
- What needs to survive a game session? (Player accounts? Match history? Leaderboards? Nothing for MVP?)
- In-game state: does it need persistence? (If server crashes mid-game, is that game lost? Probably yes for MVP)
- What's the simplest option for MVP? (SQLite? Postgres? No DB at all?)

### 7. Authentication
Questions to answer:
- Do players need accounts for MVP? (Guest/anonymous may be fine for playtesting)
- If yes: OAuth (Google/GitHub) or custom auth?
- How does player identity connect to matchmaking and leaderboards?

### 8. Testing Strategy
Questions to answer:
- How do you test the game loop without a browser? (Headless simulation)
- How do you test two-player interactions? (Scripted bots? Playwright?)
- What's the minimum viable test suite before first playtest?
- Can balance be validated programmatically? (Simulation runs → assertions on outcomes)

---

## Output Format

```
## Stack Decision: [Layer]
**Chosen:** [Technology]
**Rejected:** [A (reason), B (reason)]
**Why for Hexar:** [Specific fit against CLAUDE.md constraints]
**Risk:** [What could go wrong]
**Validate by:** [What to prototype in first week to confirm this choice works]
**Locks in:** [What downstream decisions this forces]
```

Final deliverable: complete stack table + ordered list of "validate first" prototype tasks.

---

## Rules

- **Hexar-specific, not generic.** Don't list general pros/cons of React vs Vue. Explain why a choice fits or fights Hexar's specific requirements.
- **Respect existing expertise.** If the developer already knows TypeScript, don't recommend Rust unless there's a compelling Hexar-specific reason.
- **MVP-scoped.** Choose for 1v1 with 2 concurrent players first. "But what about 1000 games" is post-MVP.
- **Reversibility matters.** Prefer choices that can be swapped later over ones that lock everything in.
