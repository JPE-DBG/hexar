---
name: game-designer
description: "Use when: analyzing game balance, finding design edge cases, suggesting mechanic changes, evaluating design holes, or refining CLAUDE.md. For pure math/number validation, use /balance-sim instead."
type: skill
---

# Game Designer Skill — Hexar Qualitative Analysis

On-demand skill for **qualitative** design iteration: finding broken interactions, dead strategies, snowball risks, undefined edge cases, and mechanic improvements.

**Not for:** running exact numeric simulations or break-even calculations — use `/balance-sim` for that.

**Prerequisites:** None. Can be invoked anytime.

## How Design Issues Surface in This Project

Design issues in Hexar typically emerge from **playtesting, not upfront analysis**. They arrive as:
- `BUG: [observed behavior] / [expected behavior]` — often contains an implicit design decision
- Rename or relabel requests that reveal a conceptual mismatch ("Economy → Gold")
- Formula corrections discovered after seeing the wrong number in the UI

When a BUG report contains an embedded design decision (e.g., "i dont want to differentiate between placing building and upgrading"), treat it as both a code fix AND a design change: implement the fix, then update CLAUDE.md to reflect the new rule before moving on.

Formal `/game-designer` invocations are for proactive analysis when the user wants to stress-test a mechanic before building it, not for responding to play-discovered issues.

---

## Workflow

### 1. Load Current State (MANDATORY)
- **Always read CLAUDE.md first.** Never rely on cached values from this skill file.
- Extract current mechanics, costs, and rules fresh each invocation.
- Identify which mechanic/system the user wants examined.

### 2. Identify Issues (Qualitative)
- Look for **broken interactions** between systems (e.g., two techs that cancel each other)
- Find **dead strategies** (options no rational player would ever pick)
- Spot **dominant strategies** (options that are always optimal regardless of opponent)
- Check for **undefined behavior** (what happens when X meets Y?)
- Test **degenerate cases** (what if a player does nothing? Rushes one thing only?)

### 3. Evaluate Design Health
For each mechanic, ask:
- Does this create meaningful player decisions? (If one choice is always better, it's fake)
- Does this interact with other systems in interesting ways? (Isolated mechanics add complexity without depth)
- Can the opponent counterplay? (No counterplay = frustrating, not fun)
- Does this fit the 30-minute session target? (Mechanics that matter at minute 45 are dead weight)

### 4. Present Findings
Use the output format below. Always show:
- The specific scenario that causes the problem
- Why it breaks the design intent
- Multiple fix options with trade-offs

### 5. Update CLAUDE.md
- Apply user's chosen fix
- Remove/update all affected cross-references
- Check that the fix doesn't introduce new contradictions

---

## Analysis Prompts (Examples)

- "Find mechanics that don't create meaningful decisions"
- "Which strategies are dominant / which are dead?"
- "What interactions between systems are undefined?"
- "Stress-test [specific mechanic] against degenerate play"
- "Review all victory conditions for achievability and counterplay"
- "What does a player do if they're losing at minute 15?"

---

## Output Format

```
## Issue: [Name]
**Problem:** [What's wrong — qualitative description]
**Scenario:** [Concrete example showing the problem]
**Design intent violated:** [Which design pillar this breaks]
**Impact:** [How this degrades player experience]

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

## Anti-Patterns to Flag

- **Fake choices:** "Pick A or B" where A is always better
- **Complexity without depth:** A rule that's hard to learn but doesn't create interesting play
- **Unfun counterplay:** "The only counter to X is to also do X"
- **Win-more mechanics:** Strong players get stronger, weak players get weaker
- **Unresolvable stalemates:** Two players can lock each other out indefinitely
- **Too-early wins:** Victory conditions achievable before the "interesting phase" of the game
- **Dead features:** Mechanics that never trigger in realistic play

---

## Notes

- All analysis must reference the CURRENT CLAUDE.md state — not remembered values
- Focus on qualitative design health, not numeric precision
- When a finding requires exact math to prove, hand off to `/balance-sim`
- Design changes should be minimal — fix the problem, don't redesign adjacent systems
