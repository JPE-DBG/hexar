# M6 Disconnect Handling Architecture

## Requirements (from MILESTONES.md M6)
- Disconnect → game continues for 30s
- If no reconnect within 30s → forfeit
- Reconnect → full snapshot resync

## State Already Reconnect-Ready (M5)

All time-sensitive state is server-authoritative and propagated via `SnapshotMsg`:

```go
// HexDTO (messages.go)
type HexDTO struct {
    // ... existing fields
    FortifyTimer  float64 `json:"fortifyTimer"`   // ✅ Already synced
    PreviousOwner int     `json:"previousOwner"`  // ✅ Already synced
}

// PlayerDTO
type PlayerDTO struct {
    // ... existing fields
    // VanguardTimer NOT synced — intentional (see below)
}

// SnapshotMsg
type SnapshotMsg struct {
    // ... existing fields
    Over   bool `json:"over"`    // ✅ Already synced
    Winner int  `json:"winner"`  // ✅ Already synced
}
```

## M6 Changes Required

### Server (Go)

1. **Room disconnect handling** (room/room.go):
   ```go
   type Room struct {
       // ... existing
       disconnectTimers map[PlayerID]*time.Timer
   }
   
   func (r *Room) OnClientDisconnect(playerID PlayerID) {
       timer := time.AfterFunc(30*time.Second, func() {
           r.ForfeitPlayer(playerID)
       })
       r.disconnectTimers[playerID] = timer
   }
   
   func (r *Room) OnClientReconnect(playerID PlayerID) {
       if timer, ok := r.disconnectTimers[playerID]; ok {
           timer.Stop()
           delete(r.disconnectTimers, playerID)
       }
       // Send full snapshot
       r.SendSnapshot(playerID)
   }
   ```

2. **Forfeit logic** (game/victory.go):
   ```go
   func ForfeitPlayer(state *GameState, playerID PlayerID) {
       // Find opponent
       var opponent PlayerID
       for pid := range state.Players {
           if pid != playerID {
               opponent = pid
               break
           }
       }
       TriggerVictory(state, opponent)  // Reuse existing victory logic
   }
   ```

### Client (TypeScript)

1. **WebSocket reconnect** (net/connection.ts):
   ```ts
   private reconnectAttempts = 0;
   private maxReconnectAttempts = 5;
   
   private onClose() {
       if (this.reconnectAttempts < this.maxReconnectAttempts) {
           setTimeout(() => {
               this.reconnectAttempts++;
               this.connect();  // Reuse existing connect logic
           }, 2000);  // 2s backoff
       } else {
           this.showDisconnectOverlay();
       }
   }
   
   private onOpen() {
       this.reconnectAttempts = 0;
       // Server sends snapshot automatically
   }
   ```

## VanguardTimer Sync Decision

**Question:** Should `VanguardTimer` be synced to clients?

**Option A: Sync VanguardTimer** (show exact countdown in UI)
- Pros: Player sees "Vanguard active: 8s remaining"
- Cons: Extra 8 bytes per player per snapshot; UI clutter

**Option B: Don't sync** (client only sees "Vanguard active" boolean)
- Pros: Simpler; timer is server-only concern
- Cons: Player doesn't know when discount expires

**Recommendation:** **Option B for M6 MVP** — client can infer "Vanguard is active" from `player.tech[4]` and show "Vanguard Ready" without exact timer. If UI polish is desired post-M6, add timer sync then.

## Testing Strategy

1. **Unit test:** `TestForfeitOnDisconnect` (victory_test.go)
2. **Integration test:** Manual two-tab test:
   - Tab 1: Start game, unlock Fortify, fortify a hex
   - Tab 1: Close tab (simulates disconnect)
   - Tab 2: Wait 30s, verify Tab 1 forfeits
   - Tab 1: Reopen, verify reconnect + state sync (fortified hex still visible)

## Known Edge Cases

1. **Both players disconnect simultaneously:** Last to disconnect after 30s forfeits (or draw if both timeout simultaneously — requires tie-break rule)
2. **Disconnect during battle resolution:** Battle resolves on server; reconnecting player sees outcome in snapshot
3. **Disconnect during auto-drop grace:** Grace continues ticking; player may lose hexes while disconnected (intended — creates risk)

## No State Model Changes Required

All M5 state additions (`FortifyTimer`, `PreviousOwner`, `Over`, `Winner`) are already in DTOs. M6 just needs to wire up reconnect logic and forfeit timers.
