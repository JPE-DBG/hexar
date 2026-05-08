# Hexar — Milestone Plan

Generated from CLAUDE.md. 6 milestones, ~6.5 weeks solo developer.

---

## Feature Inventory

| # | Feature | Category | MVP? |
|---|---------|----------|------|
| 1 | Hex grid (axial coords, adjacency) | core-loop | Yes |
| 2 | Hex ownership (owned/unclaimed, capital) | core-loop | Yes |
| 3 | Claiming unclaimed hexes (10g, instant) | core-loop | Yes |
| 4 | Gold income (2/sec base per hex) | core-loop | Yes |
| 5 | Stepped maintenance (1/2/3 tiers) | balance | Yes |
| 6 | Gold building (+50% income, upgrades) | core-loop | Yes |
| 7 | Power building (+Power per level, upgrades) | core-loop | Yes |
| 8 | Research building (+TP at 0.2 TP/sec per level, upgrades) | core-loop | Yes |
| 9 | Building demolish (50% refund) | core-loop | Yes |
| 10 | Exponential upgrade costs | balance | Yes |
| 11 | Attack enemy hex (100g, Power check) | core-loop | Yes |
| 12 | Battle timer (6-15 sec countdown) | core-loop | Yes |
| 13 | Instant takeover (Power diff > 3) | core-loop | Yes |
| 14 | Counter-spend during battle | core-loop | Yes |
| 15 | Garrison tech (adjacent owned hexes +1 Power in defense, cap +2 from Garrison, 30 TP) | core-loop | Yes |
| 16 | Tech tree (12 techs, TP spending, no prerequisites) | core-loop | Yes |
| 17 | Auto-drop (negative income, grace period) | balance | Yes |
| ~~18~~ | ~~Conquest victory (60%/10s)~~ | ~~win-condition~~ | Removed — replaced by Capital Capture |
| 19 | Fortify active action (spend 40g, prevent instant-takeover for 90s, requires Fortify tech) | core-loop | Yes |
| ~~20~~ | ~~Time limit victory (30 min)~~ | ~~win-condition~~ | Removed — no timeout tiebreaker |
| 21 | Capital Capture = sole victory condition (capital taken → instant loss, all hexes unclaimed) | win-condition | Yes |
| 22 | Tick loop (100ms, 6 phases) | networking | Yes |
| 23 | WebSocket server + client connect | networking | Yes |
| 24 | Delta state sync | networking | Yes |
| 25 | Full snapshot on connect/reconnect | networking | Yes |
| 26 | Hardcoded test map (70 hexes) | core-loop | Yes |
| 27 | Canvas hex rendering | ui | Yes |
| 28 | Click-to-hex detection (pixel→axial) | ui | Yes |
| 29 | HUD (gold, TP, hex count, timers) | ui | Yes |
| 30 | Build menu (place/upgrade/demolish) | ui | Yes |
| 31 | Battle timer UI | ui | Yes |
| 32 | Tech tree UI | ui | Yes |
| 33 | Auto-drop UI (choice prompt) | ui | Yes |
| 34 | Disconnect handling (30s forfeit) | networking | Yes |
| 35 | 2-player lobby/matchmaking | networking | Yes |

**Post-MVP (deferred):** Alliances, diplomacy, special hex types, hero units, ranked ladder, procedural map gen, 3-4 player, Svelte UI rewrite, MessagePack, deployment to VPS.

---

## Dependency Graph

```
Hex grid math (#1) ← everything else
Tick loop (#22) ← economy (#4,5), battles (#12), victory (#21), delta (#24)
Hex ownership (#2) ← claiming (#3), buildings (#6-9), combat (#11-14)
Gold income (#4) ← buildings (#6-9), combat (#11), counter-spend (#14)
Stepped maintenance (#5) ← auto-drop (#17)
Power building (#7) ← combat (#11-13), counter-spend (#14)
Research building (#8) ← tech tree (#16)
Tech tree (#16) ← all tech effects (#15 and others within #16)
Garrison (#15) ← tech tree (#16)  [Garrison effect active only after player unlocks Garrison tech]
Fortify action (#19) ← tech tree (#16)  [action available only after player unlocks Fortify tech]
Combat (#11-13) ← victory: capital capture (#21)
WebSocket (#23) ← delta sync (#24), snapshot (#25), disconnect (#34)
Canvas render (#27) ← click detection (#28), all UI (#29-33)
```

---

## Milestone Plan

### M1 — Tick Loop + Rendering Pipeline ✅ DONE (~1 week)

**Goal:** Prove the full pipeline works: Go server ticks at 100ms, state reaches browser, Canvas draws hexes in real time.

**Features:** #1 (hex math), #22 (tick loop), #23 (WebSocket), #26 (hardcoded map), #27 (Canvas hex rendering), #24 (delta sync — minimal: full state each tick is fine here)

**Stubbed:**
- No game logic (just static hex grid sent to client)
- Hardcoded 70-hex map, all unclaimed
- No player actions yet
- Full state dump every tick (delta optimization comes later)

