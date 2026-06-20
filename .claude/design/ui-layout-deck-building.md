# UI Layout Spec: Deck-Building Combat
## Hexar — Visual Design & Screen Layout

Generated: 2026-06-20.

> This document is the required output of the graphical design step.
> Client code must not be written until this spec is reviewed.

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

```
Background:       #1a1a2e   (deep navy — not pure black, reduces eye strain)
UI chrome:        #252545   (slightly lighter than background)
UI border:        #3a3a5c   (panel outlines)

Player 1 hex:     #4ecdc4   (teal — existing, keep)
Player 2 hex:     #ff6b6b   (red — existing, keep)
Unclaimed hex:    #2d3748   (dark slate — clearly subordinate)
Capital hex:      gold ring on owner color

Bar fill:         #f7c948   (amber/gold — accent color, high visibility)
Bar empty:        #3a3a5c   (same as UI border — recedes)

Unit token P1:    #4ecdc4 with #1a5f5a border
Unit token P2:    #ff6b6b with #7a1f1f border

Card — Unit:      #2a3f6f   (blue-slate background)
Card — Building:  #3d2e1e   (dark brown background)
Card — Bar boost: #1e3d2e   (dark green background)
Card — cost badge:#f7c948   (amber — same as bar fill, reinforces "this costs bar")

Status — battle:  #f7941d   (amber pulse)
Status — danger:  #e53e3e   (red — capital low HP warning)
Status — win:     #f6d860   (gold)
```

Export all as CSS variables: `--bg`, `--chrome`, `--p1`, `--p2`, `--unclaimed`, `--bar-fill`, `--bar-empty`, `--card-unit`, `--card-building`, `--card-bar`, `--cost`.

---

## Overall Screen Layout (Desktop 1280×800+)

```
┌────────────────────────────────────────────────────────────────┐
│  SHOP PANEL (60px)                                             │
│  [⚔ Soldier ×3 BUY]  [🏗 Tower ×2 BUY]  [⚡ Bar Boost ×5 BUY] │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│                        HEX MAP                                 │
│                    (full width)                                │
│                                                                │
│                  units march here                              │
│                  buildings on hexes                            │
│                                                                │
├────────────────────────────────────────────────────────────────┤
│  BAR METER (24px)                                              │
│  [▮▮▮▮▮▮░░░░]  6.0 / 10                                       │
├────────────────────────────────────────────────────────────────┤
│  DECK PANEL (190px)                                            │
│                                                                │
│  ┌──────────────┐   ┌──────────┐   ══  ══  ══                 │
│  │  TOP CARD    │   │  ASIDE   │   (3 card backs, fanned)     │
│  │              │   │  SLOT    │                               │
│  │  [type icon] │   │          │   "12 cards remain"          │
│  │  Card Name   │   │ Card Name│                               │
│  │  Effect text │   │ Effect   │                               │
│  │           [2]│   │       [1]│                               │
│  │  [PLAY ▶]    │   │ [PLAY ▶] │                               │
│  │  [PUSH →]    │   │          │                               │
│  └──────────────┘   └──────────┘                               │
└────────────────────────────────────────────────────────────────┘
```

**Proportions:**
- Shop strip: 60px (horizontal scrollable row of compact cards)
- Map area: fills remaining height minus shop, bar, and deck panel
- Bar meter strip: 24px
- Deck panel: 190px (fixed, always visible)

---

## HUD Strip (top, 40px)

Left side — your info:
- Deck card count: `⬛ 12` (small card icon + number)
- Aside slot indicator: filled/empty dot

Center:
- Game timer (MM:SS, monospace)

Right side — opponent info:
- Their bar meter: small 80px-wide segmented bar in their color
- Label: `OPP 6/10`

---

## Bar Meter (24px strip, full width)

**Segmented, 10 divisions.** Each segment = 1 bar. Fills left to right.

- Filled segment: `#f7c948` (amber)
- Empty segment: `#3a3a5c` (dark)
- Segment gap: 2px dark gap between each segment
- Text label right-aligned: `6.0 / 10` (monospace)
- At cap (10/10): entire bar pulses amber once to signal "full, spend it"

