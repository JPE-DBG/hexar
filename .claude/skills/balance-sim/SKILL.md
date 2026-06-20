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
- What scenario is being tested? (e.g., "time to play Basic Soldier if starting bar = 0")
- What are the initial conditions? (starting bar, deck order, card costs)
- What assumptions are made? (player plays top card immediately when affordable, etc.)

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

### Bar Accrual
Trace bar over time for a given strategy.
- Input: starting bar, fill rate (base + active boosts), elapsed time
- Output: bar at key time points
- Validate: can player afford key cards at intended time?

### Deck Cycle Time
How long until a given card comes back to the top after being played.
- Input: deck size, card costs, bar fill rate
- Output: seconds per full cycle
- Validate: does a 7-card starting deck cycle in a reasonable window?

### Auto-Cycle Timer Impact
When does the stuck-card timer fire relative to bar recovery?
- Input: DeckAutoCycleTimer, card cost, current bar, bar fill rate
- Output: does the timer fire before the player can afford the card?
- Validate: is the lockout timer a help (moves unaffordable cards) or a problem (moves cards just as they become affordable)?

### Combat Duration
Given unit stats, how long does a fight last?
- Input: HP, attack power, attack period for each side, unit count
- Output: time until one side is eliminated, surviving HP
- Validate: matches the design description in CLAUDE.md combat section

### Capital Kill Time
Given a steady stream of soldiers, how long to reduce capital to 0 HP?
- Input: soldiers per minute (derived from bar fill / soldier cost / deck cycle), capital HP
- Output: minimum time to win
- Validate: is this achievable in the ~15-min target game?

### Bar Speed Card ROI
Does a +15% bar speed card pay back its bar cost within its duration?
- Input: base bar fill rate, BarSpeedBonus, BarSpeedDuration, card cost
- Output: net bar gain over duration vs cost paid
- Validate: is it a net positive? By how much?

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
Bar:
  Fill rate: 0.1 bar/sec (fixed — no territory or building bonus)
  Scale: 0–10 (capped)
  +15% Bar Speed card: +0.15/sec bonus for 15s (stacks additively)
    At base rate: +15% = 0.115/sec; 2 active = 0.130/sec
  +1 Bar Burst card: instant +1 bar (0 cost, free to play)

Starting deck (7 cards):
  ×2 Hex Claim: 1 bar each
  ×2 +1 Bar Burst: 0 bar each
  ×2 +15% Bar Speed: 1 bar each
  ×1 Basic Soldier: 2 bar

Soldier:
  Bar cost: 2 | HP: 2 | Attack: 1/sec | Move: 1 hex per 2 sec
  1v1: both die after 2s
  2v1: defender dies 1s; one attacker at 1 HP, one at 2 HP

Capital:
  HP: 20 | Damage per soldier contact: 1
  Soldier vs capital: soldier consumed, capital −1 HP

Deck auto-cycle: 5s (stuck top card moves to bottom after 5s)

Shop:
  Card removal cost: 5 bar (removes top or aside card permanently)
  Shop card pool: TBD (not yet designed)
```

---

## Rules

- **Source of truth:** CLAUDE.md, loaded fresh every invocation. Use quick reference above for orientation only — if it conflicts with CLAUDE.md, CLAUDE.md wins.
- **Show your work:** Every calculation visible, no "and therefore the result is X"
- **No mechanic opinions:** If the problem is the mechanic design (not the number), say "hand off to /game-designer" and stop
- **Precision:** Use exact values. Only round for display if stated explicitly
- **Assumptions visible:** State all behavioral assumptions (e.g., "player plays top card immediately when affordable")
