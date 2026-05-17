---
tester: Tester A (Czech)
date: 2026-05-17
feedback_type: UX/UI Polish
---

# Tester Session: UI Polish & Quality-of-Life

## Issue 1: Hex Info Panel Missing Details
- **Type:** UX / Clarity
- **Observation:** When selecting a hex, the sidebar shows hex coordinates but doesn't explain what building is there
- **What I wanted:** Clear description of what each building does + its output
- **Suggested format:**
  ```
  Gold Mine (Level 2)
  +3.9 g/s
  
  [or]
  
  Barracks (Level 1)
  Pwr: 1
  
  [or]
  
  Laboratory (Level 3)
  +0.6 tech/s
  
  Position: [Q5] (on next line)
  ```

---

## Issue 2: Button Size Inconsistency
- **Type:** UX / Polish
- **Observation:** Context buttons (Attack, Upgrade, Demolish) are different size than Build buttons (Economy, Power, Research)
- **Impact:** Looks unprofessional; hard to predict button hit areas
- **Suggestion:** Standardize all button sizes (and icon sizes too — icons vary in size across buttons)

---

## Issue 3: Smart Building Feature (Quality-of-Life)
- **Type:** Feature Request / QoL Polish
- **Priority:** MEDIUM (Post-MVP, high-value for experienced players)
- **Idea:** Like League of Legends — hover hex, press shortcut key to instantly place building without clicking
  - Press `Q` while hovering → builds Economy there
  - Press `U` while hovering → upgrades building
  - Press `A` while hovering → attacks enemy hex
  - Press `F` → fortify, etc.
- **UX question:** Where should toggle be? Sidebar checkbox? Settings?
- **Benefit:** Faster gameplay for experienced players (optional, not default behavior)
- **Rationale:** Now that MVP is shipped, this is valuable polish for repeat players and speedrunners

---

## Issue 4: Attack Power Calculation Missing Garrison Info
- **Type:** Game Logic / UX
- **Observation:** When previewing an attack, the UI shows defender's Power but doesn't account for adjacent hexes (Garrison bonus)
- **Expected:** Attack validation should show: "Defender Power 2 + Garrison +1 = 3 total — you need 4 to attack"
- **Impact:** Misclick attacks that fail due to Garrison bonus

---

## Issue 5: Tech Tree Close Button Missing
- **Type:** UX / Polish
- **Observation:** Tech tree modal has no X button in top right — have to ESC or click outside
- **Suggestion:** Add standard close button (X) in top-right corner

---

## Issue 6: Empty Hex Demolish/Sell Shows No Feedback
- **Type:** UX / Clarity
- **Observation:** When selecting an unclaimed/empty hex, Demolish and Sell buttons appear but should show "0g" or be disabled
- **Current:** Buttons are clickable but do nothing
- **Better:** Show "0g" (no refund) or disable + gray out
