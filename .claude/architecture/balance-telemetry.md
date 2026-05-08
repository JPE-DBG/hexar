# Balance Telemetry Architecture

## Purpose
Log match outcomes + tech paths to detect dominant strategies during playtesting.

## Minimal Telemetry Points (add if balance issues arise)

### 1. Match Outcome Log
```go
type MatchOutcome struct {
    WinnerID       PlayerID
    Duration       float64  // seconds
    WinnerTechs    []TechID // techs unlocked by winner
    WinnerHexCount int
    LoserTechs     []TechID
    LoserHexCount  int
}
```

**When:** Called in `TriggerVictory` after `state.Over = true`  
**Where:** Append to `matches.jsonl` (one line per match)  
**Why:** Reveals if specific tech combinations correlate with wins

### 2. First Blood Timer
```go
// In GameState
FirstAttackTimestamp float64 // 0 until first enemy hex attack
```

**When:** Set in `ApplyAttack` when `hs.Owner != NoPlayer && hs.Owner != action.Player` (first time only)  
**Why:** Detects "rush" strategies (Blitz Aggressor attacking at T=3min vs Economic Builder at T=8min)

### 3. Tech Unlock Timeline
```go
// Append to match log
TechUnlockTimes map[TechID]float64 // elapsed time when each tech unlocked
```

**When:** Set in `ApplyUnlockTech`  
**Why:** Shows if winners rush expensive techs (Iron Grip at 6min) vs diversify cheap techs

## Analysis Queries (post-match)

```bash
# Dominant tech combinations
cat matches.jsonl | jq -r '[.WinnerTechs[]] | sort | @csv' | sort | uniq -c | sort -rn | head

# Average time-to-victory by archetype
cat matches.jsonl | jq 'select(.WinnerTechs | contains([9,11])) | .Duration' | stats

# Win rate when first attack happens before T=5min
cat matches.jsonl | jq 'select(.FirstAttackTime < 300) | .WinnerID' | winner_rate
```

## Decision: When to Add This?

**Not in M5** — premature. Add only if playtesting reveals:
- One tech path wins >60% of matches
- Players report "no counterplay" to specific strategies
- Match duration variance is extreme (3min stomps vs 30min stalemates)

**Trigger:** After 10+ 2-player matches, if balance concerns arise.

---

**Architecture note:** Telemetry is **append-only** and **read-only during gameplay**. It doesn't affect game logic, so it's safe to add post-M5 without changing the core game loop.
