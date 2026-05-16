---
name: go-test-writer
description: "Use when: writing or fixing tests for Hexar — Go unit tests (internal/game/), Go room integration tests (internal/room/), TypeScript type checking, or Playwright E2E tests (e2e/). Covers all 3 test layers."
type: skill
---

# Test Writer Skill — Hexar Test Suite (All 3 Layers)

Writes and structures tests across all three Hexar test layers. Each layer has a different scope, tooling, and seam.

**Policy:** Update CLAUDE.md first (design spec), implement the feature, then add tests. Tests validate the spec, not the implementation.

---

## Test Layer Map

```
internal/game/      → Layer 1: Go unit tests
                      Pure functions, no I/O, no goroutines.
                      RunTick is the seam — test it directly.
                      package game (same package — can call unexported functions)

internal/room/      → Layer 2: Go integration tests
                      Real Room + fake ClientSender.
                      Drive OnConnect/EnqueueAction/OnDisconnect.
                      Real goroutines; use time.Sleep sparingly (1-2 tick durations only)
                      package room (same package — access unexported fields directly)

e2e/*.spec.ts       → Layer 3: Playwright E2E
                      Full browser + real server + Vite dev server.
                      Test user flows, not game logic.
                      Reserve for: create→join→play, disconnect/reconnect, UI layout.

TypeScript          → Type checking (not behavioral tests)
(tsc --noEmit)        Verifies DTOs, state interfaces, and applyDelta merge patterns.
                      Run after changing state/state.ts or messages.go DTOs.
```

**Rule:** If you can test it without a goroutine, test it without a goroutine. If you can test it without a browser, test it without a browser.

---

## Layer 1: Go Unit Tests (`internal/game/`)

### Seam

`RunTick(state *GameState, dt float64, actions []Action)` — pure function, no I/O, deterministic.

```go
func TestEconomyTick(t *testing.T) {
    state := buildMinimalState() // helper in same file or shared helpers
    game.RunTick(state, 0.1, nil)
    if state.Players[pid1].Gold < 0.2 { // 2/sec × 0.1s = 0.2
        t.Errorf("gold not accruing correctly")
    }
}
```

### Table-Driven Tests (Standard Pattern)

```go
func TestAttackValidation(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*game.GameState)
        action  game.Action
        wantErr error
    }{
        {
            name:    "insufficient power",
            setup:   func(s *game.GameState) { /* attacker P1, defender P2 */ },
            action:  game.Action{Type: game.ActionAttack, Player: 1, Target: game.Hex{Q: 1, R: 0}},
            wantErr: game.ErrInsufficientPower,
        },
        {
            name:    "insufficient gold",
            setup:   func(s *game.GameState) { /* attacker has 50g, needs 100g */ },
            action:  game.Action{Type: game.ActionAttack, Player: 1, Target: game.Hex{Q: 1, R: 0}},
            wantErr: game.ErrInsufficientGold,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            state := buildMinimalState()
            tt.setup(state)
            err := game.ValidateAttack(state, tt.action)
            if err != tt.wantErr {
                t.Errorf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

### State Builder Pattern

Each test file builds its own minimal state. There is no shared `testhelpers_test.go` — helpers are local to the test file that needs them.

```go
func buildMinimalState() *game.GameState {
    state := &game.GameState{
        Players: make(map[game.PlayerID]*game.Player),
        Hexes:   make(map[game.Hex]*game.HexState),
    }
    // add capital hexes, players with gold
    return state
}
```

### Calling Unexported Functions

`package game` tests can call unexported functions directly (same package). Useful for testing selection algorithms without engineering complex RunTick scenarios:

```go
// autoDropLowestHex is unexported — call directly to test selection algorithm
dropped := autoDropLowestHex(state, pid)
if dropped.Q != expectedQ {
    t.Errorf("wrong hex dropped: got %v", dropped)
}
```

### Tick Helper

```go
func runTicks(state *game.GameState, n int, actions ...game.Action) {
    for i := 0; i < n; i++ {
        var a []game.Action
        if i == 0 { a = actions }
        game.RunTick(state, game.TickDt, a)
    }
}
```

### Existing Tests

| File | What it covers |
|---|---|
| `action_test.go` | Counter-spend cap/time-cap/flip; tech unlock via action; voluntary drop clears auto-drop flag |
| `building_test.go` | Upgrade mechanics, costs, demolish refunds |
| `economy_test.go` | Gold accrual, stepped maintenance, Prosperity + Compound Growth formula |
| `combat_test.go` | Attack validation, battle resolution, instant takeover |
| `tech_test.go` | Tech unlock, Reclamation + Vanguard 0g attack cost |
| `victory_test.go` | Capital capture → WinReason "capital"; forfeit → WinReason "forfeit" |
| `autodrop_test.go` | Grace period, forced drop, refund, Resilience, selection algorithm, protected hexes |

Run: `go test ./internal/game/...`

---

## Layer 2: Go Room Integration Tests (`internal/room/`)

### Seam

`package room` — tests access unexported fields directly:
- `room.activeClients[pid]` — current active client per player
- `room.clientPlayer[client]` — reverse map (client → player ID)
- `room.remainingGrace[pid]` — saved disconnect budget

### Fake Client

```go
type fakeClient struct {
    msgs [][]byte
    mu   sync.Mutex
}