**Done when:**
- Browser shows colored hex grid that matches server state
- Changing a hex's owner in server code causes the Canvas to update within 100ms
- Two browser tabs connect to same game, see same grid

**Design risk:** Canvas hex rendering performance with 70+ hexes at 10fps update rate. If perf is bad → reduce to key-frame updates (only redraw dirty hexes).

---

### M2 — Economy + Expansion ✅ DONE (~1 week)

**Goal:** Validate the core economic loop — does claiming hexes and earning gold feel right at CLAUDE.md's numbers?

**Features:** #2 (ownership), #3 (claiming unclaimed), #4 (gold income), #5 (stepped maintenance), #28 (click-to-hex), #25 (snapshot on connect)

**Stubbed:**
- No buildings yet (raw income only)
- No combat (only unclaimed hex claiming)
- No auto-drop (maintenance tracked but no penalty yet)
- Minimal HUD (gold counter, hex count)

**Done when:**
- Click hex → send "claim" action → server validates adjacency + gold ≥ 10 → hex turns player color
- Gold accrues visibly in HUD at correct rate (2/sec/hex - maintenance)
- Headless Go test: 10 hexes at T=30s produces correct gold (matches CLAUDE.md math)
- Two players expand from opposite corners, meet in the middle

**Design risk:** Expansion feels too fast or too slow. If land-grab is over in 30 seconds → map might be too small. If it takes 5+ minutes → too large or income too slow. Adjust map size or starting gold.

---

### M3 — Buildings + Combat ✅ DONE (~1.5 weeks)

**Goal:** First real gameplay — players build, fight, and territory changes hands. This is the first **playable** milestone.

**Features:** #6 (Gold building), #7 (Power building), #8 (Research building), #9 (demolish), #10 (exponential costs), #11 (attack enemy hex), #12 (battle timer), #13 (instant takeover), #29 (HUD), #30 (build menu), #31 (battle timer UI)

