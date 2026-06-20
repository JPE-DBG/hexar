# Combat Variant E: Roguelike-First (Zero Base Kit)
## Full Design Exploration

Generated: 2026-06-13. Last reviewed: 2026-06-13.
Parent concept: `combat-unit-rush-concept.md`

> ⚠️ DRAFT — Nothing in this document is approved or finalized.
> Every mechanic, card, number, and rule is a first-pass proposal requiring explicit approval before implementation.
> Approved changes will be marked with ✅. Open questions marked with ❓.

---

## One-Line Pitch

Players start with 4 cards and unlock a complete combat kit through fixed research thresholds. Roguelike picks (stat bonuses + behavioral evolutions) make every game play differently within the same card set.

---

## Core Mechanic Summary

| Element | Description |
|---------|-------------|
| **Action bar** | 0–10, fills automatically. Rate = hex count + economy buildings − tower maintenance. |
| **Cards** | Spend bar to play a card → spawns a unit or places/upgrades a building. |
| **Units** | Auto-march hex by hex toward their target. Fight enemy units on contact. |
| **Buildings** | Economy (boost bar fill), Research (generate research points), Tower (defend area). |
| **Research** | Research buildings accumulate points → threshold reached → pick 1 of 3 roguelike cards. |
| **Win** | A unit reaches and enters the enemy capital hex ❓. |

---

## Bar Fill Rate ✅

Bar fills from **all three sources combined**:

```
bar_fill_rate = (owned hexes × 0.1/sec)
              + (economy buildings × 0.3/sec)
              − (tower maintenance per level)
```

| Source | Rate |
|--------|------|
| Each owned hex | +0.1/sec |
| Gold Mine | +0.3/sec additional |
| Tower L1 | −0.1/sec drain |
| Tower L2 | −0.2/sec drain |
| Tower L3 | −0.3/sec drain |
| Tower L4 | −0.5/sec drain |

**Representative values:**
- 5 hexes: 0.5/sec → Basic Soldier every 8s
- 10 hexes + 2 Gold Mines: 1.6/sec → Basic Soldier every 2.5s (unit cap binding)
- 20 hexes + 5 Gold Mines + 3 Tower L2s: 2.0 + 1.5 − 0.6 = 2.9/sec

**Key property:** towers are net-negative on bar rate. Defense costs offense. A player with 5 Tower L4s loses 2.5/sec — their offensive output is severely limited. This is the anti-stalemate mechanism.

---

## Buildings ✅ (behavior confirmed, costs/numbers deferred)

### Economy Building
- Placed on an owned hex via card play
- Increases that hex's bar fill contribution
- One per hex; stays until hex is captured (enemy unit reaching it destroys the building)

### Research Building  
- Placed on an owned hex via card play
- Generates research points passively over time
- More research buildings = faster roguelike pick unlocks
- One per hex

### Tower — Four Levels ✅

| Level | Name | Chase leash | Attack range | Total coverage | Counter |
|-------|------|-------------|--------------|----------------|---------|
| L1 | Guarding Soldier | 4 hexes | 0 (melee) | 4 hexes | 1–2 soldiers |
| L2 | Archer | 2 hexes | 2 hexes | 4 hexes (same as L1) | 2–3 soldiers |
| L3 | Ballista | ❓ | 3 hexes | ❓ | 3–4 soldiers or siege |
| L4 | Fortress | 0 (immovable) | ❓ high | high | Siege unit only |

**L1 — Guarding Soldier:**
- Chases approaching enemy units up to **4 hexes** from home hex
- Melee only — must reach the enemy to fight
- Returns to home hex when enemy retreats beyond leash or is destroyed
- 1 HP base (dies in one fair fight); roguelike can upgrade HP

**L2 — Archer:**
- Moves up to **2 hexes** from home to reposition
- Fires shots at units up to **2 hexes away** (attack range)
- Total effective coverage = 4 hexes (same as Guarding Soldier, but from a distance)
- Safer than L1: fires before enemy reaches it; L1 must close the gap
- 1 HP base; roguelike can upgrade HP and/or range

**L3 — Ballista:**
- Chase leash and exact range ❓ — to be defined after L1/L2 playtesting
- Stronger shots than Archer; slower fire rate ❓

**L4 — Fortress:**
- Cannot move (0 chase leash)
- High range compensates for immobility ❓
- Very high HP — regular soldiers highly inefficient against it
- Requires Siege unit to counter effectively

**Tower upgrade:** each level is an upgrade of the previous. Costs increasing bar to upgrade. Tower stays on same hex.

