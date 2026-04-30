---
name: game-designer
description: "Use when: analyzing game balance, finding design edge cases, validating economy math, simulating scenarios, or refining CLAUDE.md mechanics"
type: skill
---

# Game Designer Skill — Hexar Refinement

On-demand skill for iterating on the Hexar game design. Helps identify balance issues, edge cases, economy math validation, and suggests refinements to CLAUDE.md.

## Workflow

### 1. **Analyze Current Design**
- Read CLAUDE.md to understand current state
- Identify which mechanic/system the user wants to examine
- Run preliminary balance checks

### 2. **Find Edge Cases & Issues**
- Simulate specific scenarios (e.g., "What if player rushes Economy early?")
- Check for snowballing risks
- Identify unfair advantages or dead strategies
- Test boundary conditions (min/max values)

### 3. **Validate Economy Math**
- Trace resource flow through scenarios
- Check early-game parity
- Verify maintenance doesn't break at different hex counts
- Ensure victory conditions are reachable

### 4. **Present Findings**
- List specific issues with examples
- Show the math or scenario that causes it
- Suggest 2-3 concrete fixes (with trade-offs)
- Highlight which systems interact unexpectedly

### 5. **Update CLAUDE.md**
- Apply user's chosen fix
- Update affected sections (balance rules, economy examples, open questions)
- Note rationale in commit message

---

## Common Refinement Tasks

**Ask me to:**
- ✅ "Analyze early-game economy balance"
- ✅ "Find edge cases in combat Power rules"
- ✅ "Validate victory condition thresholds (70% + 20s reachable?)"
- ✅ "Simulate tech rush vs military rush strategies"
- ✅ "Check if maintenance system prevents/allows snowballing"
- ✅ "Propose new tech tree entry and balance it"
- ✅ "Identify dominant strategies that need rebalancing"
- ✅ "Design test scenarios for playtesting"

**Examples:**
- "Is 50 gold/sec counter-spend too cheap?"
- "What happens if player ignores map control and only techs?"
- "Can a 5-hex player win vs a 10-hex player?"
- "Should Blitzkrieg tech cost more than 80 TP?"

---

## Key Mechanics to Stress-Test

### Economy
- Resource generation per hex (2/sec base)
- Maintenance scaling (1/sec per hex)
- Building upgrade costs and payoff
- Income cap with different hex counts

### Combat
- Power difference rules (instant takeover at >+3)
- Battle duration formula (5 + (A+D)/2)
- Attacker damage (loses levels based on diff)
- Counter-spend cost (50 gold/sec)

### Tech Tree
- Research generation (0.1 TP/sec per level)
- Tech cost scaling (should it be exponential?)
- Which techs are must-haves vs situational
- Tech synergies and combos

### Victory
- 70% threshold reachable in 30 min?
- 20-second hold time creates endgame tension?
- Tech Level 5 + 50% map achievable?
- Time limit (30 min) realistic?

---

## Output Format

When analyzing, provide:

```
## Issue: [Name]
**Problem:** [What's wrong]
**Example:** [Concrete scenario showing the problem]
**Math:** [Numbers that prove it]
**Impact:** [How this breaks balance]

## Proposed Fix
Option A: [Change description]
- Pro: [benefit]
- Con: [tradeoff]

Option B: [Alternative]
- Pro: ...
- Con: ...

**Recommendation:** Option [X] because [why]
```

---

## Notes

- All analysis refs current CLAUDE.md state
- Economy math assumes tick-based simulation (100ms or faster)
- Victory conditions tested with 30-min target playtime
- Assumes 1v1 start (balance for 2v2 later)

