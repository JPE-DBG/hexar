---
name: milestone-planner
description: "Use when: breaking Hexar design into buildable milestones, defining MVP scope, sequencing features, or deciding what to cut from first build"
type: skill
---

# Milestone Planner Skill — Hexar Build Sequence

Breaks the current CLAUDE.md design into ordered, testable milestones. Each milestone is shippable, validates a core assumption, and unlocks the next.

**Prerequisites:** Run `/tech-stack` and `/architect` first — milestone sequencing depends on knowing what you're building with and how systems connect.

## Workflow

### 1. Inventory All Features (from CLAUDE.md)
- Read CLAUDE.md fresh. Extract every distinct mechanic as a discrete feature.
- Tag each: `core-loop` / `balance` / `win-condition` / `ui` / `networking` / `polish`
- A "feature" is buildable and testable in isolation (if it's not, break it down further)

### 2. Map Dependencies
For each feature, ask:
- What other features must exist before this one can work?
- What can be stubbed? (e.g., hardcoded map before procedural generation)
- What must be tested early because it's high-risk? (most likely to surface a design flaw)

### 3. Group into Milestones
- Each milestone = 1-2 weeks of work for one developer
- Each milestone ends with something **playable or demonstrably testable**
- Rule: never more than 2 milestones without a playable checkpoint
- Separate "validation milestones" (prove the design works) from "feature milestones" (add content)

### 4. Define Acceptance Criteria
For each milestone:
- **Done when:** concrete observable behavior, not "code is written"
- **Design risk:** what CLAUDE.md assumption might break here?
- **If it breaks:** what's the fallback or redesign path?

### 5. Identify What Gets Cut
- Tag features as MVP-required vs post-MVP
- If total milestone count > 6, something needs to move to post-MVP
- Post-MVP features should not block any MVP milestone

---

## Milestone Design Principles

**Validate risky assumptions first:**
- The economy tick loop is the highest-risk system (does stepped maintenance *feel* right?)
- Combat resolution is second-highest (does the battle timer create tension or frustration?)
- Network sync is a known hard problem but well-understood — defer after local play works

**Local before networked:**
- All game logic works in local 1v1 before adding network layer
- Network bugs are 10x harder to debug than logic bugs
- Design flaws found locally cost 1/10th to fix vs found after networking is built

**Stubbing strategy:**
- Map generation → use a hardcoded test map (fixed layout, known hex positions)
- AI opponent → use a second player on same machine, or simple scripted behavior
- UI polish → placeholder rectangles/text until gameplay validates

---

## Output Format

```
## Feature Inventory
[Table of features extracted from CLAUDE.md, tagged by category]

## Dependency Graph
[Which features depend on which — text format, not visual]

## Milestone Plan

### M[N] — [Name] (~[time estimate])
**Goal:** [What assumption/system this validates]
**Features:** [List from inventory]
**Stubbed:** [What's faked in this milestone]
**Done when:** [Observable acceptance criteria]
**Design risk:** [What might break and what to do if it does]

...repeat for each milestone...

## Post-MVP (deferred)
[Features that don't make it into the initial build sequence]
```

---

## Rules

- **Never cache milestones.** Generate fresh from current CLAUDE.md each invocation.
- **No milestone without a "done when."** If you can't define acceptance criteria, the milestone is too vague.
- **Shortest path to playable.** The first playable checkpoint should be ≤ 3 milestones in.
- **Mark risks honestly.** Every milestone should name what could go wrong.
