---
name: go-test-writer
description: "Use when: writing Go tests for Hexar game logic, room lifecycle, or lobby behaviour; designing table-driven tests for RunTick; creating fake ClientSender mocks for integration tests; deciding what to unit-test vs integration-test vs leave for Playwright."
type: skill
---

# Go Test Writer Skill — Hexar Test Suite

Writes and structures Go tests for Hexar. Covers three layers: pure game logic (`internal/game/`), room/lobby integration (`internal/room/`, `internal/lobby/`), and guidance on what to leave for Playwright E2E.

**Prerequisites:** None. Can be invoked standalone.

---

## Test Layer Map

```
internal/game/      → Unit tests. Pure functions, no I/O, no goroutines.
                      RunTick is the seam — test it directly.

internal/room/      → Integration tests. Real Room + fake ClientSender.
internal/lobby/     → Drive OnConnect/EnqueueAction/OnDisconnect.
                      Use real timers sparingly; prefer injecting state.

internal/net/       → Skip. Thin serialisation layer; covered by E2E.

Playwright (E2E)    → Full browser + real server. Test user flows,
                      not game logic. Reserve for: create→join→play→victory,
                      disconnect/reconnect, already-connected tab.
```

**Rule:** If you can test it without a goroutine, test it without a goroutine. If you can test it without a browser, test it without a browser.

---

## Workflow

### 1. Identify the Seam

For `internal/game/`: the seam is `RunTick(state *GameState, dt float64, actions []Action)`. It is a pure function — no I/O, no randomness, deterministic. Every game mechanic test drives it directly.

```go
func TestEconomyTick(t *testing.T) {
    state := game.NewGameState()
    // set up state...
    game.RunTick(state, 0.1, nil)
    // assert state...
}
```

For `internal/room/`: the seam is the `ClientSender` interface. Implement a fake:

```go
type fakeClient struct {
    snapshots []*game.GameState
}

func (f *fakeClient) SendSnapshot(s *game.GameState) {
    f.snapshots = append(f.snapshots, s)
}
```

### 2. Table-Driven Tests (Standard Pattern)

Always use table-driven tests for game logic. Each row is a scenario, not a separate test function.

```go
func TestAttackValidation(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*game.GameState)
        action  game.Action
        wantErr error
    }{
        {
            name: "insufficient power",
            setup: func(s *game.GameState) { /* attacker P1, defender P2 */ },
            action:  game.Action{Type: game.ActionAttack, Player: 1, Target: game.Hex{Q: 1, R: 0}},
            wantErr: game.ErrInsufficientPower,
        },
        {
            name: "insufficient gold",
            // ...
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            state := buildState(tt.setup)
            err := game.ValidateAttack(state, tt.action)
            if err != tt.wantErr {
                t.Errorf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

### 3. State Builder Helper

Each test needs a valid starting state. Write a builder once, reuse everywhere:

```go
// testhelpers_test.go (in internal/game/ package)
func twoPlayerState() *game.GameState {
    state := game.NewGameState()
    state.Players[1] = &game.Player{ID: 1, Gold: 500}
    state.Players[2] = &game.Player{ID: 2, Gold: 500}
    // place capitals, a few hexes each
    return state
}
```

### 4. Tick Sequences

Many mechanics only trigger after multiple ticks. Helper to run N ticks:

```go
func runTicks(state *game.GameState, n int, actions ...game.Action) {
    for i := 0; i < n; i++ {
        var a []game.Action
        if i == 0 { a = actions }
        game.RunTick(state, game.TickDt, a)
    }
}
```

### 5. Room Integration Tests

```go
func TestRoomPauseOnDisconnect(t *testing.T) {
    r := room.New()
    c1, c2 := &fakeClient{}, &fakeClient{}

    r.OnConnect(c1, 1)
    r.OnConnect(c2, 2)  // starts the loop

    r.OnDisconnect(1, c1)

    // give loop one tick to broadcast
    time.Sleep(150 * time.Millisecond)

    last := c2.snapshots[len(c2.snapshots)-1]
    if !last.Paused {
        t.Error("expected game to be paused after disconnect")
    }
}
```

Note: room integration tests need real time.Sleep because the loop goroutine is real. Keep sleeps minimal (1-2 tick durations). If tests become flaky, increase sleep or inject a test clock.

---

## Priority Test List for Hexar

### Layer 1: Game Logic (`internal/game/`)

| Test | Why critical |
|---|---|
| Economy tick: gold accrues correctly for N hexes | Base of all game decisions |
| Maintenance: stepped cost triggers at 10/20 hex boundaries | Common bug surface |
| Economy + Compound Growth + Prosperity stack correctly | Triple interaction, easy to get wrong |
| Attack validation: all error conditions | Gold spent incorrectly = game-breaking |
| Battle resolution: win/lose/draw with garrison boost | Core combat outcome |
| Instant takeover: diff > 3 bypasses battle | Edge case that surprises players |
| Forfeit: all loser hexes become unclaimed | Victory condition correctness |
| Capital capture: triggers TriggerVictory with reason "capital" | Win reason must be set |
| Disconnect forfeit: sets WinReason "forfeit" | Win reason must be set |
| Auto-drop: triggers when net income goes negative | Economy safety valve |
| Vanguard timer: refreshes on each capture | Tech interaction |
| Reclamation + Vanguard stack: cost reaches 0 | CLAUDE.md explicitly calls this out |
| Tech unlock: TP deducted, bool set | Basic correctness |

### Layer 2: Room/Lobby (`internal/room/`, `internal/lobby/`)

| Test | Why critical |
|---|---|
| Waiting: loop does not start until 2nd player connects | Gold head-start bug |
| Pause: state.Paused = true immediately on disconnect | Core fairness feature |
| Unpause: state.Paused = false on reconnect | Reconnect flow |
| Cumulative grace: second disconnect gets remaining budget | Not 120s fresh each time |
| Forfeit: enqueued when PauseTimeLeft hits 0 | Must not be lost |
| Duplicate connect: IsConnected prevents second connection | Already-connected tab bug |
| Room full: third join rejected | Lobby correctness |

---

## What NOT to Test in Go

- **Rendering:** Canvas drawing is not testable in Go — leave for Playwright or visual review
- **WebSocket framing:** `net/` serialisation is thin glue — test the game logic it wraps, not the JSON encoding
- **URL hash session persistence:** Browser behaviour — Playwright only
- **Victory screen text:** DOM output — manual or Playwright

---

## Rules

- **No mocks for `internal/game/`.** It's pure — just call the functions. Mocks here mean the architecture is wrong.
- **One assertion per test name.** If the test name has "and", split it.
- **Use `t.Helper()` in helpers** so failure lines point to the call site, not the helper.
- **Table rows are documentation.** Name each row as a sentence describing the scenario, not "case1".
- **stdlib only for game tests.** `testify` is fine for integration tests. Don't add dependencies to `internal/game/`.
- **Fail fast.** If setup is wrong (e.g. player not found), `t.Fatal` not `t.Error` — subsequent assertions are meaningless.
