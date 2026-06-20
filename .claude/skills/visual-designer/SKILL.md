---
name: visual-designer
description: "Use when: choosing a visual style, color palette, or UI layout for Hexar; evaluating what hex strategy games do visually; designing the deck panel, shop panel, bar meter, or unit token rendering; deciding what animations reinforce game mechanics vs add noise."
type: skill
---

# Visual Designer Skill — Hexar UI & Aesthetics

Guides visual design decisions for Hexar: color palette, typography, hex rendering style, UI component layout, and animation priorities. Bridges game feel and technical feasibility.

**Prerequisites:** Read CLAUDE.md and `.claude/design/ui-layout-deck-building.md` before evaluating visuals — the layout doc is a draft reference, not a finalized spec.

**Important:** Visual decisions in this skill feed directly into AI image generation via ComfyUI. Style choices must be expressible as text-to-image prompt keywords. Avoid decisions that can only be described as "hand-drawn" or that require manual illustration — the asset pipeline is generative.

---

## Workflow

### 1. Establish Visual Identity
Define the tone before designing components. Hexar is:
- Strategic, not casual (avoid bright cartoon aesthetics)
- Fast-paced real-time (animations must be readable at speed, not distracting)
- 2-player competitive (player differentiation is critical — colors must be immediately distinct)

Reference games to study (visual patterns, not mechanics):
- **Clash Royale** — best reference for real-time card + battlefield split
- **mini metro** — information-dense minimalism under time pressure
- **Into the Breach** — grid strategy UI that reads clearly under time pressure
- **Polytopia** — clean flat hex style, readable at a glance

### 2. Color Palette (DRAFT — Not Approved)
The palette in `ui-layout-deck-building.md` is an **initial draft only** — not reviewed or approved. Before implementing anything, this section must produce a finalized palette decision.

Key questions still open:
- Overall tone: dark navy? earthy? high-contrast neon? muted tactical?
- Player differentiation: teal vs red is one option — are these the right hues?
- Accent color: amber/gold for bar and cost — does this fit the chosen visual style?
- Hex tile colors: unclaimed slate, owned in player color — right approach?

When the palette is finalized, document it as CSS variables and save to `.claude/design/asset-theme.md`. Every hardcoded color string before that point is throwaway code.

### 3. Screen Layout (DRAFT — Not Approved)

`ui-layout-deck-building.md` contains an initial layout proposal. It is a starting point for discussion, not a spec to implement from. Use this skill to evaluate and finalize layout decisions before any client code is written.

### 4. Information Hierarchy for New UI Elements

Every UI element must communicate state at a glance. Priority order for deck panel:
1. **Top card identity** — card type icon + name (most space)
2. **Bar cost** — amber badge top-right, immediately visible
3. **Affordability state** — dimmed card + red badge when can't afford
4. **Action buttons** — PLAY + PUSH, clearly labeled with hotkeys
5. **Aside slot** — dashed border when empty, card shown when occupied

Do not add visual complexity that obscures this hierarchy. Test by squinting — if you can't tell what the top card costs, the layout failed.

### 5. Animation Priority (by player impact)

| Effect | Gameplay value | Priority |
|---|---|---|
| Card flip/cycle when played | High — confirms action | Must have |
| Unit march (hex-to-hex slide) | High — shows army progress | Must have |
| Unit combat pulse ring | High — signals fight | Must have |
| Bar fill animation | Medium — shows resource state | Must have |
| Capital HP low warning pulse | High — danger signal | Must have |
| Card unavailable flash (insufficient bar) | Medium — prevents confusion | Should have |
| Bar full amber pulse (cap 10/10) | Medium — spend signal | Should have |
| Victory/loss overlay | High — existing, keep | Keep |
| Idle card shimmer | None | Skip |

Rule: if removing an animation doesn't affect decision-making, it's decoration. Decoration gets cut when time is tight.

### 6. Typography
- Monospace for all numeric values (bar amount, HP, card cost, timer) — prevents layout shift
- Sans-serif for labels — readable at small sizes
- No more than 2 font sizes in any one panel
- Minimum 12px for anything a player reads under pressure

---

## Evaluation Checklist

### Readability
- [ ] Can you tell who owns each hex from 1 meter away from the screen?
- [ ] Is the current bar amount readable at a glance (segmented bar)?
- [ ] Is the top card's cost badge immediately visible?
- [ ] Is it obvious when a card can't be played (dimmed state)?
- [ ] Are unit HP dots readable at hex scale?

### Visual Hierarchy
- [ ] Does the deck panel (cockpit) draw the eye more than the hex map (theatre)?
- [ ] Are interactive elements visually distinct from informational elements?
- [ ] Do overlays (pause, victory, waiting) clearly sit "above" the game board?
- [ ] Is the aside slot's empty/occupied state immediately obvious?

### Consistency
- [ ] Are all player-1 elements the same hue?
- [ ] Do all card types use their correct background colors (unit/building/bar)?
- [ ] Does the cost badge always use the same amber color?

### Performance
- [ ] Do animations run at 60fps without affecting the 100ms game-state tick?
- [ ] Does the unit march animation use pure visual interpolation (not blocking game state)?
- [ ] Does anything animate that doesn't need to? (Idle elements should be static)

---

## Output Format

```
## Visual Decision: [Topic]
**Goal:** [What player experience this serves]
**Recommendation:** [Specific choice with reasoning]
**Reference:** [Which game or pattern this borrows from]
**Implementation note:** [Canvas, CSS, or DOM approach]
**What to avoid:** [Common mistake for this element]
```

---

## Rules

- **Serve the game, not the portfolio.** Every visual choice must help players make faster, better decisions — or get cut.
- **Contrast is non-negotiable.** Player colors must pass WCAG AA at hex scale.
- **Animation budget is shared.** The 60fps animation loop and 100ms game-state loop must not fight each other. Never block the game loop for a cosmetic effect.
- **Design systems, not one-offs.** A color chosen for one card will be reused elsewhere. Establish variables before writing any values.
- **The layout doc is a draft.** `ui-layout-deck-building.md` is an initial proposal — use this skill to evaluate and finalize it, not implement from it blindly.
- **Nothing is approved until stated.** Color palette, art style, layout proportions — all are open questions until explicitly signed off.

---

## Required Output — Asset Theme Doc

When this skill produces finalized visual decisions, write them to `.claude/design/asset-theme.md`. This file is the handoff to `/asset-gen` and must include:

1. **Art style** — 2–3 sentences + 5–8 ComfyUI-compatible style keywords (e.g. `flat vector, limited palette, top-down, crisp edges`)
2. **Negative keywords** — what to exclude from all generations (e.g. `photorealistic, 3d render, gradients`)
3. **Color palette** — finalized hex codes as CSS variables, with role labels
4. **Per-asset direction** — one line per asset type (unclaimed hex, p1 hex, capital, soldier) describing the visual intent
5. **Status** — mark each decision as `approved` or `draft`

Do not write to `asset-theme.md` until decisions are actually confirmed in conversation.
