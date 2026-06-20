---
name: test-writer
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
func TestBarAccrual(t *testing.T) {
    state := buildMinimalState()
    game.RunTick(state, 0.1, nil)
    // 0.1 bar/sec × 0.1s = 0.01 per tick, per player
    if state.Players[pid1].Bar < 0.009 {
        t.Errorf("bar not accruing correctly")
    }
}
```

### Table-Driven Tests (Standard Pattern)

```go
func TestPlayTopCard(t *testing.T) {
    tests := []struct {
        name      string
        setup     func(*game.GameState)
        wantBar   float64
        wantDeck  int // expected deck length after play
    }{
        {
            name: "soldier costs 2 bar",
            setup: func(s *game.GameState) {
                s.Players[pid1].Bar = 3.0
                s.Players[pid1].Deck.Cards = []game.Card{{Type: game.CardBasicSoldier}}
            },
            wantBar:  1.0,
            wantDeck: 1, // card cycles to bottom
        },
        {
            name: "insufficient bar does nothing",
            setup: func(s *game.GameState) {
                s.Players[pid1].Bar = 1.0
                s.Players[pid1].Deck.Cards = []game.Card{{Type: game.CardBasicSoldier}}
            },
            wantBar:  1.0,
            wantDeck: 1, // card stays at top
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            state := buildMinimalState()
            tt.setup(state)
            game.RunTick(state, 0.0, []game.Action{
                {Type: game.ActionPlayTopCard, Player: pid1},
            })
            if state.Players[pid1].Bar != tt.wantBar {
                t.Errorf("bar: got %.2f, want %.2f", state.Players[pid1].Bar, tt.wantBar)
            }
        })
    }
}
```

### State Builder Pattern

Each test file builds its own minimal state. There is no shared `testhelpers_test.go` — helpers are local to the test file that needs them.

```go
const (
    pid1 = game.Player1
    pid2 = game.Player2
)

func buildMinimalState() *game.GameState {
    state := game.NewGameState()
    // Capital hexes
    cap1 := game.Hex{Q: -4, R: 0}
    cap2 := game.Hex{Q: 4, R: 0}
    state.Hexes[cap1] = &game.HexState{Owner: pid1, Capital: true}
    state.Hexes[cap2] = &game.HexState{Owner: pid2, Capital: true}
    // Players with starting deck
    state.Players[pid1] = game.NewPlayer(pid1)
    state.Players[pid2] = game.NewPlayer(pid2)
    return state
}
```

### Calling Unexported Functions

`package game` tests can call unexported functions directly (same package). Useful for testing internal helpers:

```go
// cycleTopCard is unexported — call directly to test deck cycling
cycleTopCard(p)
if p.Deck.Cards[0].Type != game.CardHexClaim {
    t.Errorf("wrong card at top after cycle")
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

### Test Files and Coverage (What to Write)

| File | Test cases |
|---|---|
| `bar_test.go` | Bar accrues at 0.1/tick; caps at 10; bar deducted on PlayTopCard; BarBoost ticks down; stacking two BarBoosts adds rates correctly |
| `deck_test.go` | Play from top cycles card to bottom; aside slot blocks on push; aside card plays and goes to bottom; cannot push when slot occupied; auto-cycle timer moves stuck card to bottom; starting deck initialized with correct 7 cards; CardCost returns correct costs |
| `units_test.go` | Unit spawns at capital hex; unit advances toward enemy capital each move tick; 1v1 combat both die after 2s; 2v1 correct HP after fight; unit consumed on capital contact; capital HP decrements; dead units removed from list |
| `shop_test.go` | BuyCard deducts bar; bought card goes to bottom of deck; shared pool quantity decrements; second player blocked when quantity 0; RemoveTopCard removes card from deck; RemoveAsideCard removes aside card |
| `victory_test.go` | Capital at 0 HP sets WinReason "capital"; damage is permanent between ticks; forfeit sets WinReason "forfeit" |

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
type MockClientSender struct {
    mu        sync.Mutex
    snapshots []*game.GameState
}

func (m *MockClientSender) SendSnapshot(state *game.GameState) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.snapshots = append(m.snapshots, state)
}
```

### Room Integration Test Pattern

```go
func TestPauseOnDisconnect(t *testing.T) {
    r := New()
    c1, c2 := &MockClientSender{}, &MockClientSender{}

    r.OnConnect(c1, pid1)
    r.OnConnect(c2, pid2)  // starts the loop

    time.Sleep(150 * time.Millisecond)
    r.OnDisconnect(pid1, c1)

    if !r.state.Paused {
        t.Error("expected game to be paused after disconnect")
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
| `TestWaitingStateNoBarAccrual` | Bar unchanged while waiting |
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

Add E2E for: disconnect/reconnect flows, new UI screens, new responsive layouts, full game flow changes. Skip E2E for: game logic (covered by Go), serialization (thin layer).

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

  await p1.goto('http://localhost:5173');
  await p1.click('#create-btn');
  const code = await p1.locator('#room-code').textContent();

  await p2.goto('http://localhost:5173');
  await p2.fill('#join-code', code!);
  await p2.click('#join-btn');

  await p1.waitForSelector('#hud');
  await p2.waitForSelector('#hud');

  await p2.close();

  await expect(p1.locator('#pause-banner')).toBeVisible();
});
```

### E2E Tests to Write (Client Not Yet Built)

Once the client is implemented, priority E2E specs:

| Spec | Tests |
|---|---|
| `lobby.spec.ts` | Create shows 4-char code; join with bad code errors; waiting overlay; HUD on P2 join |
| `gameplay.spec.ts` | Canvas visible; bar meter fills; deck panel shows top card; SPACE plays top card; E pushes to aside |
| `connection.spec.ts` | Duplicate tab → "Already Connected"; URL hash reconnects; pause banner on disconnect |

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
- Changing `units` list — always present, never optional; code accessing it directly is safe
- DeckDTO fields (cards, asideCard, autoTimer) — asideCard can be null; code must handle `null` not just `undefined`

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
