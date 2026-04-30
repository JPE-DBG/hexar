# Hexar — Real-Time Hex Strategy Game

## Overview

Hexar is a fast-paced, real-time multiplayer hex strategy game inspired by Antiyoy. Players conquer hexes through economic growth and strategic upgrades, competing on small maps with 30-minute target game length.

**Core Loop:** Earn resources → Upgrade buildings → Gain power → Conquer adjacent hexes → Repeat

---

## Game Concept

**Genre:** Real-time strategy (RTS), economic conquest  
**Players:** 2-4 (starting with 1v1)  
**Game Duration:** ~25-30 minutes  
**Map:** Hexagonal grid, 15-25 hexes per player starter territory  
**Win Conditions:** Conquest, Tech Dominance, or Time Limit (see Victory Conditions)

### Design Pillars

1. **Economic Competition:** Victory through resource management and smart upgrades, not reflexes
2. **Meaningful Decisions:** Each building choice and tech unlock trades off against alternatives
3. **Tense Endgame:** Border warfare and tight resource battles in final minutes
4. **No Snowballing:** Maintenance costs and battle attrition keep powerful players vulnerable

---

## Core Mechanics

### Hexes & Ownership

- **Player hexes (owned):** Controlled by a player, generates income, has buildings
- **Unclaimed hexes (empty):** No owner, Power 0, no buildings
- **Starting position:** Each player starts with exactly 1 hex (their capital)
- **Owned hex generates:** 2 resources/sec (base)
- **Maintenance cost:** 1 resource/sec per hex controlled (owned)
- **Net income per hex:** +1 resource/sec (before buildings/upgrades)

### Resources

- **Gold:** Main currency for buildings, upgrades, and attacks
- **Tech Points (TP):** Earned from Research buildings (0.1 TP/sec per Research level), spent on tech tree

### Buildings (One per Hex)

Each hex can have **one building** of three types. Building can be demolished (refund 50%) and replaced.

#### Economy Building
- **Effect:** +50% resources/sec from that hex (stacks additively with other bonuses)
- **Build cost:** 80 gold
- **Upgrade cost (Exponential):** L1=40, L2=80, L3=160, L4=320, L5=640 (doubles each level)
- **Reward (Linear):** +0.5 resources/sec per level (constant gain)
- **Max level:** Unlimited (incremental)

#### Defense Building
- **Effect:** +1 Power per level (affects combat)
- **Build cost:** 60 gold
- **Upgrade cost:** 30 gold/level
- **Max level:** Unlimited (incremental)

#### Research Building
- **Effect:** +0.1 TP/sec per level (fuel for tech tree)
- **Build cost:** 80 gold (reduced from 120)
- **Upgrade cost (Exponential):** L1=40, L2=80, L3=160, L4=320, L5=640 (doubles each level)
- **Reward (Linear):** +0.1 TP/sec per level (constant gain)
- **Max level:** Unlimited (incremental)

---

## Combat System

### Attack Rules

- **Requirement:** Attacker Power > Defender Power (strictly greater than, ≤ fails with no cost)
- **Attack cost:** 100 gold (only attacker pays)
- **Adjacent only:** Can only attack hexes touching your hex

### Battle Resolution

**Attacking Empty Hex (Power 0):**
- If attacker Power ≥ 1: **Instant takeover** (no battle, no resource cost beyond attack)
- Reason: Early game expansion should be fast; bottleneck is economy, not combat

**Attacking Enemy Hex (Owned):**
- **Requirement:** Attacker Power > Defender Power (strictly greater)
- If Power diff > 3: **Instant takeover** (no battle timer, immediate conquest)
- If Power diff = 1, 2, or 3: **Standard battle** (6-15 sec, see below)
- If Power diff ≤ 0: **Attack fails** (attacker wasted 100 gold, can retry)

**Standard Battle (1-3 power difference):**

```
Battle Duration = 5 + (Attacker Power + Defender Power) / 2 seconds

Min: 6 seconds (Power 1 vs 1)
Max: ~15 seconds (Power 10 vs 8)
```

- Battle timer visible to both players during countdown
- Attacker takes hex if battle completes
- Defender's hex becomes neutral (not captured by attacker) if they own it

### Attacker Hex Damage (Cost of Victory)

Winning attacker's hex loses levels based on power differential:

```
Power diff = +1:  Attacker loses 0 levels (clean win)
Power diff = +2:  Attacker loses 0 levels (solid win)
Power diff = +3:  Attacker loses 1 level (costly)
Power diff = +4+: Attacker loses 2+ levels (pyrrhic)
```

- Example: Your Power 5 hex attacks enemy Power 2 hex (diff +3). You win but drop to Power 4.

### Defender Options (Active Defense)

**Counter-Spend (Capped)**

During the battle countdown, defender can spend resources to boost Defense Power:
- Cost: 50 gold/second during battle
- Effect: +1 Power per second spent (capped at +3 power total)
- Example: Battle at 9 seconds, you're losing 3 vs 5. Spend 150 gold over 3 sec → Power jumps to 6, you win battle
- Cap prevents defender from "buying" complete victory; makes battles tactical instead of pay-to-win

