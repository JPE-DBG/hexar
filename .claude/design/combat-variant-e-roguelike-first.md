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
| **Action bar** | 0–10, fills automatically. Rate depends on economy buildings owned. |
| **Cards** | Spend bar to play a card → spawns a unit or places a building. |
| **Units** | Auto-march hex by hex toward their target. Fight enemy units on contact. |
| **Buildings** | Static structures: Economy (speeds bar fill) or Tower (shoots units in range). |
| **Research** | Triggered by owning N hexes. Each unlock = pick 1 of 3 new cards ❓. |
| **Win** | A unit reaches and enters the enemy capital hex ❓. |

---

## Action Bar

- Scale: 0–10
- Hard unit cap: **3 simultaneous units per player** (chaos prevention) ❓

### ❓ Bar Fill Rate — UNRESOLVED

The original proposal (+0.5/sec per Economy building) is broken at map scale.
On a 50+ hex map a player could own 20+ hexes with Economy buildings → bar fills in ~4s → units sent every 4s → game becomes unreadable spam.

**Three candidate models — needs decision:**

**Option A: Hard cap on bar fill rate**
- Economy buildings: +0.1/sec each, max fill rate 1.0/sec (hard cap)
- Bar fills in minimum 10s regardless of economy size
- Pro: simple ceiling; Con: economy investment loses value past ~8 buildings

**Option B: Economy buildings reduce unit costs, bar fills at fixed rate**
- Bar fills at fixed 0.5/sec (always 20s to full bar)
- Each Economy building: −0.2 bar cost on all units (floor: 1)
- With 5 buildings: Soldier costs 3 bar → sent every 12s at full bar
- Pro: economy remains meaningful throughout; Con: two levers may confuse

**Option C: Bar fills from hex count, not buildings**
- Each hex owned: +0.015/sec (50 hexes = 0.75/sec → fills in ~13s)
- Economy buildings serve a separate purpose (reduce unit cost or boost unit strength)
- Pro: territory directly = military power; Con: economy buildings lose identity

---

## Starting State (Turn 0)

**Cards in kit:** 2 ❓
- **Settler** (cost: 2): spawns a unit that marches from capital toward nearest unclaimed hex, claims it on arrival, then dissolves
- **Basic Soldier** (cost: 4): spawns a unit that marches toward enemy territory; captures hexes en route; fights enemy units on contact

Both players start identical with both cards. First direct interaction possible from the first seconds of the game.

**Buildings available:** none yet (unlocked via research) ❓

**Bar fill rate:** see unresolved section above ❓

---

## Research System ❓

Research points accumulate from **territory size** — owning more hexes progresses research faster. This ties exploration (Settlers) directly to kit development.

| Hexes owned | Research pick unlocked |
|-------------|----------------------|
| 5 | Pick 1 (Tier 1) |
| 9 | Pick 2 (Tier 1 or 2) |
| 13 | Pick 3 (Tier 2) |
| 17 | Pick 4 (Tier 2 or 3) |
| 21 | Pick 5 (Tier 3) |
| 25 | Pick 6 (Tier 3 or 4) |
| 29+ | Pick 7+ (Tier 4) |

Each pick: player sees **3 random cards from the next eligible tier**, chooses 1. The card is permanently added to their kit. Unpicked options are discarded.

The game knows the player's current kit and skews options toward synergistic picks (e.g. if player has Soldier, Tier 2 options bias toward Soldier upgrades, not Tower upgrades).

---

## Card Library ❓

### Tier 0 (Always Available)
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

## Win Condition ❓

A unit enters the enemy capital hex for the **second time** (or first time if capital health was already at 1 from a previous hit). Capital health resets to 2 after 30 seconds if not hit again — a failed push is survivable but sets the stage for the next attempt.

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
