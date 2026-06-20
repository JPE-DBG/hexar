# UI Layout Spec: Deck-Building Combat
## Hexar — Visual Design & Screen Layout

Generated: 2026-06-20. Updated to reflect M0 implementation: 2026-06-21.

> This document tracks the approved layout. Dimensions reflect what is actually built in M0.
> Color palette superseded by parchment theme — see `asset-theme.md` for authoritative CSS variables.

---

## Design Vocabulary

The new system inverts MVP1's interaction model:

| | MVP1 (Power Gate) | New (Deck-Building) |
|---|---|---|
| Primary interaction | Click hex → click action | Watch top card → click to play |
| Player focus | Map (select and act) | Deck (react and cycle) |
| Information priority | Hex state (power, building) | Card queue + bar level |

The hex map becomes a **theatre** — the player watches units march and buildings appear there. The **deck panel is the cockpit** — that's where decisions are made.

Reference games studied:
- **Clash Royale** — best reference for real-time card + battlefield split
- **mini metro** — information-dense minimalism under time pressure
- **Into the Breach** — grid strategy readable at speed

---

## Color Palette

Parchment theme (approved in `asset-theme.md`). CSS variables defined in `client/src/style.css`:

```css
--bg:           #2c1e0f;   /* dark walnut table */
--chrome:       #f0e6d0;   /* warm parchment */
--chrome-dark:  #e0d0b8;   /* panel subdivisions */
--ui-border:    #b8956a;   /* warm sienna */
--ui-accent:    #c8a020;   /* gold accents */
--text-primary: #2c1810;   /* dark brown ink */
--text-muted:   #8a6840;   /* medium brown */
--bar-fill:     #d4860a;   /* deep amber */
--bar-empty:    #c4aa88;   /* pale tan (now unused — trough PNG handles empty) */
--cost:         #d4860a;   /* amber cost badges */
--danger:       #c03030;
--win:          #c8a020;
```

Hex tile colors come from generated PNG assets, not CSS.

---

## Overall Screen Layout (Desktop)

```
┌────────────────────────────────────────────────────────────────┐
│  SHOP STRIP (230px)                                            │
│  [160×210 card]  [160×210 card]  [160×210 card]  [REMOVE]     │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│                        HEX MAP                                 │
│                    (full width, flex: 1)                       │
│                  pointy-top hex grid                           │
│                                                                │
├────────────────────────────────────────────────────────────────┤
│  BAR METER (32px)                                              │
│  [6]  [████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]                     │
├────────────────────────────────────────────────────────────────┤
│  DECK PANEL (230px)                                            │
│                                                                │
│  ┌──────────────┐   ┌──────────┐   ══  ══  ══                 │
│  │  TOP CARD    │   │  ASIDE   │   (card backs, fanned)        │
│  │  160×210px   │   │ 120×160px│                               │
│  │  [type icon] │   │          │   "12 cards remain"           │
│  │  Card Name   │   │ Card Name│                               │
│  │  Effect text │   │ Effect   │                               │
│  │           [2]│   │       [1]│                               │
│  │  [PLAY ▶]    │   │ [PLAY ▶] │                               │
│  │  [PUSH →]    │   │          │                               │
│  └──────────────┘   └──────────┘                               │
└────────────────────────────────────────────────────────────────┘
```

**Proportions:**
- Shop strip: **230px** (scrollable row of full-size cards)
- Map area: fills remaining height (`flex: 1`)
- Bar meter strip: **32px**
- Deck panel: **230px** (fixed, always visible)

---

## Bar Meter (32px strip, full width)

**Design: integer counter + single fill segment.**

```
┌───────────────────────────────────────────────────────────────┐
│  [6]  [████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] │
└───────────────────────────────────────────────────────────────┘
   ↑                    ↑
integer part         fractional fill (0 → 1)
(payment value)
```

- **Left label:** `floor(bar)` — the integer count of whole bar units. This is the primary payment value the player reads. Bold monospace, 16px, dark brown.
- **Fill segment:** `bar - floor(bar)` as a fraction 0→1. Fills amber left-to-right. When it reaches 1, the integer counter increments and the segment resets.
- **Background:** `bar-meter-trough.png` (stone/wood channel texture), tiled with `background-size: auto 100%; repeat-x`. Segment is transparent when empty so trough shows through.
- **Why this design:** The player's question is always "do I have 2 bar to play a soldier?" — the integer on the left answers this without parsing a segmented strip. The fill segment shows momentum (how close the next unit is).

---

## Deck Panel (bottom, 230px)

### Top Card (leftmost, largest)

Dimensions: **160px × 210px**