---

## Tech Tree

Research buildings generate Tech Points. Spend TP to unlock perks (global bonuses):

| Tech | Cost | Effect |
|------|------|--------|
| Iron Grip | 30 TP | All hexes +1 Power |
| Production Boom | 40 TP | All hexes +30% resource generation |
| Efficient Conquest | 35 TP | Attack cost reduced to 75 gold |
| Fortified Borders | 25 TP | Enemy attacks cost them +25 gold |
| Economic Synergy | 50 TP | Economy buildings give +60% (instead of +50%) |
| Blitzkrieg | 40 TP | Reduce battle duration by 2 seconds (min 5 sec) |
| Garrison | 60 TP | During battle, adjacent hexes can reinforce defender (add their Power) |

**Tech Level:** Total number of techs unlocked. Reaching Tech Level 3 is significant (see Victory Conditions).

**Rationale:** Costs reduced by ~40% from original to make tech tree achievable in 30-min games. First few techs are cheap to encourage early tech investment as a viable alternative to pure military.

---

## Economy Example (Early Game)

```
T=0s:      Control 1 hex (capital), Power 1, no buildings
           Income: 2/sec, Maintenance: 1/sec → Net +1/sec
           Gold: 0

T=0-20s:   Accumulate 20 gold (1/sec × 20)

T=20s:     Spend 100 gold → Attack adjacent unclaimed hex (instant, Power 1 > 0)
           Now control 2 hexes
           Income: 4/sec, Maintenance: 2/sec → Net +2/sec
           Gold: 0

T=20-40s:  Accumulate 40 gold (2/sec × 20)

T=40s:     Spend 80 gold → Build Economy on hex 1
           Income now: hex 1 = 2 × 1.5 = 3/sec, hex 2 = 2/sec → Total 5/sec
           Maintenance: 2/sec → Net +3/sec
           Gold: 0

T=40-80s:  Accumulate 240 gold (3/sec × 80)

T=80s:     Spend 60 → Build Defense on hex 2
           Spend 30 → Upgrade Defense Level 1 (Power now 2)
           Gold remaining: 150

T=80-150s: Expand to 3-4 hexes via instant conquest of unclaimed hexes
           Now 4 hexes = 4 maintenance, ~8-10 income/sec
           Spend resources upgrading buildings (L2 Defense costs 60 gold for hex 2)

T=150s+:   Meet enemy around this time
           Accumulated ~300+ gold for upgrades
           Have ~4-5 hexes, Power 2-3 depending on defense investment
           Ready for first border skirmish
```

---

## Victory Conditions

### 1. Conquest Victory (Primary)
- Hold **50% of map for 10 consecutive seconds**
- Timer resets if you drop below 50%
- This is the most common win condition

**Strategic implications:** Aggressive players win by pushing early; defensive players must hold the line and counter-attack to reset timer.

**How to stop opponent:**
- Attack their border hexes aggressively
- Bring them below 50%, reset their timer
- Race to 50% yourself

### 2. Tech Dominance (Mid-Game Alternative)
- Reach **Tech Level 3 AND hold 35% map simultaneously** for 10 consecutive seconds
- OR reach **Tech Level 5 alone** (no map control requirement)
- Rewards investing in research as a viable win condition

**Strategic implications:** Tech rush is faster than pure conquest (tech scaling matters). Early investment in Research building pays off. Opponent can counter by military pressure.

**How to stop opponent:**
- Rush military early while they tech
- Attack their Research hexes specifically
- Keep them pinned defending, prevent expansion
- Example: While opponent gets Tech 2-3, you've conquered 40% through aggression, win by Conquest

### 3. Time Limit (Tiebreaker)
- At 30 minutes, highest hex count wins
- Rarely triggers in well-balanced games

---

## Balance Rules & Constraints

### Exponential Upgrade Costs (Prevents Snowballing)
- Economy and Research buildings use exponential cost scaling: L1=40, L2=80, L3=160, L4=320, L5=640 (doubles each level)
- **Linear rewards:** Each level provides constant +0.5 resources/sec (or +0.1 TP/sec for Research)
- **Effect:** Early game upgrades are cheap (quick power spikes). Late game upgrades cost exponentially more, capping power growth
- **Math:** L1-L4 costs ~840 gold total. L5 adds 640. Each additional level doubles cost.
- **Example:** Player A has 10 hexes earning 20 resources/sec. Upgrading to L5 on 3 hexes = 1920 gold = 96 seconds of savings. By then, Player B has scaled up too.

### Maintenance System (Prevents Extreme Expansion)
- Each hex costs 1 maintenance/sec to hold
- Forces quality over quantity
- At ~10+ hexes, income caps out without production upgrades
- If maintenance exceeds income, slowest hexes auto-drop (player chooses order)

### Combat Attrition
- Attackers lose levels based on power difference (0 levels at +1-2 diff, 1 level at +3, 2+ at +4+)
- Discourages overkill attacks
- Creates cost to conquest