**Notes (actual vs planned):** Research building (#8) was added during this milestone rather than deferred. All three building types (Gold/Power/Research) are fully functional with correct cost formulas and UI labels. Significant polish iteration: Economy→Gold/Defense→Power rename, delta label fix (L1 shows +1.9/s not +0.9/s), constants extracted from magic numbers.

**Stubbed:**
- No counter-spend (defender is passive)
- No auto-drop (just track negative income in logs)
- No victory conditions (play until bored)

**Done when:**
- Player can build Gold/Power/Research, see income increase / Power increase / TP increase
- Player attacks adjacent enemy hex, battle timer counts down, hex flips on win
- Power diff > 3 = instant takeover (no timer)
- Power diff ≤ 0 = attack rejected (no gold spent)
- Exponential upgrade costs visible in build menu
- Two humans can play a full land-grab + border fight session

**Design risk:** Combat feels too deterministic (higher Power always wins, no counterplay). This is expected without counter-spend — M4 adds tension. If Power stacking makes fights one-sided even at diff=1, battle duration formula may need tuning.

---

### M4 — Counter-Spend + Auto-Drop ✅ DONE (~1 week)

**Goal:** Complete the defensive gameplay loop and economic pressure system. Games now have real tension and economic collapse risk.

**Features:** #14 (counter-spend), #17 (auto-drop with grace period + UI), #32 (tech tree UI — stub: all 12 slots visible, TP spending works, no effects yet), #33 (auto-drop UI), #5 (maintenance fully enforced)

**Note:** Research building (#8) already done in M3. TP rate is 0.2 TP/sec per level.

**Stubbed:**
- All 12 tech effects inactive — UI shows slots, costs, and TP balance; unlocking spends TP and marks the tech as owned, but nothing changes in gameplay yet (effects wired in M5)
- No victory conditions yet

**Done when:**
- Defender can spend gold during battle to boost Power (+1/sec, total cap +3 or seconds remaining — whichever is lower)
- Player with 20+ hexes and no Economy buildings hits negative income → 10s grace → forced drop
- Auto-drop UI prompts player to choose which hex to shed
- Research building generates TP visibly in HUD at 0.2 TP/sec per level
- Headless test: verify counter-spend cannot exceed +3 total

**Design risk:** Counter-spend might feel unfair (defender always wins by outspending). The +3 cap and time-remaining limit should prevent this; if battles always flip to defender → review cap or attack cost. Auto-drop grace period (10s) might be too short if income fluctuates from battles — may need to extend or add hysteresis.

---

### M5 — Tech Tree + Victory Conditions ✅ COMPLETE

**Goal:** Complete game with win conditions. A full 15-30 minute match is playable end-to-end.

**Features:** #15 (Garrison), #16 (all 12 techs functional), #19 (Fortify action), #21 (Capital Capture victory)

**Stubbed:**
- No disconnect handling (both players must stay connected)
- No lobby (hardcoded 2-player join)
- Minimal victory screen (text overlay)

**Done when:**
- All 12 techs unlock and apply their effects globally
  - Blitz: unclaimed claims cost 0g
  - Vanguard: next attack within 12s costs 50g after a capture
  - Iron Grip: all hexes +1 Power
  - Siege Mastery: attacker battle timers ×0.6; tie → attacker wins
  - Prosperity: each Gold building +1/sec
  - Supply Lines: maintenance 0.9/1.8/2.7 per tier
  - Compound Growth: Gold buildings ×1.25 output
  - Garrison: adjacent owned hexes +1 Power in defense (cap: +2 from Garrison; +3 total with counter-spend)
  - Fortify: 40g action grants instant-takeover immunity for 90s on one hex
  - War Chest: capture enemy hex → recover 30g
  - Reclamation: recapture own hex for 50g
  - Resilience: drop grace 20s, drop refund 70%
- Game ends when: player captures the enemy capital → instant win, all loser hexes unclaimed
- Tech tree UI (built in M4) now shows unlocked techs as active with visual distinction; all 12 effects apply
- Victory screen shows winner and reason
- Full 15-30 min game is completable between two human players

**Power display conventions (decided during M5 playtesting):**
- Hex canvas labels show **effective combat power**, not raw building level for Power buildings
  - Power L1 + Iron Grip → label shows "P2" (effective), not "P1" (level)
  - Capital + Power L1 + Iron Grip → label shows "P3" (1 innate + 1 level + 1 IG)
- Non-Power owned hexes show a small yellow power badge if power > 0 (capital innate, Iron Grip)
- Garrison is **excluded from static power display** — it's a defense-battle-only bonus, shown as `"+N def"` note in build menu
- Fortify timer rendered as **shrinking lime border segments** (clockwise from top), no text overlay

**Design risk:** 12 techs is significantly more scope than original 4. Consider implementing the 6 cheap techs (≤40 TP) first and the 6 expensive techs second. Fortify introduces the first "active action from a tech" — may need new UI affordance. Siege Mastery's 40% timer reduction needs verification that the minimum 3s floor doesn't create degenerate battles.

---

### M6 — Polish + Networking Robustness (~1 week)

**Goal:** Game is reliably playable by two people on the same network without crashes or desyncs.

**Features:** #24 (proper delta sync — only changed hexes), #25 (reconnect with snapshot), #34 (disconnect → 30s forfeit), #35 (simple lobby: create/join game)

**Stubbed:**
- No matchmaking (manual room codes)
- No auth (anonymous players)
- localhost only

**Done when:**
- Delta messages are <500 bytes per tick during steady state (not full state dumps)
- Player disconnects → game continues 30s → forfeit if no reconnect
- Player reconnects → receives full snapshot → resumes seamlessly
- Simple lobby: player 1 creates game, gets code; player 2 joins with code
- No observed desyncs over a full 30-min game

**Design risk:** Delta sync bugs causing ghost state (client shows hex as owned, server disagrees). Mitigation: periodic full-state checksum comparison (every 5s), force resync on mismatch.

---

## Post-MVP (deferred)

- **Tech Tree Benefit Display (UX Enhancement):** Show quantitative benefits in tech tree UI before unlocking
  - Dynamic calculation based on current game state
  - Example displays: "Prosperity: +3.0 gold/s total (3 Economy buildings)", "Supply Lines: -2.0 gold/s maintenance (current: 15 hexes)", "Iron Grip: All 12 owned hexes +1 Power"
  - Helps players evaluate tech ROI before spending TP
  - Implementation: Add benefit calculator function to techtree.ts, pass GameState to update(), display benefit below static description
  - Estimated effort: 2-3 hours
  - Rationale for deferral: Static tech descriptions are sufficient for MVP gameplay; quantitative feedback is nice-to-have polish
- Procedural map generation (replace hardcoded map)
- 3-4 player mode
- Alliances / diplomacy
- Special hex types (terrain)
- Hero units
- Ranked ladder / matchmaking
- Persistent accounts (OAuth)
- Match history (SQLite)
- Remote deployment (Fly.io)
- Svelte UI rewrite (if DOM overlays get complex)
- MessagePack serialization (if JSON bandwidth is a problem)
- Playwright E2E tests

---

## Timeline Summary

| Milestone | Duration | Cumulative | Playable? |
|-----------|----------|------------|-----------|
| M1 — Tick + Render | ~1 week | Week 1 | Visual only |
| M2 — Economy + Expansion | ~1 week | Week 2 | Expansion game |
| M3 — Buildings + Combat | ~1.5 weeks | Week 3.5 | **First real game** |
| M4 — Counter-Spend + Auto-Drop | ~1 week | Week 4.5 | Defensive play |
| M5 — Tech Tree + Victory | ~1.5 weeks | Week 6 | **Complete game** |
| M6 — Polish + Networking | ~1 week | Week 7 | Robust game |

Critical path: M1 → M2 → M3 (each depends on the previous). M4/M5 could partially overlap if combat and tech are developed in parallel, but counter-spend needs the battle system from M3.
