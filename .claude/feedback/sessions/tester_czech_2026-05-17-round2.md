---
tester: Tester A (Czech)
date: 2026-05-17
feedback_type: UX / Game Mechanics
status: RAW
issues_total: 8
issues_implemented: 0
issues_rejected: 0
issues_deferred: 0
---

# Tester Session 2: UX Refinements & Game Mechanic Proposals

## Issue 1: Empty Hex Has No Name
- **Type:** UX / Clarity
- **Observation:** Unclaimed hexes have no label in the sidebar — selecting one shows nothing meaningful
- **Suggestion:** Give them a display name, e.g. "Empty Field", "Unclaimed", or "Wilderness"
- **Design question:** What name fits the game tone? "Unclaimed Hex" is functional; "Empty Field" is more evocative
- **Status:** ⬜ Open

---

## Issue 2: Capital Power Label Position
- **Type:** UX / Layout
- **Observation:** Capital hex shows `Pwr: 1` above the building name in the sidebar — but capitals can simultaneously have a Power, Economy, or Research building on top of the innate Power 1
- **Expected:** Power label should appear below building info (not above), consistent with how other hexes render it; innate capital power should included from building-added power. i.e. in case capital have gold mine it should show "+x g/s Pwr: 1" under hex name.
- **Status:** ⬜ Open

---

## Issue 3: Smart Build ON Should Disable Sticky Tool Mode
- **Type:** UX / Behavior
- **Observation:** With Smart Build ON, pressing Q/W/E still enters the sticky tool mode (highlights button, waits for click). This creates two competing interaction models active at the same time.
- **Expected:** When Smart Build is ON, Q/W/E should only execute instantly on the hovered hex (or do nothing if no valid target) — sticky tool mode should be fully disabled; toggling the tool by pressing hotkey twice would be confusing when Smart is ON
- **Related:** Smart Build toggle in BUILD section (current: OFF by default)
- **Status:** ⬜ Open

---

## Issue 4: Build and Upgrade via Same Hotkey
- **Type:** UX / Feature
- **Observation:** Separate build vs upgrade flow is redundant — Q should build an Economy building on an empty hex OR upgrade an existing one, depending on hex state
- **Proposed behavior:**
  - Empty own hex + press Q → build Economy L1
  - Own hex with Economy L1 + press Q → upgrade to L2
  - Own hex with non-Economy building + press Q → no action (or show error)
- **Design question:** How to communicate the cost? Build cost (60g) vs upgrade cost (120g, 240g…) need to be visible before the hotkey fires. Sidebar context area or tooltip on the hex could show the pending cost.
- **Status:** ⬜ Open

---

## Issue 5: Attack by Clicking Opponent Hex Directly
- **Type:** UX / Feature
- **Observation:** Current flow: select own hex → select enemy hex → press A. Tester expects: hover enemy hex → click it to attack (when attacker has sufficient power)
- **Design question:** Attack cost (100g) and power diff preview need to be visible before committing. Without an explicit A key confirmation step, the UI must show cost/result in a hover tooltip or the sidebar must update on hover (not just on click-select).
- **Alternative:** Keep A key but make it fire on hovered enemy hex when Smart Build is ON, consistent with the Smart Build model
- **Status:** ⬜ Open

---

## Issue 6: Auto Counter-Spend Should Cap at Minimum Needed to Win
- **Type:** Game Mechanic / Feature
- **Observation:** Tester expects an "auto defend" button that spends just enough gold to flip the battle outcome — not the full cap
- **Current:** Counter-spend is manual (click button or C key per 50g burst)
- **Proposed:** When counter-spend is not needed but there are still pointts left game will do nothing to save gold
- **Design considerations:**
  - Risk: removes defender agency / tactical bluffing
- **Status:** ⬜ Open

---

## Issue 7: Buildings Take Time to Upgrade
- **Type:** Game Mechanic / Feature
- **Observation:** Buildings upgrade instantly — tester expected a construction delay
- **Proposal:** Each upgrade takes 5 seconds (constant across all levels and building types)
  - During construction: building shows previous level stats; hex is marked "upgrading" with timer like fortify
  - Construction can not be canceled, but hex can be demolished or sold completly
- **Design considerations:**
  - Needs timer visualization (hex timer, same pattern as Fortify?)
  - Interacts with battle: can you upgrade while defending? While attacking from that hex?
  - Significant mechanic change — needs playtesting to validate 5s constant feels right across all levels
- **Status:** ⬜ Open

---

## Issue 8: Mobile — Hex Info Panel Placement
- **Type:** Mobile UX / Layout
- **Observation:** On mobile, hex info (building name, power, output) is hard to reach or obscured in the current bottom-bar layout
- **Suggestion:** Move hex info to below the player resource info (gold/hexes/research counters) at the top of the screen
- **Design question:** Does this conflict with HUD layout on small screens? Portrait vs landscape may need different treatments.
- **Status:** ⬜ Open