### Tech Scaling (Achievable in 30 min)
- Tech costs reduced for viability (first 3 techs = 130 TP total, ~1300 seconds with 1 Research hex)
- Tech progression is a viable win condition (Tech Level 3 + 35% map in ~15 min)
- Opponents must balance military pressure with allowing tech growth

### Early Game Parity
- All players start with 1 capital hex, Power 1, no buildings
- First minute: rapid hex expansion (unclaimed hexes taken instantly)
- 5-10 minute mark: players meet at borders
- Game stays competitive until ~15-minute mark if balanced play

---

## Design Decisions

### No "Neutral" Hex State (Only Unclaimed)
- **Simplified model:** Hexes are either **owned by a player** or **unclaimed (empty)**
- **Dropped "neutral" concept:** Initial design had unclaimed hexes revert to "neutral" after conquest, but this adds complexity without benefit
- **Why:** 
  - When you attack an enemy hex and win, you capture it directly (you own it)
  - No intermediate "neutral" state where other players can immediately reclaim it
  - Simpler UI: only 2 states instead of 3
  - Matches Antiyoy's design (conquered = owned)

### Instant Takeover on Unclaimed Hexes Only
- Players can instantly claim any adjacent unclaimed hex if they have Power ≥ 1
- Conquering **enemy hexes** always requires battles (6-15 sec or instant if diff > 3)
- Effect: Early game is about **speed of expansion** into unclaimed territory. Midgame is about **tactical battles** over enemy hexes.

### Adjacent-Only Conquest
- Reinforces territorial control and front-line battles
- Prevents "teleport attacks" across map
- Creates natural borders and defensive positions

### Only Attacker Pays
- Defender can react (counter-spend, capped at +3) rather than commit upfront
- Asymmetric costs reward smart positioning
- Encourages aggressive expansion but punishes poor decisions

### Power Difference Auto-Takeover (diff > 3, Enemy Hexes Only)
- Instant conquests feel rewarding (massive advantage in hex battles)
- Close fights (1-3 diff) are tense and tactical
- Prevents "grinding" low-power battles between evenly matched players
- Limited to enemy hexes only; unclaimed hexes always taken instantly

### Exponential Upgrade Costs with Linear Rewards
- **Cost:** Doubles each level (40, 80, 160, 320, 640...)
- **Reward:** Constant per level (+0.5 resources/sec or +0.1 TP/sec)
- **Effect:** Early power spikes are cheap and quick. Late-game upgrades cost exponentially more, capping power growth
- **Why:** Prevents snowballing by making the path to dominance prohibitively expensive. No player can spam upgrades to 10+ levels without massive time investment

---

## Open Questions / TODO

### Map Size & Pacing Analysis

**Recommended:** 60-80 total hexes per 1v1 map

**Rationale:**
- Players start with 1 hex each (2 total)
- ~30-40 unclaimed hexes available
- ~25-30 hexes per player after midgame (realistic distribution)

**Pacing with 70 hexes (example):**
- **T=0-3min:** Each player conquers 3-4 unclaimed hexes → 4-5 hexes each
- **T=3-5min:** Meet at borders (players are ~5 hexes from each other)
- **T=5-15min:** Border skirmishes, some territory trades hands
- **T=15-25min:** One player pushes toward 50% (35 hexes) or secures Tech Level 3
- **T=25-30min:** Final race to victory condition

**Too small (30 hexes total):** Players meet at T=1min, constant warfare, no economy buildup, RNG-heavy
**Too large (120+ hexes):** Players farm 15+ minutes unopposed, snowball guaranteed, long game

### Playtesting Needs

- [ ] **Exponential cost feel:** Does progression curve feel right? Too fast/slow?
- [ ] **Counter-spend cap:** Is +3 power cap balanced? Create interesting battles?
- [ ] **Map size:** Does 70-hex map hit 30-min target? Adjust if needed.
- [ ] **Conquest threshold:** Does 50% + 10 sec create tense endgame?
- [ ] **Tech viability:** Do players build Research? Or pure military?

### Mechanical Unknowns

- [ ] **Player starting position:** Do all players start equidistant? Random corners? (Recommend: opposite corners for 1v1)
- [ ] **Capital hex rules:** Should capital be undestroyable? Or can it be taken like any hex? (Recommend: can be taken, increases risk)
- [ ] **Building demolish refund:** Is 50% refund fair or should it be 100%?
- [ ] **Multiple battles:** Can hex be attacked by multiple enemies simultaneously? (Recommend: one attacker at a time, queue battles)

### Future Mechanics (Post-MVP)

- [ ] Alliances (2v2 mode with shared resources)
- [ ] Diplomacy (trade, temporary truces)
- [ ] Special hex types (mountains, water, resources)
- [ ] Hero units (unique units on hexes)
- [ ] Persistent progression (ranked ladder, cosmetics)

---

## Implementation Notes

- **Start with 1v1** — Easier to balance before adding 3-4 player variants
- **Hex grid rendering:** Use offset or axial coordinates for grid math
- **Real-time simulation:** Tick-based economy (every 100ms or server-dependent)
- **Networking:** Send only delta state (hex ownership changes, resource updates) not full state
- **UI priority:** Show hex Power prominently, battle timer clearly, resource flow transparent