**Why segmented over smooth:** The player's question is always "do I have 2 bar for a soldier?" — 10 visible blocks answer this instantly without reading the number.

---

## Deck Panel (bottom, 190px)

### Top Card (leftmost, largest)

Dimensions: **140px × 170px**

```
┌──────────────────────┐
│ [icon]          [2▶] │  ← type icon top-left, cost top-right (amber badge)
│                      │
│   Basic Soldier      │  ← card name, bold, centered
│                      │
│   Deploy a soldier   │  ← effect, 1 line, small
│   toward capital     │
│                      │
│  [▶ PLAY]  [→ PUSH]  │  ← two action buttons at bottom
└──────────────────────┘
```

- Background color by card type (unit/building/bar)
- Cost badge: amber `#f7c948` circle, white number inside, top-right corner
- `[▶ PLAY]` button: primary action, highlighted (accent color)
- `[→ PUSH]` button: secondary, pushes to aside slot; grayed if aside slot is occupied
- Hotkey labels shown on buttons: `[SPACE]` / `[E]`
- Disabled state (insufficient bar): card dimmed, play button grayed, cost badge red

### Aside Slot (to the right of top card)

Dimensions: **110px × 140px**

Same layout as top card at reduced scale. Label "ASIDE" in small muted text above card. `[▶ PLAY]` button only — no push button. If empty: dashed border, "empty" label, muted.

### Card Backs (deck depth indicator)

3 card backs visible to the right of the aside slot, fanned slightly (each offset ~8px right and 4px down). Size decreases: 90px → 75px → 60px. Shows deck depth without revealing information.

Below the fan: `"12 cards"` in muted small text.

---

## Card Anatomy (all contexts)

```
┌─────────────────┐
│ [icon]     [cost]│   Row 1: type icon (16px) | cost badge (amber circle, 20px)
│                 │
│   Card Name     │   Row 2: name, bold, 14px, centered
│                 │
│  Effect text    │   Row 3: effect, 11px, muted, centered, max 2 lines
│                 │
│  [action btn]   │   Row 4: action buttons (only in deck panel)
└─────────────────┘
```

**Type icons (16px):**
- Unit card: ⚔ (sword)
- Building card: 🏗 (or simple square outline)
- Bar boost card: ⚡ (lightning)

**Background by type:**
- Unit: `#2a3f6f` blue-slate
- Building: `#3d2e1e` dark brown
- Bar boost: `#1e3d2e` dark green

**In shop context:** replace action buttons with `[BUY]` + quantity badge `×3`.

---

## Shop Panel (top strip, 60px)

Always visible — not collapsible. Horizontal row across the full width. Scrolls horizontally if cards exceed screen width.

```
┌──────────────────────────────────────────────────────────────────────────► scroll
│ ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ │⚔ Soldier  ×3│  │🏗 Tower    ×2│  │⚡ Bar Boost ×5│  │  REMOVE [5▶] │
│ │  cost: 2     │  │  cost: 4     │  │  cost: 3     │  │  top | aside │
│ │  [BUY]       │  │  [BUY]       │  │  [BUY]       │  │              │
│ └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘
└─────────────────────────────────────────────────────────────────────────────────┘
```

- Each shop card: compact (~80px tall), shows icon, name, cost, quantity
- Quantity shown as `×3` badge, grays out when `×0` (sold out)
- BUY button disabled if insufficient bar or quantity 0
- Scrollable if card list exceeds panel height
- Remove Card section at the bottom, always visible
- BUY button disabled mid-interaction (Hex Claim pending): prevent accidental purchase while placing a hex

---

## Unit Token on Hex Grid

Units must be distinguishable from hex ownership color while sharing the same color family.

**Design:** Filled circle, 38% of hex diameter, centered on the hex.
- Fill: player color at 90% opacity
- Border: 2px, player color darkened 40% (e.g. P1: `#2a7a74`)
- Letter inside: `S` for Soldier (white, bold, 10px)
- HP dots: row of small circles below the letter (filled = alive, empty = lost)
  - 2 HP: `●●`
  - 1 HP: `●○`

