---
name: balance-sim
description: "Use when: running economy projections, simulating upgrade cost payoffs, tracing TP timelines, or verifying balance numbers from CLAUDE.md before they get hardcoded"
type: skill
---

# Balance Sim Skill — Hexar Numeric Validation

Runs traced simulations of Hexar economy, combat, and tech timelines using exact CLAUDE.md numbers. Finds numeric imbalances before they become bugs in running code.

## Workflow

### 1. Load Current Numbers
- Read CLAUDE.md for all constants: base income, maintenance tiers, building costs, upgrade curves, TP rates, tech costs
- Note any numbers marked as estimates or unverified

### 2. Run Requested Simulation
- Trace step-by-step with explicit math at each step
- Show intermediate values (not just final result)
- Flag where assumptions were made

### 3. Compare Against Design Intent
- Does the result match the CLAUDE.md design goal?
- Is a victory condition reachable in the intended time window?
- Does a strategy feel rewarding or punishing in the right ways?

### 4. Identify Numeric Issues
- Values that make a strategy dominant or useless
- Thresholds that are never reached in practice
- Costs that pay off too fast or too slowly

### 5. Propose Adjustments
- Suggest specific number changes (not mechanic redesigns)
- Show the simulation result after the proposed change
- Let game-designer skill handle mechanic-level redesigns

---

## Common Simulations

### Economy Projection
Trace gold over time for a given strategy:
```
Input: hex expansion rate, building timing, tech unlocks
Output: gold/sec at T=2min, T=5min, T=10min, T=20min
Validate: can player save 100g for first enemy attack at the right time?
```

### Upgrade Payoff (Break-Even)
How many seconds until an upgrade pays for itself:
```
Economy L1: costs 40g, gains +0.75/sec (at base ×1.5 multiplier)
Break-even: 40 / 0.75 = 53 seconds
Question: is 53 seconds a good investment in a 30-min game? (Yes — pays off 27x)
```

### Maintenance Cap Trace
At what hex count does income go negative without Economy buildings:
```
Trace: 1 hex → 2 hex → ... → 31 hex
Show net income at each step using stepped maintenance tiers
Validate: CLAUDE.md claims auto-drop triggers at 31 hexes — confirm
```

### TP Timeline
How long to reach Tech Level 4 with different Research investment:
```
Input: number of Research buildings, upgrade levels, start time
Output: time to unlock each tech, time to Tech Level 4
Validate: Tech Level 4 + 35% map reachable in 15-18 min (CLAUDE.md claim)
```

### Battle Outcome Table
Given two Power values, show all outcomes:
```
Input: attacker Power, defender Power
Output: result (fail / standard battle duration / instant), counter-spend options remaining
```

---

## Output Format

```
## Simulation: [Name]
**Setup:** [Initial conditions]
**Constants used:** [Which CLAUDE.md values]

Step-by-step trace:
  T=Xs: [state] → [calculation] → [new state]
  ...

**Result:** [Final value]
**Design target:** [What CLAUDE.md says it should be]
**Match:** YES / NO — [explanation if no]

## Finding (if any)
**Issue:** [What the numbers reveal]
**Suggested fix:** [Specific number change]
**Revised result:** [Simulation after fix]
```

---

## Notes

- Always show the math, not just conclusions
- Use CLAUDE.md numbers exactly — don't round unless stated
- Economy formula: `(base + 0.5 × level) × 1.5` for Economy buildings
- Maintenance tiers: 1/sec (hexes 1-10), 2/sec (11-20), 3/sec (21+)
- Battle duration: `5 + (AttackerPower + DefenderPower) / 2` seconds
- Counter-spend cap: `min(+3, seconds_remaining)`
- All 4 tech costs: Iron Grip 50, Production Boom 40, Efficient Conquest 35, Garrison 60 = 185 TP total
