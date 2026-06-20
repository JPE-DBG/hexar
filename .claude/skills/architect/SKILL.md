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
- Server owns all game state — bar, deck, units, hex ownership, capital HP
- Client never modifies game state — only sends intentions (actions)
- Client maintains a mirror of server state via snapshot+delta sync
- Effects (unit march animation, card flip) are client-side only, never in game state

---

## Key Invariants

These must be preserved across any change:

1. **`internal/game/` zero external imports.** `RunTick(state, dt, actions)` is a pure function. No I/O, no goroutines, no network. Fully testable with `go test` alone.

2. **One goroutine per room.** Actions arrive via WebSocket, queued into a channel, processed at the next tick boundary.

3. **Delta packets < 500 bytes at idle.** Player state (bar, deck) always sent — it changes every tick. Hex changes only sent when owner or building status changes.

4. **Player state always sent in full.** Bar, deck order, aside slot, and BarBoosts change continuously. Never omitempty these — stale client state would silently break card display.

5. **Units always sent in full.** The unit list is sent complete every tick. An empty list correctly clears all units from the client (units die and positions change every tick).

6. **Canvas + DOM split preserved.** Hex grid and unit tokens on Canvas, UI panels as DOM. The deck panel, shop panel, bar meter, and HUD are DOM overlays.

---

## Tick Execution Order (6 phases, strict sequence)

```
Phase 1: PROCESS ACTIONS — drain queued player actions (PlayCard, BuyCard, PushAside, ClaimHex)
Phase 2: BAR            — accrue bar for each player (×0.1 per tick); tick down BarBoosts
Phase 3: UNITS          — advance unit positions, resolve combat, deal capital damage
Phase 4: DECK           — advance auto-cycle timers, move stuck cards to bottom
Phase 5: VICTORY        — check capital HP, trigger win
Phase 6: DELTA          — diff vs previous tick, broadcast to clients
```

**Why this order matters:** Actions (Phase 1) run before bar accrues (Phase 2) — a player cannot spend bar that accrues in the same tick. Units advance (Phase 3) before victory is checked (Phase 5) — a unit reaching the capital deals damage before the win condition is evaluated.

---

## State Sync Model

**Snapshot:** Full state dump. Sent on initial connect and reconnect (`prevSnapshot = nil`).

**Delta:** Only what changed. Built by `buildDelta(prev, curr)` in `internal/net/delta.go`.

Delta merge on client (`state/state.ts`):
```typescript
// Players: always sent in full (bar/deck change every tick)
players.set(id, p);
// Units: always sent in full (positions change; empty list clears all)
units = msg.units;
// Hexes: only changed hexes sent; spread preserves unchanged hexes
for (const hex of msg.hexChanges ?? []) hexes.set(key(hex), hex);
```

---

## Module Responsibilities

| Module | Owns | Does not own |
|---|---|---|
| `internal/game/` | Game rules, state transitions, validation | Network, goroutines, I/O |
| `internal/room/` | Room lifecycle, tick goroutine, pause/forfeit | Game rules |
| `internal/net/` | WebSocket, HTTP, JSON serialization, delta building | Game rules |
| `client/src/state/` | Client state mirror, apply delta/snapshot | Rendering, UI |
| `client/src/render/` | Canvas drawing (hex grid, unit tokens, effects) | Game state mutation |
| `client/src/ui/` | DOM: HUD (bar meter), deck panel, shop panel, overlays | Canvas |

---

## Adding a New Game Feature (Checklist)

1. **Update CLAUDE.md first** — write the design spec before code
2. **`internal/game/constants.go`** — add numeric constants
3. **`client/src/constants.ts`** — mirror the constants
4. **`internal/game/state.go`** — add fields to `GameState`, `Player`, `Deck`, or `Unit` as needed
5. **`internal/game/`** — implement the mechanic in the appropriate file (bar.go, deck.go, units.go, shop.go)
6. **`internal/net/messages.go`** — add/update DTO fields; never omitempty continuously-changing fields
7. **`internal/net/delta.go`** — add diffing only for fields that rarely change (hex ownership, HasBuilding); always-changing fields (bar, deck, units) skip diffing
8. **`client/src/state/state.ts`** — update DTOs; update `applyDelta` merge logic
9. **`client/src/ui/`** — expose the new action in the deck panel or shop panel
10. **Tests** — game logic unit test in `internal/game/`, room test if session behavior changes

---

## Adding a New Card Type

A new card type touches several files:

1. Add `CardXxx CardType = iota` constant in `state.go`
2. Add `CostXxx = x.x` in `constants.go` and `constants.ts`
3. Add case to `CardCost()` in `deck.go`
4. Add case to `resolveCardEffect()` in `deck.go` (or `isCardPlayable()` if it has special playability conditions)
5. Add case to `InitStartingDeck()` if it belongs in the starting deck
6. Add `CardDTO` type mapping in `state.ts` on the client
7. Add card rendering to the deck panel UI
8. Test in `deck_test.go`

---

## Common Pitfalls

| Pitfall | Where | How to avoid |
|---|---|---|
| Adding omitempty to bar or deck fields | `messages.go` | Bar and deck change every tick — omitempty would leave stale client values |
| Sending unit list as delta | `delta.go` | Units always sent in full; partial unit delta would leave dead units on client |
| Hex Claim card playing immediately without target | `deck.go` | HexClaim deducts bar and cycles card, but hex assignment happens via a separate ActionClaimHex sent by client |
| Adding goroutines in `internal/game/` | Any game logic | Package must remain pure; goroutines live in `internal/room/` |
| Importing non-stdlib packages in `internal/game/` | Go import | Zero external imports; use stdlib only |
| Using `Distance()` as standalone function | `units.go` | Distance is a method on `Hex`: `from.Distance(target)` |

---

## Rules

- **Derive, don't assume.** Every architectural change must trace to a CLAUDE.md requirement.
- **Testability is a constraint.** If you can't test it without a browser or goroutine, the architecture boundary is wrong.
- **Shortest path.** Don't add abstraction layers unless the current pattern actively creates bugs.
- **CLAUDE.md wins.** If the quick reference above conflicts with CLAUDE.md, CLAUDE.md is authoritative.