**Tower maintenance cost:** drains bar rate per tick. Scales with level. **Values deferred** — set after unit mechanics and bar numbers finalized.

**Siege unit implication:** L4 Fortress requires Siege unit as counter → must be unlockable via roguelike research. Strategic read: if opponent turtles Fortresses, research Siege.

---

## Action Bar

- Scale: 0–10
- Unit cap: **3 simultaneous units per player** ✅ — prevents bar overflow at full expansion; bar fills but no unit spawns if 3 are active

### Bar Fill Rate ✅ (resolved above)

See Bar Fill Rate section. Formula: `hexes × 0.1 + economy buildings × 0.3 − tower maintenance`.

---

## Starting State (Turn 0) ✅

**Cards in kit:** 4 — always available from second 0, no research required

| Card | Cost | Effect |
|------|------|--------|
| **Settler** | 2 bar | Unit marches to nearest unclaimed hex, claims it, dissolves |
| **Basic Soldier** | 4 bar | 1 HP, marches toward enemy capital, captures hexes en route, fights on contact |
| **Gold Mine** | 5 bar | Places Economy building — +0.3/sec bar fill (16.7s payback — real tradeoff vs Soldiers) |
| **Research Lab** | 3 bar | Places Research building — generates research points passively |

Both players start identical. First direct interaction possible from the first seconds (Basic Soldier available immediately).

**Economy and Research buildings:** always available — no research required to place them.

**Towers and army unit upgrades:** require research to unlock (see Research System below).

**Bar fill rate:** see Bar Fill Rate section ✅.

---

## Research System ✅

### How research works ✅
- **Research buildings generate research points** passively (per second)
- More Research Labs built = faster point accumulation
- Two types of events when threshold is crossed:
  1. **Card unlock** (deterministic) — a new card becomes available to all players at fixed thresholds
  2. **Roguelike pick** (player choice) — player sees 3 random stat bonuses/evolutions, picks 1 permanently
- Picks are **unlimited** — player controls pace through Research Lab investment

### Roguelike structure ✅
- **Cards auto-unlock** at fixed research thresholds — deterministic, same for all players, no choice
- **Roguelike picks are stat bonuses and behavioral evolutions** (not cards)
- Picks differentiate players: same card pool, different power modifiers

### Card unlock progression ✅

| Threshold | Card Unlocked | Bar Cost | Category |
|-----------|--------------|----------|----------|
| 1 | **Scout** | 1 | Unit — cycle/pressure |
| 2 | **Tower L1 (Watchtower)** | 3 | Building — first defense |
| 3 | **Raider** | 3 | Unit — economy disruptor |
| 4 | **Tower L2 (Archer)** | 4 | Building — ranged defense (requires L1) |
| 5 | **Heavy Soldier** | 6 | Unit — tank win condition (2 HP) |
| 6 | **Catapult** | 5 | Unit — splash support |
| 7 | **Tower L3 (Ballista)** | 5 | Building — strong defense (requires L2) |
| 8 | **Assassin** | 5 | Unit — tower-bypass win condition |
| 9 | **Siege Engine** | 7 | Unit — tower destroyer |
| 10 | **Tower L4 (Fortress)** | 6 | Building — ultimate defense (requires L3) |

Siege Engine (threshold 9) unlocks before Tower L4 (threshold 10) — counter always available before ultimate defense.

### Thresholds ❓
- Exact point values deferred — set after implementation
- Model: increasing cost per pick (each pick requires more points than previous)

---

## Unit Types ✅

### All Bar Costs

