---
name: balance-sim
description: "Use when: running economy projections, simulating upgrade cost payoffs, tracing TP timelines, or verifying balance numbers from CLAUDE.md before they get hardcoded"
type: skill
---

# Balance Sim Skill — Hexar Quantitative Validation

Runs traced, step-by-step numeric simulations using **exact current CLAUDE.md values**. Proves or disproves claims about game balance with math.

**Not for:** suggesting mechanic redesigns or evaluating design health — use `/game-designer` for that.

**Prerequisites:** None. Can be invoked anytime.

## Workflow

### 1. Load Constants from CLAUDE.md (MANDATORY)
- **Read CLAUDE.md at the start of every simulation.** Extract all relevant values fresh.
- Never use numbers remembered from a previous conversation or cached in this file.
- List the constants used at the top of your output so the user can verify.

### 2. Define Simulation Parameters
- What scenario is being tested? (e.g., "time to reach Tech Level 4 with 3 Research buildings")
- What are the initial conditions? (starting hexes, gold, time)
- What assumptions are made? (player behavior, expansion rate, opponent interference)

### 3. Run Step-by-Step Trace
- Show every intermediate value — not just final results.
- Use exact arithmetic, no rounding unless stated.
- Mark each time step clearly (T=Xs: state → calculation → new state).
- If the trace branches (player makes a choice), show both paths.

### 4. Compare Against Design Target
- What does CLAUDE.md say the result SHOULD be?
- Does the simulation match? If not, by how much?
- Is the discrepancy a problem or within acceptable range?

### 5. Report Finding
- If numbers match design intent: confirm with proof.
- If numbers diverge: quantify the gap and suggest a specific constant adjustment.
- Hand off to `/game-designer` if the problem requires a mechanic redesign rather than a number tweak.

---

## Common Simulation Types

### Economy Projection
Trace gold accumulation over time for a given expansion/building strategy.
- Input: expansion rate, building placement timing, tech unlocks
- Output: gold/sec and total gold at key time points
- Validate: can player afford key actions at intended time?

### Upgrade Payoff (Break-Even)
Calculate how long until an upgrade pays for itself.
- Input: upgrade cost, income delta from upgrade
- Output: break-even time in seconds
- Validate: is break-even fast enough to matter in a 30-min game?

### Maintenance Threshold
Trace at what hex count income goes negative (auto-drop trigger).
- Input: number of hexes, number/level of Economy buildings
- Output: exact hex count where net income ≤ 0
- Validate: does it match CLAUDE.md's claimed threshold?

### TP Timeline
Calculate time to reach specific Tech Levels.
- Input: number of Research buildings, their levels, start time
- Output: time to unlock each tech tier
- Validate: is Tech Dominance achievable in the intended time window?

### Battle Outcome Matrix
Given Power values, enumerate all possible outcomes including counter-spend.
- Input: attacker Power, defender Power, Garrison status
- Output: result (fail/standard/instant), battle duration, counter-spend options
- Validate: do battles resolve in intended time frame?

---

## Output Format

```
## Simulation: [Name]
**Question:** [What we're trying to prove/disprove]
**Constants loaded from CLAUDE.md:**
  [list each value used and where it appears in CLAUDE.md]

**Setup:** [Initial conditions and assumptions]

**Trace:**
  T=0s:   [state] → [calculation] → [result]
  T=10s:  [state] → [calculation] → [result]
  ...

**Result:** [Final answer]
**Design target (from CLAUDE.md):** [What it should be]
**Verdict:** MATCH / MISMATCH — [explanation]

## Adjustment (if mismatch)
**Suggested constant change:** [specific value → new value]
**Revised trace (abbreviated):** [show key steps with new value]
**New result:** [confirm it hits design target]
```

---

## Quick Reference (verify against CLAUDE.md — these may be outdated)

```
Economy:
  Base income: 2/sec per hex
  Maintenance: 1/sec (hexes 1-10), 2/sec (11-20), 3/sec (21+)
  Net per hex: T1=+1/sec, T2=0/sec, T3=-1/sec (without Economy building)
  Economy building: +50% bonus, formula = (base + 0.6 × level) × 1.5
    L1=3.9/sec, L2=4.8/sec, L3=5.7/sec; delta +0.9/sec per level
  Upgrade formula (all buildings): BuildCost × 2^currentLevel
  Economy upgrade costs (BuildCost=60): L1=60, L2=120, L3=240, L4=480
  Defense upgrade costs (BuildCost=60): L1=60, L2=120, L3=240, L4=480
  Research upgrade costs (BuildCost=80): L1=80, L2=160, L3=320, L4=640
  All demolish refund: BuildCost × (2^level - 1) × 0.5
  No separate build action — upgrading empty hex to L1 costs BuildCost

Combat:
  Unclaimed hex: 10 gold, instant
  Enemy hex: 100 gold, battle
  Battle duration: 5 + (Attacker Power + Defender Power) / 2 seconds
  Counter-spend: 50 gold/sec, cap = min(+3, seconds_remaining)
  Garrison shares counter-spend cap (+2 max from Garrison, +3 total)
  Instant takeover: Power diff > 3

Tech (12 total, 460 TP to unlock all — no game unlocks everything):
  Research buildings: +0.2 TP/sec per level
  Blitz: 20 TP | Fortify: 20 TP | Prosperity: 25 TP | Reclamation: 25 TP
  Vanguard: 30 TP | Garrison: 30 TP | War Chest: 30 TP | Supply Lines: 40 TP
  Resilience: 45 TP | Iron Grip: 55 TP | Compound Growth: 65 TP | Siege Mastery: 75 TP
  Tech bonuses:
    Prosperity: +1.0/sec per Economy building (flat, after Compound Growth multiplier)
    Compound Growth: Economy output ×1.25 (before Prosperity flat bonus)
    Supply Lines: maintenance 0.9/1.8/2.7 per tier (vs 1.0/2.0/3.0)
    Iron Grip: all owned hexes +1 Power (static, always-on)
    Garrison: adjacent owned hexes +1 Power in defense, cap +2 (defense-battle only)
    Reclamation: recapture own hex costs 50g (−50g)
    Vanguard: after capture, next attack within 12s costs 50g (−50g)
    Reclamation + Vanguard stack: 100 − 50 − 50 = 0g (free attack)

Victory:
  Capital Capture only — capture the enemy capital hex to win instantly.
  All loser hexes become unclaimed immediately.
```

---

## Rules

- **Source of truth:** CLAUDE.md, loaded fresh every invocation. Use quick reference above for orientation only — if it conflicts with CLAUDE.md, CLAUDE.md wins.
- **Show your work:** Every calculation visible, no "and therefore the result is X"
- **No mechanic opinions:** If the problem is the mechanic design (not the number), say "hand off to /game-designer" and stop
- **Precision:** Use exact values. Only round for display if stated explicitly
- **Assumptions visible:** State all behavioral assumptions (e.g., "player expands 1 hex every 5 seconds")
