---
name: milestone-planner
description: "Use when: planning the next feature or post-MVP addition to Hexar, sequencing work, estimating scope, or deciding what to build vs defer"
type: skill
---

# Milestone Planner Skill — Hexar Feature Planning

Plans and sequences Hexar features during the active deck-building combat redesign. MVP1 (Power Gate) is frozen on `release_mvp_1`. Active work is on `mvp_fix_combat_2026_06_13`.

**Prerequisites:** Read CLAUDE.md "Current Status" and "Open Questions" sections first.

---

## Context: Current State

- **MVP1** complete and frozen on `release_mvp_1` branch (deployed to Fly.io)
- **Active branch** `mvp_fix_combat_2026_06_13`: deck-building combat redesign in progress
- Stack: Go server + TypeScript client, WebSocket, Fly.io deployment, GitHub Actions CI/CD
- Tests: Go unit tests (`internal/game/`), room integration (`internal/room/`), Playwright E2E (`e2e/`) — game tests not yet written for new system
- All decisions and architecture are documented in CLAUDE.md

**Implementation progress (check CLAUDE.md Current Status for latest):**
- Design spec, CLAUDE.md rewrite, UI layout spec: done
- Old game logic deleted, new Go structs stubbed, compiles: done
- Core backend (bar, deck, units, win), shop, client: pending

**Source of truth:** CLAUDE.md (this file). Do not create a separate MILESTONES.md.

---

## Workflow

### 1. Load Context from CLAUDE.md

- Read "Current Status" implementation table — what's done, what's pending
- Read "Open Questions" — what's TBD that might block a feature
- Read "Game Rules" for the mechanic context of any feature being added
- Check "Rejected Alternatives" — don't re-litigate settled decisions

### 2. Define the Feature Clearly

Before estimating or sequencing:
- What is the observable player-facing change?
- What new server state is needed (if any)?
- What new client state/rendering is needed?
- Does this need a new test file or extend an existing one?
- Does CLAUDE.md need updating first? (Policy: update CLAUDE.md before implementing)

### 3. Map Dependencies

- What existing system does this touch? (`internal/game/`, `internal/room/`, `client/src/ui/`, etc.)
- What must exist before this works?
- What can be stubbed for a first pass?
- What's the testing seam? (Can it be unit tested? Integration? E2E?)

### 4. Estimate and Sequence

Group by milestone (1–2 weeks of solo work):
- Each milestone ends with something **playable or demonstrably testable**
- Validate risky assumptions early (design, balance) before polish
- Playtesting feedback often invalidates assumptions — keep milestones short

### 5. Define Acceptance Criteria

- **Done when:** concrete observable behavior (not "code is written")
- **Balance risk:** what CLAUDE.md assumption might this break?
- **Test coverage:** which new tests are required?

---

## Remaining Implementation Sequence

The implementation plan from CLAUDE.md, in dependency order:

1. **Core backend** — bar fill, deck cycling, unit march, combat, capital damage, win condition
   - `bar.go`, `deck.go`, `units.go`, `victory.go` already stubbed; needs full implementation + tests
   - Done when: `go test ./internal/game/...` passes with all test cases from test-writer skill

2. **Shop card pool design** — which cards, costs, quantities (design session, not code)
   - Can run in parallel with core backend implementation
   - Output: CLAUDE.md shop card pool section filled in
   - Blocks shop backend implementation

3. **Shop backend** — shared pool, BuyCard, card removal
   - Blocked by: shop card pool design
   - Done when: `shop_test.go` passes

4. **Client rewrite** — state.ts, bar meter, deck panel, shop panel, unit token rendering
   - Blocked by: core backend (needs DTOs to build against)
   - Reference: `.claude/design/ui-layout-deck-building.md` for layout
   - Done when: two browser tabs can claim hexes, deploy soldiers, and capital takes damage

5. **Playtest core loop** — validate bar rate, deck cycle feel, soldier march speed
   - Blocked by: client rewrite
   - Output: adjusted constants in CLAUDE.md and constants.go

6. **Shop integration** — add shop cards one at a time, playtest each
   - Blocked by: playtest core loop + shop backend

---

## Milestone Design Principles

**Short milestones beat long ones.** Balance flaws found at week 1 are 10x cheaper than at week 4.

**Update CLAUDE.md before implementing.** Design spec first, then code, then tests.

**Post-MVP sequencing priorities:**
1. Core loop playability first (bar rate, soldier feel, capital HP)
2. Shop card pool design — needed before shop can be built
3. Shop implementation
4. Polish and additional card types

---

## Output Format

```
## Feature: [Name]
**What it is:** [One-sentence player-facing description]
**CLAUDE.md section to update:** [Where to write the spec first]
**Files affected:** [List]
**Testing approach:** [Unit / Integration / E2E — which layer and why]
**Done when:** [Observable acceptance criteria]
**Balance risk:** [What could unexpectedly break]

## Sequence Recommendation
[Ordered list of features if multiple were requested]
**Why this order:** [Dependencies or risk that drives the sequence]

## Deferred (and why)
[Features that should wait]
```

---

## Rules

- **No milestone without "done when."** If you can't define acceptance criteria, the feature is too vague.
- **Mark risks honestly.** Every feature should name what could go wrong.
- **CLAUDE.md first.** If a feature changes game rules, write the spec in CLAUDE.md before writing code.
- **One test layer per feature.** Don't add Playwright for something fully covered by Go unit tests.
- **Check open questions.** If a TBD in CLAUDE.md blocks this feature, resolve it first.