**When two units from the same player are on the same hex:** stack indicator — small `×2` badge top-right of the circle.

**When units are fighting (contact):** amber pulse ring around the hex (same pattern as MVP1 battle ring).

**Movement animation:** unit slides from center of one hex to center of next over 200ms (one tick). Must not block game state updates — purely visual interpolation.

---

## Mobile Landscape (896×414px)

```
┌────────────────────────────────────────────────────────────────┐
│ HUD (20px)  [▮▮▮░░ 5/10]              [opp ░░░░░░ 2/10]       │
├───────────────────────────────────┬───────────────────────────┤
│                                   │ SHOP (collapsible →)      │
│        HEX MAP                    │ [card] [BUY]              │
│                                   │ [card] [BUY]              │
│                                   │ [▼ hide]                  │
├───────────────────────────────────┴───────────────────────────┤
│ BAR: [▮▮▮▮▮░░░░░] 5/10                                        │ 18px
├────────────────────────────────────────────────────────────────┤
│ ┌────────────┐  ┌─────────┐  ══ ══     [→ PUSH]     [SHOP ▲] │ 140px
│ │ TOP CARD   │  │  ASIDE  │                                   │
│ │ [2] Sold.  │  │[1] +Bar │                                   │
│ │ Deploy a   │  │ 15% 15s │                                   │
│ │ soldier    │  │         │                                   │
│ │  [PLAY]    │  │ [PLAY]  │                                   │
│ └────────────┘  └─────────┘                                   │
└────────────────────────────────────────────────────────────────┘
```

- Shop collapses to a right-edge button `[SHOP ▲]`; tapping opens a bottom sheet overlay
- Card size reduced: top card ~110px × 130px, aside ~90px × 110px
- Bar strip: 18px (thinner but still segmented)
- HUD: 20px, minimal info only
- `[→ PUSH]` button moves to right of deck panel (thumb-reachable)

---

## Interaction States

| State | Visual |
|---|---|
| Bar insufficient for top card | Card dimmed 50%, cost badge red, PLAY button disabled |
| Aside slot empty | Dashed border, "push here" label |
| Aside slot occupied | Card shown, PUSH button on top card grayed |
| Hex Claim pending | Map highlighted with claimable hexes (bright border); deck panel shows "choose a hex" prompt |
| Card just played | Card flips away animation (200ms), next card slides up from behind |
| Unit fighting | Amber pulse ring on combat hex |
| Capital HP low (<5) | Capital hex ring pulses red |
| Win/loss | Full-screen overlay (same pattern as MVP1) |

---

## Implementation Notes

### Canvas (renderer.ts)
- Hex grid, buildings on hexes, unit tokens (animated march)
- Unit combat pulse ring reuses MVP1 battle ring pattern
- Capital hex: colored ring with HP counter label

### DOM overlays
- `#hud` — top strip (40px fixed)
- `#bar-meter` — 24px strip above deck panel
- `#deck-panel` — bottom 190px (fixed, always visible)
- `#shop-panel` — right sidebar 200px
- All panels: `background: var(--chrome)`, `border: 1px solid var(--ui-border)`

### CSS variables to define immediately
```css
:root {
  --bg: #1a1a2e;
  --chrome: #252545;
  --ui-border: #3a3a5c;
  --p1: #4ecdc4;
  --p2: #ff6b6b;
  --unclaimed: #2d3748;
  --bar-fill: #f7c948;
  --bar-empty: #3a3a5c;
  --card-unit: #2a3f6f;
  --card-building: #3d2e1e;
  --card-bar: #1e3d2e;
  --cost: #f7c948;
  --text-primary: #e8e8f0;
  --text-muted: #6b7280;
  --danger: #e53e3e;
  --win: #f6d860;
}
```

### Keyboard shortcuts
| Key | Action |
|---|---|
| Space | Play top card |
| E | Push top card to aside slot |
| F | Play aside card |
| B | Focus shop panel |
| ESC | Cancel pending hex claim |
