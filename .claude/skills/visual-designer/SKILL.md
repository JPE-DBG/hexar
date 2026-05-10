---
name: visual-designer
description: "Use when: choosing a visual style, color palette, or UI layout for Hexar; evaluating what hex strategy games do visually; designing the build menu, lobby screen, or effect system; deciding what animations reinforce game mechanics vs add noise."
type: skill
---

# Visual Designer Skill — Hexar UI & Aesthetics

Guides visual design decisions for Hexar: color palette, typography, hex rendering style, UI component layout, and animation priorities. Bridges game feel and technical feasibility.

**Prerequisites:** Read CLAUDE.md for current mechanics before evaluating visuals — effects should reinforce what matters in the game, not just look good in isolation.

---

## Workflow

### 1. Establish Visual Identity
Define the tone before designing components. Hexar is:
- Strategic, not casual (avoid bright cartoon aesthetics)
- Fast-paced real-time (animations must be readable at speed, not distracting)
- 2-player competitive (player differentiation is critical — colors must be immediately distinct)

Reference games to study (visual patterns, not mechanics):
- **Polytopia** — clean flat hex style, readable at a glance
- **Antiyoy** — minimal but legible, good contrast model
- **Catan Universe** — richer hex textures, shows what "premium" hex looks like
- **mini metro** — excellent information-dense minimalism
- **Into the Breach** — grid strategy UI that reads clearly under time pressure

### 2. Color Palette First
All visual decisions flow from the palette. Define before touching any code:
- **Background:** dark neutral (not pure black — use deep navy or slate)
- **Player 1 hex:** current teal `#4ecdc4` — evaluate if it holds contrast against background
- **Player 2 hex:** current red `#ff6b6b` — check accessibility (red/green colorblind safe?)
- **Unclaimed hex:** muted, clearly subordinate to owned hexes
- **UI chrome:** dark panels, semi-transparent overlays
- **Accent / action:** single highlight color for interactive elements
- **Status colors:** battle (amber), pause (purple), victory (gold), danger (red)

Export as CSS variables immediately. Every hardcoded color string is a debt.

### 3. Hex Rendering Hierarchy
Every hex must communicate its state at a glance, in priority order:
1. **Ownership** — player color fills, must be unmistakable
2. **Power** — label or badge, readable at hex scale
3. **Building type** — icon or shape indicator
4. **Status** — battle ring, fortify border, selection highlight

Do not add visual complexity that obscures this hierarchy. Test by squinting — if ownership is unclear, the palette failed.

### 4. Animation Priority (by player impact)
Score each animation by: does the player *need* this feedback to make decisions?

| Effect | Gameplay value | Priority |
|---|---|---|
| Capture flash on hex takeover | High — confirms action resolved | Must have |
| Battle ring pulse | High — signals ongoing fight | Must have |
| Floating `+gold` on income | Medium — reinforces economy loop | Should have |
| Fortify border segments | High — existing, keep | Keep |
| Hex selection highlight | High — existing, keep | Keep |
| Screen shake on capital fall | Low — pure feel | Nice to have |
| Idle hex shimmer | None | Skip |

Rule: if removing an animation doesn't affect decision-making, it's decoration. Decoration gets cut when time is tight.

### 5. UI Component Evaluation
For each UI panel (build menu, tech tree, HUD, lobby, overlays):
- What is the player's job here? (Buy something? Read info? Make a choice?)
- What's the primary action? Make it the largest touch target.
- What information is secondary? Move it to a secondary zone, smaller.
- What can be removed entirely?

Build menu priority order: attack → upgrade → demolish → fortify (matches frequency of use)

### 6. Typography
- Monospace for all numeric values (gold, power, timers) — prevents layout shift
- Sans-serif for labels — readable at small sizes
- No more than 2 font sizes in any one panel
- Minimum 12px for anything a player reads under pressure

---

## Evaluation Checklist

### Readability
- [ ] Can you tell who owns each hex from 1 meter away from the screen?
- [ ] Is the battle timer readable during a contested fight?
- [ ] Does the selected hex stand out clearly from adjacent owned hexes?
- [ ] Is the gold counter legible at 60fps when it's changing fast?

### Visual Hierarchy
- [ ] Does the most important information have the most visual weight?
- [ ] Are interactive elements visually distinct from informational elements?
- [ ] Do overlays (pause, victory, waiting) clearly sit "above" the game board?

### Consistency
- [ ] Are all player-1 elements the same hue?
- [ ] Are all timer animations the same direction (clockwise) and speed model?
- [ ] Do all action buttons look the same regardless of which menu they're in?

### Performance
- [ ] Do animations run at 60fps without affecting the 100ms game-state tick?
- [ ] Are gradients and shadows within Canvas/CSS GPU budget?
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

## Palette Proposal
Background:   #______
Player 1:     #______
Player 2:     #______
Unclaimed:    #______
UI chrome:    #______
Accent:       #______
[Status colors...]
```

---

## Rules

- **Serve the game, not the portfolio.** Every visual choice must help players make faster, better decisions — or get cut.
- **Contrast is non-negotiable.** Player colors must pass WCAG AA at hex scale. Test with a colorblind simulator.
- **Animation budget is shared.** The 60fps animation loop and 100ms game-state loop must not fight each other. Never block the game loop for a cosmetic effect.
- **Design systems, not one-offs.** A color chosen for one overlay will be reused elsewhere. Establish variables before writing any values.
- **Reference before inventing.** Hex strategy games have solved most layout problems. Look at what works before designing from scratch.
