---
name: architect
description: "Use when: understanding the existing Hexar architecture, evaluating a proposed structural change, adding a new system, or debugging a cross-layer issue"
type: skill
---

# Architect Skill — Hexar System Reference

Reference guide for Hexar's implemented architecture. Use when evaluating structural changes, adding new systems, or understanding how components connect.

**The architecture is implemented.** This skill is not for designing from scratch — it's for navigating and changing what exists.

---

## Architecture Overview

**Pattern:** Server-authoritative, tick-based, delta state sync

```
Client (TypeScript)          Server (Go)
─────────────────           ─────────────
Canvas hex grid    ←─────── delta/snapshot (WebSocket JSON)
DOM overlays       ─────────→ actions (WebSocket JSON)
applyDelta/Snapshot          RunTick(state, dt) — 100ms
state.ts (mirror)            internal/game/ (pure functions)
```

**Authority model:**
- Server owns all game state — gold, Power, hex ownership, battle timers, tech unlocks
- Client never modifies game state — only sends intentions (actions)
- Client maintains a mirror of server state via snapshot+delta sync
- Effects (capture flash, gold floaters) are client-side only, never in game state

---

## Key Invariants

These must be preserved across any change:

1. **`internal/game/` zero external imports.** `RunTick(state, dt)` is a pure function. No I/O, no goroutines, no network. Fully testable with `go test` alone.

2. **One goroutine per room.** Actions arrive via WebSocket, queued into a channel, processed at the next tick boundary.

3. **Delta packets < 500 bytes at idle.** Tech field diffed (only send when changed); battles omitempty; FortifyTimer quantized to 1-second boundaries.

4. **Reverting PlayerDTO fields always sent.** `AutoDropActive`, `AutoDropGrace`, `VanguardTimer` can go back to zero/false — they must never be omitempty or the client's delta merge spread will retain stale values.

5. **Canvas + DOM split preserved.** Hex grid on Canvas, UI overlays as DOM. backdrop-filter on `#sidebar` makes `position: fixed` children relative to it — `#destructive-panel` lives outside `#sidebar` to avoid this.

6. **Capture flash compares against client state, not wire data.** `previousOwner` is set on capture and never reset by server — any future delta for that hex would re-fire the flash spuriously if you compare against it.

---

## Tick Execution Order (6 phases, strict sequence)

```
Phase 1: PROCESS ACTIONS — drain queued player intentions, validate, apply
Phase 2: ECONOMY      — accrue gold/TP (×0.1 per tick)
Phase 3: BATTLES      — decrement timers 0.1s, resolve expired
Phase 4: AUTO-DROP    — check negative income, manage grace, drop if expired
Phase 5: VICTORY      — check capital capture, unclaim loser hexes
Phase 6: DELTA        — diff vs previous tick, broadcast
```

**Why this order matters:** Counter-spend (Phase 1) runs before battles resolve (Phase 3). Economy accrues before victory check. Phase order is not changeable without careful analysis of phase interactions.

---

## State Sync Model

**Snapshot:** Full state dump. Sent on initial connect and reconnect (`prevSnapshot = nil`).

**Delta:** Only what changed. Built by `buildDelta(prev, curr)` in `internal/net/delta.go`.

Delta merge on client (`state/state.ts`):
```typescript
// Players: spread preserves fields not in delta; explicit tech fallback
players.set(id, { ...(existing ?? {}), ...p, tech: p.tech ?? existing?.tech ?? [] });
// Battles: null-coalesce to preserve active battles when field omitted
battles: msg.battles ?? state.battles
```

---

## Module Responsibilities

| Module | Owns | Does not own |
|---|---|---|
| `internal/game/` | Game rules, state transitions, validation | Network, goroutines, I/O |
| `internal/room/` | Room lifecycle, tick goroutine, pause/forfeit | Game rules |
| `internal/net/` | WebSocket, HTTP, JSON serialization, delta building | Game rules |
| `client/src/state/` | Client state mirror, apply delta/snapshot | Rendering, UI |
| `client/src/render/` | Canvas drawing, effects | Game state mutation |
| `client/src/ui/` | DOM: HUD, sidebar, tech tree, overlays | Canvas |

---

## Adding a New Game Feature (Checklist)

1. **Update CLAUDE.md first** — write the design spec before code
2. **`internal/game/constants.go`** — add numeric constants
3. **`client/src/constants.ts`** — mirror the constants
4. **`internal/game/state.go`** — add fields to `GameState`, `Player`, or `Hex` as needed
5. **`internal/game/`** — implement the mechanic (action validation, tick phase, economy)
6. **`internal/net/messages.go`** — add/update DTO fields; check omitempty safety
7. **`internal/net/delta.go`** — add diffing for new fields that shouldn't be sent every tick
8. **`client/src/state/state.ts`** — update DTOs; update `applyDelta` if new fields need special merge
9. **`client/src/ui/sidebar.ts`** or other UI files — expose the new action in the UI
10. **Tests** — game logic unit test in `internal/game/`, room test if session behavior changes, Playwright if player flow changes

---

## Adding a New Tech

A tech touches several files. Reference the Iron Grip or Garrison implementation as a model:
1. Add constant in `constants.go` / `constants.ts`
2. Add tech index in `tech.go` (unlock validation uses the index)
3. Add effect in the relevant phase (economy, combat, or action validation)
4. Update `sidebar.ts` if the tech enables new UI buttons (like Fortify)
5. Add test case to `tech_test.go`
6. Document in CLAUDE.md Tech Tree table

---

## Common Pitfalls

| Pitfall | Where | How to avoid |
|---|---|---|
| Adding omitempty to a reverting field | `messages.go` | Only omitempty fields that are one-directional (e.g., tech: false→true never reverts) |
| Capture flash using `previousOwner` from wire | `main.ts` / `applyDelta` | Compare `hex.owner` vs current `state.hexes` owner |
| Querying `data-action` attribute for destructive buttons | `sidebar.ts` | Demolish/Sell buttons exist twice in DOM — queries find both nodes |
| Adding goroutines in `internal/game/` | Any game logic | Package must remain pure; goroutines live in `internal/room/` |
| Importing non-stdlib packages in `internal/game/` | Go import | Zero external imports; use stdlib only |

---

## Rules

- **Derive, don't assume.** Every architectural change must trace to a CLAUDE.md requirement.
- **Testability is a constraint.** If you can't test it without a browser or goroutine, the architecture boundary is wrong.
- **Shortest path.** Don't add abstraction layers unless the current pattern actively creates bugs.
- **CLAUDE.md wins.** If the quick reference above conflicts with CLAUDE.md, CLAUDE.md is authoritative.
