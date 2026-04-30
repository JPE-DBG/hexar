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

- Each hex can be owned by a player, be neutral, or empty
- **Owned hex generates:** 2 resources/sec (base)
- **Maintenance cost:** 1 resource/sec per hex controlled (owned or neutral claimed)
- **Net income per hex:** +1 resource/sec (before buildings/upgrades)

### Resources

- **Gold:** Main currency for buildings, upgrades, and attacks
- **Tech Points (TP):** Earned from Research buildings (0.1 TP/sec per Research level), spent on tech tree

### Buildings (One per Hex)

Each hex can have **one building** of three types. Building can be demolished (refund 50%) and replaced.

#### Economy Building
- **Effect:** +50% resources/sec from that hex (stacks additively with other bonuses)
- **Build cost:** 80 gold
- **Upgrade cost:** 40 gold/level
- **Max level:** Unlimited (incremental)

#### Defense Building
- **Effect:** +1 Power per level (affects combat)
- **Build cost:** 60 gold
- **Upgrade cost:** 30 gold/level
- **Max level:** Unlimited (incremental)

#### Research Building
- **Effect:** +0.1 TP/sec per level (fuel for tech tree)
- **Build cost:** 120 gold
- **Upgrade cost:** 50 gold/level
- **Max level:** Unlimited (incremental)

---

## Combat System

### Attack Rules

- **Requirement:** Attacker Power > Defender Power (strictly greater than, ≤ fails with no cost)
- **Attack cost:** 100 gold (only attacker pays)
- **Adjacent only:** Can only attack hexes touching your hex

### Battle Resolution

**Power Difference Rule:**

```
Power diff = Attacker Power - Defender Power

If diff > 3:
  → Instant takeover (no battle timer)
  → Attacker takes hex immediately

If diff = 1, 2, or 3:
  → Standard battle (see below)

If diff ≤ 0:
  → Attack fails (attacker wasted 100 gold but can retry)
```

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

**Option A: Counter-Spend (Recommended start)**

During the battle countdown, defender can spend resources to boost Defense Power:
- Cost: 50 gold/second during battle
- Effect: +1 Power per second spent
- Example: Battle at 9 seconds, you're losing 3 vs 5. Spend 100 gold over 2 sec → Power jumps to 5, you win

---

## Tech Tree

Research buildings generate Tech Points. Spend TP to unlock perks (global bonuses):

| Tech | Cost | Effect |
|------|------|--------|
| Iron Grip | 50 TP | All hexes +1 Power |
| Production Boom | 80 TP | All hexes +30% resource generation |
| Efficient Conquest | 60 TP | Attack cost reduced to 75 gold |
| Fortified Borders | 40 TP | Enemy attacks cost them +25 gold |
| Economic Synergy | 100 TP | Economy buildings give +60% (instead of +50%) |
| Blitzkrieg | 80 TP | Reduce battle duration by 2 seconds (min 5 sec) |
| Garrison | 100 TP | During battle, adjacent hexes can reinforce defender (add their Power) |

**Tech Level:** Total number of techs unlocked. Reaching Tech Level 8 is significant (see Victory Conditions).

---

## Economy Example (Early Game)

```
T=0s:      Control 1 hex, Power 1, no buildings
           Income: 2/sec, Maintenance: 1/sec → Net +1/sec
           Gold: 0

T=0-20s:   Accumulate 20 gold (1/sec × 20)

T=20s:     Spend 100 gold → Attack adjacent empty hex (instant, Power 1 > 0+3)
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

T=80-100s: Enemy attacks hex 2 (their Power 2 vs your Power 2)
           Battle lasts 5 + (2+2)/2 = 9 seconds
           You counter-spend 50 gold → Power +1, you win battle
           Enemy's hex dropped to Power 1
           You control 2 hexes, opponent still has theirs
```

---

## Victory Conditions

### 1. Conquest Victory (Primary)
- Hold **70% of map for 20 consecutive seconds**
- Timer resets if you drop below 70%
- This is the most common win condition

**How to stop opponent:**
- Attack their border hexes aggressively
- Bring them below 70%, reset their timer
- Race to 70% yourself

### 2. Tech Dominance (Long-Game)
- Reach **Tech Level 5 AND hold 50% map simultaneously** for 15 consecutive seconds
- OR reach **Tech Level 8 alone** (no map control requirement)
- Rewards investing in research but requires map presence or extreme tech lead

**How to stop opponent:**
- Rush military early while they tech
- Attack their Research hexes specifically
- Keep them pinned defending, prevent expansion
- Example: While opponent gets Tech 3-4, you've conquered 55% through aggression

### 3. Time Limit (Tiebreaker)
- At 30 minutes, highest hex count wins
- Rarely triggers in well-balanced games

---

## Balance Rules & Constraints

### Maintenance System (Prevents Snowballing)
- Each hex costs 1 maintenance/sec to hold
- Forces quality over quantity
- At ~10+ hexes, income caps out without production upgrades
- If maintenance exceeds income, slowest hexes auto-drop (player chooses order)

### Combat Attrition
- Attackers lose levels based on power difference
- Discourages overkill attacks
- Creates cost to conquest: attacking with Power 10 vs Power 2 (instant) is free, but Power 5 vs Power 3 costs attacker a level

### Tech Scaling
- Tech is slow (0.1 TP/sec per Research level)
- To reach Tech Level 8: ~800 TP needed = 8,000 seconds = 2+ hours (impossible in 30-min game)
- Realistic max: Tech Level 4-5 in a game
- Tech advantage is real but not insurmountable with military pressure

### Early Game Parity
- All players start identical (1 hex, Power 1, no buildings)
- First 2 minutes are roughly equal economy
- First conquest (empty hex) at ~20 sec gives first advantage
- Game stays competitive until ~10-minute mark

---

## Design Decisions

### Adjacent-Only Conquest
- Reinforces territorial control and front-line battles
- Prevents "teleport attacks" across map
- Creates natural borders and defensive positions

### Only Attacker Pays
- Defender can react (counter-spend) rather than commit upfront
- Asymmetric costs reward smart positioning
- Encourages aggressive expansion but punishes poor decisions

### Power Difference Auto-Takeover (diff > 3)
- Instant conquests feel rewarding (massive advantage)
- Close fights (1-3 diff) are tense and tactical
- Prevents "grinding" low-power battles

### Unlimited Incremental Upgrades
- No hard level caps on buildings
- Economy scales toward 30-min endgame (exponential growth slows down due to maintenance)
- Players with 8 hexes + Economy buildings can sustain constant upgrades
- Prevents "tech walls" where players get stuck

---

## Open Questions / TODO

### Playtesting Needs

- [ ] **Battle timing:** Does 6-14 sec range feel right? Too long/short?
- [ ] **Counter-spend balance:** Is 50 gold/sec too cheap or too expensive?
- [ ] **Conquest threshold:** Does 70% + 20 sec endgame feel tense?
- [ ] **Tech progression:** Should Tech costs scale (exponential) to slow late-game tech rush?
- [ ] **Map size impact:** How many hexes total? 30? 50? 100?

### Mechanical Unknowns

- [ ] **Player starting position:** Do all players start equidistant? Random corners?
- [ ] **Neutral hex power:** Do unclaimed hexes have Power 0 or random?
- [ ] **Building demolish refund:** Is 50% refund fair or should it be 100%?
- [ ] **Defender escape option:** Should there be a "Retreat" button to lose hex but save 50% gold?
- [ ] **Multiple battles:** Can hex be attacked by multiple enemies simultaneously?

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