```
┌──────────────────────┐
│ [icon]          [2▶] │  ← type icon top-left, cost top-right (amber badge)
│                      │
│   Basic Soldier      │  ← card name, bold, centered
│                      │
│   Deploy a soldier   │  ← effect, small, muted
│   toward capital     │
│                      │
│  [▶ PLAY]  [→ PUSH]  │  ← two action buttons at bottom
└──────────────────────┘
```

- Background: `card-frame.png` stretched 100% × 100%
- Cost badge: amber circle, white number, top-right corner, 22px diameter
- Type icon: emoji (⚔/🏗/⚡), top-left, 14px
- `[▶ PLAY]` — primary action, gold accent
- `[→ PUSH]` — secondary, pushes to aside slot; grayed if aside occupied

### Aside Slot (to the right of top card)

Dimensions: **120px × 160px**

Same layout as top card at reduced scale. Label "ASIDE" above card. `[▶ PLAY]` button only — no push button. If empty: dashed border, "empty" label, muted.

### Card Backs (deck depth indicator)

3 card backs visible to the right of aside slot, fanned (each offset ~8px right and 4px down). Size decreases: 90px → 75px → 60px. Deck count label below.

---

## Card Anatomy (all contexts)

```
┌─────────────────┐
│ [icon]     [cost]│   Row 1: type icon (14px emoji) | cost badge (amber circle, 22px)
│                 │
│   Card Name     │   Row 2: name, bold, 12px, centered
│                 │
│  Effect text    │   Row 3: effect, 10px, muted, centered
│                 │
│  [action btn]   │   Row 4: action buttons (PLAY/PUSH in deck; BUY in shop)
└─────────────────┘
```

**Type icons:** ⚔ unit · 🏗 building · ⚡ bar boost

**In shop context:** BUY button + `×3` quantity badge. Cost badge top-right.

---

## Shop Strip (top, 230px)

Always visible. Horizontal row, scrolls horizontally if cards exceed width.

```
┌──────────────────────────────────────────────────────────────────────────► scroll
│ ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ │⚔  [2]        │  │🏗  [4]        │  │⚡  [3]        │  │✂  [5]        │
│ │ Soldier      │  │ Tower        │  │ Bar Boost    │  │ REMOVE       │
│ │ Deploy…      │  │ Defends…     │  │ +15% rate…   │  │ top · aside  │
│ │ ×3 left      │  │ ×2 left      │  │ ×5 left      │  │              │
│ │  [BUY]       │  │  [BUY]       │  │  [BUY]       │  │              │
│ └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘
└─────────────────────────────────────────────────────────────────────────────────┘
```

- Each card: **160px × 210px** (same size as deck top card)
- Cost badge: top-right corner, 22px amber circle (same as deck cards)
- BUY button at bottom, full card width
- Remove card: cost badge top-right, scissors icon + label centered

---

## Hex Grid

- **Orientation:** pointy-top (vertex at top and bottom, flat edges left and right)
- **Coordinate system:** axial, pointy-top formula in `hexmath.ts`
- **Tile rendering:** flat-top PNG tiles rotated 30° in Canvas context to align with pointy-top clip path
- **Clip path:** procedural pointy-top hexagon per tile
- **Unit token:** filled circle, 28% of hex size, centered; `S` label; player color

---

## Unit Token on Hex Grid

**Design:** Filled circle, 28% of hex diameter, centered on the hex.
- Fill: player color at 90% opacity
- Border: 2px, darkened player color
- Letter inside: `S` for Soldier (white, bold)

---

## Mobile Landscape

Not yet designed for M0. Defer to M1+ polish.

---

## Interaction States (planned for M1+)

| State | Visual |
|---|---|
| Bar insufficient for top card | Card dimmed 50%, cost badge red, PLAY button disabled |
| Aside slot empty | Dashed border, "push here" label |
| Aside slot occupied | Card shown, PUSH button on top card grayed |
| Hex Claim pending | Map highlighted with claimable hexes; prompt in deck panel |
| Capital HP low (<5) | Capital hex ring pulses red |
| Win/loss | Full-screen overlay |

---

## Implementation Notes

### Canvas (`renderer.ts`)
- Hex grid, buildings on hexes, unit tokens
- Tile images: draw with `ctx.rotate(Math.PI / 6)` + `ctx.translate(px, py)` to align flat-top PNG with pointy-top clip
- Unit circle: drawn after `ctx.restore()`, not affected by tile rotation

### DOM overlays
- `#shop-strip` — top 230px, `background: var(--chrome)`, `overflow-x: auto`
- `#bar-meter` — 32px strip, trough PNG background tiled `repeat-x auto 100%`
- `#deck-panel` — bottom 230px, `background: var(--chrome)`
- All panels: parchment chrome with sienna border

### Keyboard shortcuts (planned for M1+)
| Key | Action |
|---|---|
| Space | Play top card |
| E | Push top card to aside slot |
| F | Play aside card |
| ESC | Cancel pending hex claim |
