---
name: visual-designer
description: "Use when: choosing a visual style, color palette, or UI layout for Hexar; evaluating what hex strategy games do visually; designing the deck panel, shop panel, bar meter, or unit token rendering; deciding what animations reinforce game mechanics vs add noise."
type: skill
---

# Visual Designer Skill — Hexar UI & Aesthetics

Guides visual design decisions for Hexar: color palette, typography, hex rendering style, UI component layout, and animation priorities. Bridges game feel and technical feasibility.

**Prerequisites:** Read CLAUDE.md and `.claude/design/ui-layout-deck-building.md` before evaluating visuals — the layout spec is the reference for all panel decisions.

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

### 2. Color Palette (Defined — Verify Against Design Doc)
All visual decisions flow from the palette. Defined in `ui-layout-deck-building.md`:

```css
--bg: #1a1a2e       /* deep navy background */
--chrome: #252545   /* UI panel surfaces */
--ui-border: #3a3a5c
--p1: #4ecdc4       /* Player 1 teal */
--p2: #ff6b6b       /* Player 2 red */
--unclaimed: #2d3748
--bar-fill: #f7c948 /* amber — the accent color */
--bar-empty: #3a3a5c
--card-unit: #2a3f6f
--card-building: #3d2e1e
--card-bar: #1e3d2e
--cost: #f7c948
--text-primary: #e8e8f0
--text-muted: #6b7280
--danger: #e53e3e
--win: #f6d860
```

Export as CSS variables immediately. Every hardcoded color string is a debt.

### 3. Screen Layout (Defined — See Design Doc)

The layout is specified in `ui-layout-deck-building.md`:
- **HUD strip:** 40px top — deck count, timer, opponent bar
- **Map + shop area:** fills remaining height
- **Shop sidebar:** 200px right, always visible
- **Bar meter strip:** 24px above deck panel
- **Deck panel:** 190px bottom, always visible

Do not redesign the layout — implement from the spec. Use this skill for:
- Component-level visual decisions not covered by the spec
- Mobile adaptation questions
- Animation decisions

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
- **The layout spec is the source.** For panel dimensions and positioning, `ui-layout-deck-building.md` wins. Use this skill for the decisions that spec doesn't cover.