func (f *fakeClient) Send(b []byte) error {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.msgs = append(f.msgs, b)
    return nil
}

func (f *fakeClient) Close() {}
```

### Room Integration Test Pattern

```go
func TestPauseOnDisconnect(t *testing.T) {
    r := newTestRoom()
    c1, c2 := &fakeClient{}, &fakeClient{}

    r.OnConnect(c1, pid1)
    r.OnConnect(c2, pid2)  // starts the loop

    r.OnDisconnect(pid1, c1)
    time.Sleep(150 * time.Millisecond) // 1-2 tick durations

    state := r.State()
    if !state.Paused {
        t.Error("expected game to be paused after disconnect")
    }
    if state.PauseTimeLeft < 119.0 {
        t.Errorf("PauseTimeLeft too low: %v", state.PauseTimeLeft)
    }
}
```

### When to Use time.Sleep

Room tests need real `time.Sleep` because the loop goroutine is real. Keep sleeps minimal:
- 150ms (1-2 ticks) to let the loop broadcast a state change
- 300ms+ when testing timer decrement (e.g., verifying PauseTimeLeft decremented)
- Never sleep more than 500ms in a test — it's a sign the test design is wrong

### Existing Tests

| Test | What it asserts |
|---|---|
| `TestWaitingState` | `Waiting=true` until 2nd player; `Waiting=false` on 2nd connect |
| `TestWaitingStateNoGoldAccrual` | Gold unchanged while waiting |
| `TestPauseOnDisconnect` | `state.Paused=true`; `PauseTimeLeft` ≈ 120s |
| `TestUnpauseOnReconnect` | `state.Paused=false`; grace saved to `remainingGrace` |
| `TestCumulativeGrace` | 2nd disconnect uses saved grace, not reset to 120s |
| `TestForfeitEnqueued` | `PauseTimeLeft→0` → `WinReason="forfeit"` |
| `TestDuplicateConnectReplaces` | 2nd `OnConnect` replaces old client; new client active; old evicted |
| `TestDuplicateConnectDuringPause` | Reconnect during pause succeeds (not rejected as duplicate) |

Run: `go test ./internal/room/...`

---

## Layer 3: Playwright E2E (`e2e/`)

### When to Add E2E Tests

Add E2E for: new tech (player behavior changes), disconnect/reconnect flows, new UI screens, new responsive layouts. Skip E2E for: game logic (covered by Go), serialization (thin layer).

### Setup

Playwright's `webServer` config auto-starts both servers. Multi-player scenarios use two `BrowserContext` objects.

```typescript
// e2e/playwright.config.ts
webServer: [
  { command: 'go run cmd/server/main.go', url: 'http://localhost:8080/health' },
  { command: 'npm run dev --prefix client', url: 'http://localhost:5173' },
]
```

### Two-Player Test Pattern

```typescript
test('pause banner appears on opponent disconnect', async ({ browser }) => {
  const ctx1 = await browser.newContext();
  const ctx2 = await browser.newContext();
  const p1 = await ctx1.newPage();
  const p2 = await ctx2.newPage();

  // create game with p1
  await p1.goto('http://localhost:5173');
  await p1.click('#create-btn');
  const code = await p1.locator('#room-code').textContent();

  // join with p2
  await p2.goto('http://localhost:5173');
  await p2.fill('#join-code', code!);
  await p2.click('#join-btn');

  // wait for game to start
  await p1.waitForSelector('#hud');
  await p2.waitForSelector('#hud');

  // disconnect p2
  await p2.close();

  // p1 should see pause banner
  await expect(p1.locator('#pause-banner')).toBeVisible();
});
```

### Responsive Layout Testing

```typescript
test('mobile portrait layout', async ({ browser }) => {
  const ctx = await browser.newContext({
    viewport: { width: 390, height: 844 },
  });
  const page = await ctx.newPage();
  // ... join game, verify sidebar at bottom
  await expect(page.locator('#sidebar')).toHaveCSS('bottom', '0px');
});
```

### Existing Specs

| Spec | Tests |
|---|---|
| `lobby.spec.ts` | Create shows 4-char code; join with bad code errors; waiting overlay; HUD on P2 join |
| `gameplay.spec.ts` | Canvas visible; gold increases; sidebar visible |
| `connection.spec.ts` | Duplicate tab → "Already Connected"; URL hash reconnects; pause banner on disconnect |
| `responsive.spec.ts` | Desktop/mobile portrait/landscape layout verification |

Run: `make test-e2e`

---

## TypeScript Type Checking

Not behavioral tests — verifies DTOs match between Go server messages and client state interfaces.

**When to run:** After changing `internal/net/messages.go` DTOs, `client/src/state/state.ts`, or `applyDelta` merge logic.

```bash
cd client && npx tsc --noEmit
```

**Common type errors after DTO changes:**
- Adding optional `?` fields to `PlayerDTO` in `state.ts` → downstream code using them without `?? 0` null coalescing will fail
- Changing `battles` in `DeltaMsg` from required to optional → code accessing `msg.battles` directly (not `msg.battles ?? state.battles`) will have type issues

**Type safety rule for optional fields:** Fields that can be absent from delta but are always present on the full client `GameState` interface must use null coalescing at the merge point in `applyDelta`. Never use `!` non-null assertions on optional DTO fields.

---

## Run All Tests

```bash
go test ./...           # Layer 1 + Layer 2
make test-e2e           # Layer 3 (also starts servers)
cd client && npx tsc --noEmit   # TypeScript type check
```

---

## What NOT to Test

| Thing | Reason |
|---|---|
| Canvas rendering | Not testable in Go; visual review or manual check |
| WebSocket JSON framing | Thin serialisation glue; covered by E2E |
| URL hash session persistence | Browser behaviour; Playwright only |
| Victory screen text | DOM output; manual or Playwright |
| `internal/net/` package | Skip; thin layer; test the game logic it wraps |

---

## Rules

- **No mocks for `internal/game/`.** It's pure — just call the functions.
- **One assertion per test name.** If the test name has "and", split it.
- **Use `t.Helper()` in helpers** — failure lines point to the call site, not the helper.
- **Table rows are documentation.** Name each row as a scenario, not "case1".
- **stdlib only in `internal/game/` tests.** `testify` fine for integration; no external deps in game tests.
- **`t.Fatal` for bad setup.** If state setup fails (e.g. player not found), `t.Fatal` — subsequent assertions are meaningless.
- **For room tests:** always use `t.Errorf` for assertions, not `t.Logf` — `Logf` never fails the test.
