---
name: milestone-planner
description: "Use when: planning the next feature or post-MVP addition to Hexar, sequencing work, estimating scope, or deciding what to build vs defer"
type: skill
---

# Milestone Planner Skill — Hexar Feature Planning

Plans and sequences future Hexar features. M1–M9 are complete and deployed. This skill is for evaluating and ordering post-MVP work.

**Prerequisites:** Read CLAUDE.md "Open Questions & Future Work" and "Current Status" sections first.

---

## Context: Current State

- M1–M9 complete and deployed to Fly.io
- Stack: Go server + TypeScript client, WebSocket, Fly.io deployment, GitHub Actions CI/CD
- Tests: 3 layers — Go unit tests (`internal/game/`), room integration (`internal/room/`), Playwright E2E (`e2e/`)
- All decisions and architecture are documented in CLAUDE.md

**Source of truth:** CLAUDE.md (this file). Do not create a separate MILESTONES.md — it duplicates content and drifts. Record milestone info in git commit messages if tracking is valuable.

---

## Workflow

### 1. Load Context from CLAUDE.md

- Read "Open Questions & Future Work" — what's already identified as next
- Read "Known Issues" — are there blockers that need fixing before new features?
- Read "Balance Rules" and "Game Rules" for the mechanic context of any feature being added
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

## Milestone Design Principles

**Short milestones beat long ones.** Balance flaws found at week 1 are 10x cheaper than at week 4.

**Update CLAUDE.md before implementing.** Design spec first, then code, then tests. Tests validate the spec, not the implementation.

**Don't re-plan completed work.** M1–M9 history is in git. Focus on what's next.

**Post-MVP sequencing priorities:**
1. Balance/playtesting gaps first (these can invalidate subsequent features)
2. Known issues that block player experience
3. Tech debt that creates friction for future features
4. New mechanics that expand archetypes

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