| Card | Bar | HP | Speed | Behavior |
|------|-----|----|-------|---------|
| Scout | 1 | 0 (no combat) | 1 hex/sec | Claims 1 unclaimed hex toward opponent, dissolves |
| Settler | 2 | 0 (no combat) | 0.5 hex/sec | Claims nearest unclaimed hex, dissolves |
| Basic Soldier | 4 | 1 | 0.5 hex/sec | Marches to capital |
| Raider | 3 | 1 | 0.5 hex/sec | Targets Economy buildings first, then capital |
| Heavy Soldier | 6 | 2 | 0.5 hex/sec | Marches to capital, survives 2 tower hits |
| Catapult | 5 | 1 | 0.3 hex/sec | Advances to capital; on death: 1 damage to all units/towers in 1-hex radius |
| Assassin | 5 | 1 | 0.5 hex/sec | Ignores towers (doesn't trigger them), goes to capital |
| Siege Engine | 7 | 3 | 0.3 hex/sec | Targets towers first, takes 0 damage while attacking a tower |

### Unit Target Priority

| Unit | Priority 1 | Priority 2 | Priority 3 |
|------|-----------|-----------|-----------|
| Basic Soldier | Enemy units | Towers | Capital |
| Heavy Soldier | Enemy units | Towers | Capital |
| Scout | Nearest unclaimed hex toward opponent | Dissolves on claim | — |
| Raider | Economy buildings | Research Labs | Capital |
| Assassin | Capital directly | (ignores towers and enemy units) | — |
| Catapult | Advances to capital | Explodes on death (1-hex splash) | — |
| Siege Engine | Towers | Then normal advance | Capital |

### Counter Triangle

```
3× Basic Soldiers (swarm)  →  countered by  →  Catapult (splash kills cluster)
Catapult (slow, 1 HP)      →  countered by  →  Heavy Soldier (survives, kills before arrival)
Heavy Soldier (tank)       →  countered by  →  Tower L2+ (ranged fire kills 2-HP unit)
Tower (static defense)     →  countered by  →  Siege Engine (destroys towers)
Siege Engine (slow)        →  countered by  →  Basic Soldiers (fast, cheap, overwhelm)
                                                    ↕ loops back to Catapult countering swarm
```

---

## Roguelike Pick Pool ✅

### Behavioral Evolutions (behavior changes, not stat boosts)

| Pick | Effect |
|------|--------|
| "Raider Network" | Raider now targets ALL buildings (Research Labs + Economy) |
| "Ghost Protocol" | Assassin ignores enemy units too (full ghost — passes through everything) |
| "Siege Armor" | Siege Engine takes 0 damage from towers it is currently attacking |
| "Tower Refund" | When one of your towers kills an enemy unit → +1 bar (**defensive win path**) |
| "Pack Tactics" | Basic Soldiers gain +1 HP when 3 are active simultaneously |

### Stat Bonuses

| Pick | Effect |
|------|--------|
| "Hardened" | All units +1 HP |
| "Swift" | All units +25% movement speed |
| "Sniper" | All towers +1 attack range |
| "Reinforced" | All towers +1 HP |
| "Alert" | Guarding Soldier chase leash +2 hexes |
| "Efficiency" | Gold Mine +25% bar fill (0.3 → 0.375/sec) |
| "Accelerate" | Research Labs +25% points/sec |

**Tower Refund** is the defensive win path: kill enemy unit with tower → gain +1 bar → counter-push while opponent's bar regenerates.

---

## Archetypes ✅

| Archetype | Core Strategy | Key Cards | Key Picks |
|-----------|--------------|-----------|-----------|
| **Aggressor** | Constant Soldier pressure + Raider economy destruction | Soldier, Raider, Assassin | Swift, Raider Network |
| **Turtle-Pusher** | Tower line → Tower Refund farm → Siege breach | Towers, Siege Engine | Tower Refund, Siege Armor |
| **Economist** | Gold Mine first, delayed but powerful army | Gold Mine × 3, Heavy Soldier | Efficiency, Hardened |
| **Rusher** | Scout cycle + 3× Soldiers before opponent builds towers | Scout, Basic Soldier | Pack Tactics, Swift |

---

## Card Library ✅

> The tier-based card system below this section is **obsolete**. Current structure:
> - 4 starting cards (always available): Settler, Basic Soldier, Gold Mine, Research Lab
> - 10 cards auto-unlock through research thresholds (see Research System above)
> - Roguelike picks are **stat bonuses + behavioral evolutions** (see Roguelike Pick Pool above)

---

### Bonus Pool ✅ (see Roguelike Pick Pool section above)

---

### Old Tier System (obsolete — preserved for reference only, do not use)

#### Tier 0 (Always Available)
| Card | Cost | Effect |
|------|------|--------|
| Settler | 2 | Unit marches to nearest unclaimed hex, claims it |

### Tier 1 (First Pick — Sets Archetype)
| Card | Cost | Effect |
|------|------|--------|
| **Soldier** | 4 | Unit marches toward enemy capital, captures enemy hexes en route |
| **Gold Mine** | 3 | Places Economy building on target owned hex (+0.5/sec bar fill) |
| **Watchtower** | 3 | Places Tower on target owned hex (shoots units within 2 hexes, auto-aim) |

This first pick signals the player's early archetype: Aggressor, Economist, or Defender.

### Tier 2 (Second + Third Pick)
| Card | Cost | Archetype | Effect |
|------|------|-----------|--------|
| Heavy Soldier | 5 | Aggressor | Soldier variant; survives 2 tower hits instead of 1 |
| Raider | 4 | Aggressor | Unit targets Economy buildings first; destroys one on contact, then continues |
| Squad | 5 | Aggressor | Spawns 2 light Soldiers (each dies to 1 tower hit) simultaneously |
| Advanced Mine | 4 | Economist | Upgrades one existing Gold Mine to +1.0/sec (replaces +0.5) |
| Foundry | 4 | Economist | Economy building that reduces all unit spawn costs by 1 bar |
| Upgraded Tower | — | Defender | Upgrade an existing Tower: range 3 instead of 2 |
| Garrison | 3 | Defender | Spawns a stationary defender unit on target hex; fights any unit that enters |

### Tier 3 (Fourth + Fifth Pick)
| Card | Cost | Archetype | Effect |
|------|------|-----------|--------|
| Commander | 6 | Aggressor | Unit that buffs all friendly units within 2 hexes (+1 effective health) |
| Assassin | 6 | Aggressor | Ignores towers entirely; targets capital directly; fragile (1 hit from any unit) |
| War Factory | 5 | Economist | Economy building: generates 1 free Soldier per 30s automatically |
| Tower Network | — | Defender | Upgrade: all your towers gain shared vision — a unit spotted by one tower is targeted by all in range |
| Fortress | 4 | Defender | Reinforces one of your hexes: takes 3 units to capture instead of 1 |

### Tier 4 (Sixth Pick Onward — Late Game)
| Card | Cost | Effect |
|------|------|--------|
| Blitz | 7 | Spawns 3 Soldiers simultaneously toward capital — overwhelming push |
| Economic Sabotage | 5 | Raider variant that targets all Economy buildings in a path, not just one |
| Siege Engine | 7 | Slow heavy unit; destroys towers in one hit; takes 3 unit hits to kill |
| Rally | 0 | All currently active friendly units gain +1 health and speed up for 10s |

---

## Early Game (0–5 min) ❓

**What happens:**
Both players use Settler repeatedly to claim unclaimed hexes. Bar fills slowly (0.2/sec base). A Settler can be sent roughly every 10 seconds at base fill rate. Map fills with owned hexes quickly in all directions. When both Settlers head for the same unclaimed hex, first arrival wins — this is the only player interaction in this phase.

**Decisions:**
- Which unclaimed hexes to prioritize (hexes that expand toward opponent vs hexes that build a wider economy base)
- When to spend bar on Settler vs save for the first research pick's card

**Research Pick 1 (~2 min, at 5 hexes owned):**
This is the most important moment of the game. Three options appear:
- Pick Soldier → first combat capability, but no economy yet; Soldiers cost 4 bar and bar fills at 0.2/sec (takes 20 seconds per Soldier)
- Pick Gold Mine → no combat yet; but bar will fill much faster, enabling more Settlers and future cards
- Pick Watchtower → defensive positioning; can't attack but can hold territory

**Player interaction starts:** If both players pick Soldier, the first unit exchanges happen by minute 3. If one picks Gold Mine and the other picks Soldier, the Soldier player gets early pressure but the economic player will outpace them later.

---

## Mid Game (5–15 min) ❓

**What happens:**
2–4 research picks have fired. Each player has a kit of 3–5 cards. Bar fill rates diverge significantly based on economy investment. Borders form and are contested by Soldiers and Towers. Raiders start destroying Economy buildings, creating economic disruption. Players must both expand (Settlers) and defend (Towers/Garrisons) and attack (Soldiers).

**Key interactions:**
- Soldier vs Soldier: both die (equal), first arrival to enemy territory is a threat
- Soldier vs Tower: tower wins (Soldier dies), unless Heavy Soldier or Siege Engine
- Raider reaches opponent's Gold Mine: mine is destroyed, opponent's bar fill drops — very impactful moment
- Garrison on a hex: that hex requires a unit to fight through, slowing any advance

**Decisions per minute:**
- Which card to play with current bar
- Where to place new buildings (closer to border for Towers; closer to economy cluster for Mines)
- Whether to send Soldier toward capital now or wait for a better combo

**Research Pick 3 (~8 min, at 13 hexes):**
Tier 2 or 3 cards available. Player's archetype is now clear. The pick either deepens the archetype (Heavy Soldier on top of Soldier) or pivots (Garrison on top of Soldier — becoming more defensive).

---

## Late Game (15–30 min) ❓

**What happens:**
Both players have 5–6 card kits. Bar fill rates are high (multiple Gold Mines built). Units are more powerful (Tier 3–4). Capital pushes become real threats. Fortress hexes guard the approach to capitals. Blitz or Siege Engine cards available for decisive moments.

**The endgame pattern:**
Player A saves bar up to 10, plays Blitz (spawns 3 Soldiers simultaneously). Player B sees 3 units incoming and has ~20 seconds to respond — Tower shots reduce them, Garrison stops one, but the third reaches capital.

OR: Player A uses Assassin (ignores towers) pointed at capital while Player B's attention is on a different front.

**Why games end before 30 min:**
- Economic dominance → fast bar → more units more often → slow territorial gain → eventually capital exposed
- A decisive Tier 4 card played at the right moment breaks through
- Raider destroying all of opponent's Gold Mines collapses their bar rate → they can't respond to attacks

**Capital Defense:**
Capital hex has **2 built-in health** (takes 2 units to capture). This prevents a single Soldier from ending the game in minute 5. First unit to reach capital reduces it to 1 health; second unit captures it. Defender has one unit-travel-time to intercept the second unit.

---

## Win Condition ✅ (approach confirmed, HP value TBD)

A unit enters the enemy capital hex → unit dissolves, capital loses 1 HP. When capital HP reaches 0, the game ends immediately.

**Capital HP:** TBD — initial draft = 2 (requires 2 units to reach the capital). Subject to playtesting.

**No HP reset** — capital damage is permanent. A failed push still weakens the capital; grinding it down over multiple pushes is a valid strategy. Reaserch can add option to heal capital.

**Future:** Capital may have a built-in upgradeable tower as a defensive option — deferred, not in initial implementation.

---

## Player Interactions ❓

| Phase | Type | Frequency |
|-------|------|-----------|
| Early game | Settler racing for unclaimed hexes | Every 10–20 seconds |
| Mid game | Unit vs unit combat at borders | Every 30–60 seconds |
| Mid game | Raider destroying economy buildings | Occasional, high impact |
| Mid game | Tower shooting incoming units | Continuous passive |
| Late game | Capital push attempts | 1–3 per game, high drama |
| Throughout | Research picks (indirect) | Every 3–5 min |

**Max simultaneous "battles":** 3 (unit cap per player). Map never has more than 6 units total in 1v1. Readable.

---

## Roguelike Feel Per Archetype ❓

**Aggressor run:** Soldier → Heavy Soldier → Squad → Blitz. Constant pressure. Bar must be good enough to afford frequent Soldiers. Games end fast — either you rush them down or you run out of steam.

**Economist run:** Gold Mine → Advanced Mine → Foundry → War Factory. Slow early game, dominant mid-late. Bar fills so fast that Settlers and Soldiers play non-stop. Opponent must destroy economy buildings before this snowballs.

**Defender run:** Watchtower → Upgraded Tower → Tower Network → Fortress. No offensive pressure early. Builds an impenetrable turtle. Must eventually pivot to attack via Raider or Assassin — or wins by opponent running out of units attacking towers.

**Hybrid run (most common):** Gold Mine → Soldier → Garrison → Heavy Soldier. Balanced. Less extreme than the above. More reactive to opponent's strategy.

---

## Known Risks and Mitigations ❓

| Risk | Mitigation |
|------|-----------|
| Early game is parallel solitaire | Settlers racing for same hexes creates real competition. Research Pick 1 at 2 min brings first divergence fast. |
| Economy snowball (rich player always wins) | Raiders exist specifically to destroy economy buildings. An Aggressor-first player can cripple an Economist before they reach late game. |
| Turtle problem (Defender never loses) | Assassin card ignores towers entirely. Siege Engine destroys towers. Defender must eventually be beatable. |
| Unit cap (3) feels too limiting | Playtest: 3 may be right for readability. Can tune to 4. Key is that player never feels like they're losing units they didn't see. |
| Research pacing too slow | Territory-based pacing rewards fast Settler expansion. If research feels slow, lower the hex thresholds. |
| Capital captured too fast in early game | 2-hit capital health + 30-second reset gives defender enough time to respond to a first hit. |
| No defense option in first 5 minutes | Settler claiming hexes around capital creates a buffer zone. First enemy unit has to travel through several hexes before reaching capital — traversal time IS the early defense. |

---

## Open Questions for Prototyping

1. Does bar fill from total hexes OWNED or from Economy buildings only? (Current: buildings. Alternative: both contribute.)
2. What speed do units move? 1 hex/sec? 0.5 hex/sec? This defines how much reaction time the defender has.
3. Can a player recall a unit after spawning it? (Probably no — simplicity.)
4. When a Soldier captures an enemy hex en route, does the enemy building get destroyed? (Probably yes — it's now your hex with no building.)
5. Can the player place a Tower on a hex that already has a Gold Mine? (Probably no — one building per hex.)
6. Do units from the same player block each other, or can two Soldiers occupy the same hex? (Probably stack on same hex, fight as a group.)
