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
| Go tick loop | < 5ms per tick (100ms interval) | Lightweight | Breaks if RunTick O(n²) on hex count |
| Delta message | < 500 bytes steady state | ~300-400 bytes | Breaks at 4-player or large maps |
| Canvas render (game state) | < 16ms at 100ms intervals | Fine at 70 hexes | Breaks with particle effects |
| rAF animation loop | < 16ms every frame (60fps) | Not yet implemented | Will matter once effects added |
| WebSocket round-trip | < 100ms | Fine on LAN | Matters for remote players |

---

## Workflow

### 1. Identify Which Layer Is Slow

**Symptoms → Layer:**

| Symptom | Likely layer |
|---|---|
| Game state updates feel laggy (>100ms stale) | Go tick loop or WebSocket |
| Animations stutter even when game state is fine | Canvas rAF loop |
| High CPU in browser, game responsive | Canvas over-drawing |
| High CPU in Go process | RunTick hot path or GC |
| Delta messages growing over time | Delta diff logic, hex state accumulation |
| Memory grows indefinitely | Goroutine leak or slice append without cap |

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
- `json.Marshal` in `client.SendSnapshot` — can be hot if called every tick for many clients
- `runtime.GC` — if > 5% of CPU, you have allocation pressure in the hot path

### 3. Profiling: Canvas Renderer

Open Chrome DevTools → Performance → Record 5 seconds of gameplay.

Look for:
- **Long tasks (> 50ms)** in the main thread during game renders
- **Forced reflows** — reading layout properties after DOM writes in the same frame
- **Canvas `drawImage` or `fillRect` call count** — should scale linearly with hex count, not quadratically
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

If deltas are growing:
- **Players always included** — intentional, they change every tick (gold/TP)
- **Battles always included** — intentional, timers change every tick
- **Hex changes** — should be near zero in a static board; spike during combat is expected

Optimisation options if needed (in order of impact):
1. **MessagePack** instead of JSON — 40-60% size reduction, no schema change needed
2. **Omit unchanged player fields** — only send gold/TP/vanguardTimer when they change
3. **Quantise floats** — gold to 1 decimal, timers to 2 decimals (already reasonable)

### 5. Canvas vs PixiJS Decision

Migrate to PixiJS **only if:**
- rAF loop consistently > 12ms with effects enabled (measured, not estimated)
- You need effects that require GPU: particle systems, shaders, glow filters, fog of war
- Canvas `drawImage` on a 70-hex board is measurably slow (it won't be)

Stay on Canvas if:
- Adding gradients and border glow only (CSS-style operations, fine on Canvas)
- Effect layer is limited to capture flash + floating numbers
- Performance budget is met

Migration cost: ~500 lines rewrite of `client/src/render/renderer.ts`. Only justified by a profiling result, not anticipated need.

### 6. Common Hexar-Specific Bottlenecks

**Hot path allocation in RunTick:**
```go
// Bad: allocates slice every tick
func (r *Room) drainActions() []game.Action {
    var pending []game.Action  // new alloc every call
    ...
}

// Better: pre-allocate with expected capacity
pending := make([]game.Action, 0, 16)
```

**Canvas over-drawing:**
- Only redraw hexes that changed (dirty rect tracking)
- Skip rendering when `state.Waiting = true` or `state.Paused = true` — board doesn't change

**Goroutine leak — room never stopped:**
```go
// Verify rooms are cleaned up
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
- [ ] `RunTick` completes in < 1ms on a 70-hex board
- [ ] No allocation in the tick hot path (use `go test -benchmem`)
- [ ] No goroutine leaks after game over (check pprof goroutine endpoint)
- [ ] Delta messages < 500 bytes steady state (existing log covers this)

### Canvas Renderer
- [ ] Game-state render (100ms) < 16ms
- [ ] rAF animation loop < 16ms per frame
- [ ] No forced reflows in render path
- [ ] Static board (pause/waiting) skips unnecessary redraws

---

## Rules

- **Measure first, always.** "This might be slow" is not a reason to optimise.
- **One change at a time.** Two changes = you don't know which one helped.
- **The 500-byte delta budget is the server performance KPI.** If it's met, the server is fine.
- **Canvas is fast enough for 70 hexes.** Don't pre-emptively migrate to PixiJS — wait for a measured reason.
- **GC pressure matters in Go.** Allocating in the tick hot path (called 10/sec) accumulates. Use benchmarks with `-benchmem` to catch it.
