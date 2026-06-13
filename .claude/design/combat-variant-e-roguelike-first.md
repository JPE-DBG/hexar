# Combat Variant E: Roguelike-First (Zero Base Kit)
## Full Design Exploration

Generated: 2026-06-13. Last reviewed: 2026-06-13.
Parent concept: `combat-unit-rush-concept.md`

> ⚠️ DRAFT — Nothing in this document is approved or finalized.
> Every mechanic, card, number, and rule is a first-pass proposal requiring explicit approval before implementation.
> Approved changes will be marked with ✅. Open questions marked with ❓.

---

## One-Line Pitch

Players start with only one card (Settler) ❓ and build their entire combat kit through roguelike research picks. Every game feels different because the kit is assembled from scratch.

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
bar_fill_rate = (owned hexes × base_per_hex)
              + (economy buildings × economy_bonus)
              − (tower maintenance per level)
```

- **Every owned hex**: +base/sec regardless of building on it
- **Economy building**: additional +bonus/sec on that hex (amplifies it)
- **Tower (any level)**: ongoing maintenance drain from bar rate
- Numbers (base_per_hex, economy_bonus, maintenance) deferred until unit mechanics settled

**Key property:** towers are net-negative on bar rate. Every tower built slows unit production. Defense costs offense. This is the anti-stalemate mechanism — a player who builds 15 towers has almost no bar to send soldiers.

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
- Unit cap: no hard cap initially — revisit if playtesting shows chaos ✅

### Bar Fill Rate ✅ (resolved above)

See Bar Fill Rate section. Formula: `hexes × base_per_hex + economy buildings × bonus − tower maintenance`. Numbers deferred.

---

## Starting State (Turn 0) ✅

**Cards in kit:** 4 — always available from second 0, no research required

| Card | Cost | Effect |
|------|------|--------|
| **Settler** | 2 bar | Unit marches to nearest unclaimed hex, claims it, dissolves |
| **Basic Soldier** | 4 bar | Unit marches toward enemy capital, captures hexes en route, fights on contact (1 HP) |
| **Gold Mine** | TBD | Places Economy building on owned hex — boosts bar fill rate |
| **Research Lab** | TBD | Places Research building on owned hex — generates research points passively |

Both players start identical. First direct interaction possible from the first seconds of the game (Basic Soldier available immediately).

**Economy and Research buildings:** always available — no research required to place them.

**Towers and army unit upgrades:** require research to unlock (see Research System below).

**Bar fill rate:** see Bar Fill Rate section above ✅. Exact numbers (base_per_hex, economy_bonus) deferred to playtesting.

---

## Research System ✅

### How research works ✅
- **Research buildings generate research points** passively (per second)
- More Research Labs built = faster point accumulation
- When a point threshold is reached: player sees **3 random stat bonuses**, picks 1 permanently
- Picks are **unlimited** — player controls the pace entirely through Research Lab investment

### Roguelike structure ✅
- **Cards auto-unlock** at fixed research thresholds — deterministic, no choice involved
- **Roguelike picks are stat bonuses only** (not cards)
- Bonuses can improve any aspect: economy, army stats, bar fill, towers, research speed
- Example bonuses: "+1 Soldier HP", "+20% bar fill rate", "+1 tower attack range", "+25% research points/sec"

### Card unlock progression ✅ (structure confirmed, order = initial draft)
Towers and army unit upgrades unlock automatically as research thresholds are crossed. Economy buildings (Gold Mine) and Research Lab are always available and do NOT require research.

Initial draft unlock order (needs full brainstorm before finalizing):
| Approx threshold | Card | Type |
|-----------------|------|------|
| 1 (early) | Watchtower (Tower L1) | Building |
| 2 | Raider | Unit — targets Economy buildings |
| 3 | Tower L2 (Archer) | Building upgrade |
| 4 | Squad | Unit — replaces Basic Soldier |
| 5 | Tower L3 (Ballista) | Building upgrade |
| 6 | Assassin | Unit — bypasses towers |
| 7 | Siege Engine | Unit — destroys towers |
| 8 (late) | Tower L4 (Fortress) | Building upgrade |

**⚠️ This card order is an initial draft — needs dedicated brainstorm session before implementation.**

### Thresholds ❓
- Exact point thresholds deferred — set after bar rate numbers and bonus pool finalized
- Model: increasing cost per pick (each pick requires more points than previous)

---

## Card Library ⚠️ INITIAL DRAFT — NEEDS COMPLETE REDESIGN

> The tier-based card system below is **obsolete**. The new structure is:
> - 4 starting cards (always available): Settler, Basic Soldier, Gold Mine, Research Lab
> - Remaining cards auto-unlock through research thresholds (no choice)
> - Roguelike picks are **stat bonuses**, not cards
>
> The card list below is preserved as a brainstorm reference only. Nothing in it is confirmed.

---

### Bonus Pool ⚠️ INITIAL DRAFT — needs confirmation

Each research pick shows 3 random bonuses from this pool. Player picks 1.

**Economy:**
- "Efficiency": Economy buildings +20% bar fill
- "Territory Yield": Each owned hex +10% bar contribution
- "Quick Build": Gold Mine and Research Lab cost 1 less bar

**Research:**
- "Accelerate": Research buildings generate 25% more points/sec
- "Deep Focus": First research building placed generates double points

**Army:**
- "Hardened": All units +1 HP
- "Swift": All units move 25% faster
- "Raider Mastery": Raider destroys 2 economy buildings per contact instead of 1
- "Siege Expert": Siege Engine takes 25% less damage

**Tower:**
- "Alert": Guarding Soldier chase leash +2 hexes
- "Sniper": All towers +1 attack range
- "Reinforced": All towers +1 HP
- "Tower Network": Towers share vision — unit spotted by one is targeted by all in range

---

### Old Tier System (obsolete — for brainstorm reference only)

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

**No HP reset** — capital damage is permanent. A failed push still weakens the capital; grinding it down over multiple pushes is a valid strategy.

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
