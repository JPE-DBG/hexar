---
name: performance
description: "Use when: the Canvas renderer feels slow or drops frames, the 100ms tick loop is running late, WebSocket delta messages are growing too large, profiling the Go server under load, or evaluating whether to migrate from Canvas to PixiJS."
type: skill
---

# Performance Skill — Hexar Profiling & Optimisation

Diagnoses and fixes performance issues in Hexar across three layers: Go server (tick loop, GC pressure), WebSocket (message size, frequency), and client rendering (Canvas 60fps animation vs 100ms game-state ticks).

**Prerequisites:** Measure before optimising. Never optimise without a profiling result showing the bottleneck.

---

## Performance Budgets

| Layer | Budget | Current state | Break point |
|---|---|---|---|
| Go tick loop | < 5ms per tick (100ms interval) | Lightweight | Breaks if RunTick O(n²) on unit count |
| Delta message | < 500 bytes steady state | TBD (client not yet built) | Breaks at 4-player or many units |
| Canvas render (game state) | < 16ms at 100ms intervals | Fine at 70 hexes | Breaks with many unit tokens + effects |
| rAF animation loop | < 16ms every frame (60fps) | Not yet implemented | Will matter once unit march effects added |
| WebSocket round-trip | < 100ms | Fine on LAN | Matters for remote players |

---

## Workflow

### 1. Identify Which Layer Is Slow

**Symptoms → Layer:**

| Symptom | Likely layer |
|---|---|
| Game state updates feel laggy (>100ms stale) | Go tick loop or WebSocket |
| Unit tokens stutter even when game state is fine | Canvas rAF loop |
| High CPU in browser, game responsive | Canvas over-drawing |
| High CPU in Go process | RunTick hot path or GC |
| Delta messages growing over time | Player/unit list serialization |
| Memory grows indefinitely | Goroutine leak or unit slice append |

### 2. Profiling: Go Server

```bash
# CPU profile: run server with pprof enabled
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Heap profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine leak check
curl http://localhost:8080/debug/pprof/goroutine?debug=1
```

Add to `cmd/server/main.go` during profiling sessions:
```go
import _ "net/http/pprof"
```

Key things to look for:
- `game.RunTick` — should be < 1ms per call
- `json.Marshal` in `client.SendSnapshot` — hot if many clients or large unit lists
- `runtime.GC` — if > 5% of CPU, you have allocation pressure in the hot path

### 3. Profiling: Canvas Renderer

Open Chrome DevTools → Performance → Record 5 seconds of gameplay.

Look for:
- **Long tasks (> 50ms)** in the main thread during game renders
- **Forced reflows** — reading layout properties after DOM writes in the same frame
- **Canvas `drawImage` or `fillRect` call count** — should scale linearly with hex+unit count
- **`requestAnimationFrame` callback duration** — must stay < 16ms

Quick Canvas profiling in code:
```typescript
const t0 = performance.now();
renderer.render();
const dt = performance.now() - t0;
if (dt > 16) console.warn(`slow render: ${dt.toFixed(1)}ms`);
```

### 4. Delta Message Size

Current budget: 500 bytes steady state. Check with the existing log:
```go
if _, isDelta := msg.(*DeltaMsg); isDelta && len(data) > maxDeltaBytes {
    log.Printf("delta over budget: %d bytes", len(data))
}
```

**What's always included in deltas (intentional — they change every tick):**
- Players (bar, deck, BarBoosts) — continuously changing
- Units (all units, all positions) — moves and dies every tick; empty list correctly clears client
- Elapsed — always changes

**What's diffed (only sent when changed):**
- Hex changes (owner, HasBuilding) — sparse in steady state

If deltas are growing:
- **Many units on screen** — unit list scales linearly with active soldiers; expected spike during pushes
- **Deck serialization** — 7–20 cards with type IDs; should be small
- **BarBoosts array** — only grows when +15% cards are active; bounded

Optimisation options if needed (in order of impact):
1. **MessagePack** instead of JSON — 40-60% size reduction, no schema change needed
2. **Omit BarBoosts when empty** — only send boosts array when at least 1 is active
3. **Quantise bar** — send bar to 2 decimal places (already reasonable)

### 5. Canvas vs PixiJS Decision

Migrate to PixiJS **only if:**
- rAF loop consistently > 12ms with effects enabled (measured, not estimated)
- You need effects that require GPU: particle systems, shaders, glow filters
- Canvas is measurably slow on the 70-hex board + unit tokens (it won't be)

Stay on Canvas if:
- Unit token rendering is a circle + letter + HP dots — trivial Canvas ops
- Effect layer is limited to pulse rings and bar fill animations
- Performance budget is met

Migration cost: ~500 lines rewrite of `client/src/render/renderer.ts`. Only justified by a profiling result.

### 6. Common Hexar-Specific Bottlenecks

**Hot path allocation in RunTick — unit slice operations:**
```go
// Bad: allocates new slice every tick in removeDeadUnits
alive := make([]*Unit, 0, len(state.Units))

// Better: reuse slice backing with [:0]
alive := state.Units[:0]
for _, u := range state.Units {
    if u.HP > 0 { alive = append(alive, u) }
}
state.Units = alive
```
(This is already the pattern in `units.go` — verify it stays this way.)

**Unit march animation blocking game state:**
- Unit positions update every 2 seconds in game state (1 hex / SoldierMovePeriod)
- Visual march must be a client-side interpolation between positions, not tied to server ticks
- Never use server tick timing to drive CSS/canvas animations — they're independent loops

**Canvas over-drawing:**
- Skip rendering when `state.Waiting = true` or `state.Paused = true` — board doesn't change
- Only redraw unit tokens that moved since last frame

**Goroutine leak — room never stopped:**
```go
// After game over, r.Stop() must be called or the ticker goroutine leaks
```

**JSON marshal on every tick:**
- `BuildSnapshot` allocates a new `SnapshotMsg` every tick per client — normal, but watch under load
- Consider snapshot pooling if > 100 concurrent rooms

---

## Optimisation Checklist

### Before Optimising
- [ ] Profiling result in hand — specific function + % CPU or bytes identified
- [ ] Baseline measurement recorded (before state)
- [ ] Change is isolated to one variable

### Go Server
- [ ] `RunTick` completes in < 1ms on a 70-hex board with 10+ active units
- [ ] No allocation in the tick hot path (use `go test -benchmem`)
- [ ] No goroutine leaks after game over (check pprof goroutine endpoint)
- [ ] Delta messages < 500 bytes steady state (existing log covers this)

### Canvas Renderer
- [ ] Game-state render (100ms) < 16ms
- [ ] rAF animation loop < 16ms per frame
- [ ] No forced reflows in render path
- [ ] Static board (pause/waiting) skips unnecessary redraws
- [ ] Unit march is visual interpolation — never blocks on server tick

---

## Rules

- **Measure first, always.** "This might be slow" is not a reason to optimise.
- **One change at a time.** Two changes = you don't know which one helped.
- **The 500-byte delta budget is the server performance KPI.** If it's met, the server is fine.
- **Canvas is fast enough for 70 hexes + 20 unit tokens.** Don't pre-emptively migrate to PixiJS — wait for a measured reason.
- **GC pressure matters in Go.** Allocating in the tick hot path (called 10/sec) accumulates. Use benchmarks with `-benchmem` to catch it.
