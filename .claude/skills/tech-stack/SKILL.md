---
name: tech-stack
description: "Use when: evaluating a proposed tech change to the existing Hexar stack, assessing a new library or tool, or deciding if a technology swap is worth the cost"
type: skill
---

# Tech Stack Skill — Hexar Technology Evaluation

Evaluates proposed changes to Hexar's technology stack against specific requirements. The stack is already decided and deployed — this skill is for assessing whether a change is worth the cost, not for choosing from scratch.

---

## Current Stack (Implemented)

| Layer | Choice | Potential future swap |
|---|---|---|
| Server | Go | — |
| Client | TypeScript | — |
| Rendering | HTML Canvas + DOM overlays | Sprite-based hex rendering (future pass) |
| Bundler | Vite | — |
| Networking | WebSocket, JSON | MessagePack for 4-player bandwidth |
| WebSocket lib | `github.com/coder/websocket` | — |
| Game loop | Goroutine + `time.Ticker` (100ms) | — |
| Deployment | Fly.io (distroless Docker) | — |
| Database | None (in-memory rooms) | SQLite via `modernc.org/sqlite` for persistence |
| Auth | None | OAuth when accounts needed |
| Testing | Go stdlib + Playwright | — |

See CLAUDE.md "Server Architecture → Tech Stack" for the full rationale behind each decision.

---

## Workflow

### 1. Load the Proposed Change

Clearly state:
- What technology is being added or swapped
- What problem it solves that the current stack does not
- What layer it touches (server, client, rendering, networking, deployment, testing)

### 2. Check Against Existing Constraints

Hexar-specific non-negotiables:
- Server must maintain WebSocket connections — no serverless/scale-to-zero
- Game loop must be a single goroutine per room with 100ms tick
- `internal/game/` must remain zero external imports (pure functions, stdlib only)
- Client rendering: Canvas for hex grid, DOM for overlays — any change must preserve this split
- Delta packets target < 500 bytes/tick at idle

### 3. Evaluate the Trade-off

For the proposed change:
- **What problem does it solve?** Is the problem real (measured) or hypothetical?
- **What does it break?** Existing tests, architecture invariants, deployment model?
- **What's the migration cost?** How many files change? Which systems does it touch?
- **Is it reversible?** Can you swap back if it doesn't work out?
- **Does it conflict with planned future work?** (e.g., sprite rendering, SQLite, 4-player)

### 4. Recommend

Choose one:
- **Accept:** Change is worth the cost. Document the new choice in CLAUDE.md and update "Rejected Alternatives" if relevant.
- **Defer:** Change is valid but not urgent. Add to CLAUDE.md "Future Mechanics" with rationale for when to revisit.
- **Reject:** Change solves a hypothetical problem, or cost exceeds benefit. Document rejection reason.

---

## Evaluation Questions by Layer

**Rendering change (Canvas → PixiJS / Sprite):**
- Is the rAF loop measured > 12ms? (Don't migrate without measured perf problem)
- Does the new approach preserve Canvas/DOM split?
- `drawHexFill` is intentionally minimal — sprite-based hex fill is planned; see CLAUDE.md

**Networking change (JSON → MessagePack):**
- Are delta packets actually hitting the 500-byte limit frequently?
- Does it change the client merge pattern in `applyDelta`?
- Worth doing when 4-player support needs it; not before

**UI framework (DOM → Svelte/React):**
- Is live reactivity actually needed? (Gold threshold buttons, tech tree)
- Does it break the Canvas/DOM split invariant?
- Deferred until reactivity is genuinely needed — DOM is currently sufficient

**Database (in-memory → SQLite):**
- Needed when room persistence across server restarts matters
- `modernc.org/sqlite` (pure Go, no CGO) is the planned choice — no Cgo, works in distroless

**Auth (none → OAuth):**
- Only needed when accounts/matchmaking/leaderboards are being built
- Don't add OAuth before there's something to protect

---

## Output Format

```
## Tech Change Evaluation: [Proposed Change]
**Layer:** [Which stack layer this affects]
**Problem being solved:** [Specific measured issue or anticipated need]
**Current approach limitation:** [Why the current solution falls short]
**Proposed solution:** [What changes]
**Conflicts with:** [Existing architecture, invariants, or planned work]
**Migration cost:** [Files affected, breaking changes]
**Reversibility:** [Can you undo this easily?]
**Recommendation:** Accept / Defer / Reject
**Rationale:** [Why]
```

---

## Rules

- **Hexar-specific, not generic.** Don't list general pros/cons. Explain the impact on THIS game's architecture.
- **Measure before migrating.** Don't swap rendering or networking without a measured problem.
- **One layer at a time.** Bundled tech swaps compound risk. Evaluate and migrate one layer before touching the next.
- **CLAUDE.md is authoritative.** Check "Rejected Alternatives" before proposing something already considered.
