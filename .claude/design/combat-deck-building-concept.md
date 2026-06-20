# Combat: Deck-Building Concept
## Full Design Specification

Generated: 2026-06-20. Last reviewed: 2026-06-20.

> ⚠️ DRAFT — Nothing in this document is approved or finalized.
> Approved items marked ✅. Open questions marked ❓.

---

## One-Line Pitch

Players build a cycling deck of cards in real-time. Bar fills automatically and is spent to play cards or buy new ones from a shared shop. Units march to the enemy capital on a hex map. Deck composition is the primary strategic axis.

---

## Core Loop

1. Bar fills at a fixed rate
2. Top card of your deck is always visible — pay its cost to play it, or push it to an aside slot
3. Played cards go to the bottom of the deck and cycle back
4. Spend bar in the shared shop to add new cards to the bottom of your deck
5. Units you deploy march across the hex map toward the enemy capital
6. First player to reduce the enemy capital to 0 HP wins

---

## Bar ✅

- **Base fill rate:** 0.1 bar/sec (fixed — does not change with territory)
- **Scale:** 0–10 (presumed; TBD if cap changes)
- **Single currency:** bar pays for everything — playing cards, buying from shop, removing cards
- Bar does not reset or drain; it accumulates until spent or capped

---

## Deck Mechanics ✅

### Playing cards
- The top card of your deck is always visible
- Pay its bar cost → effect resolves → card goes to the bottom of the deck
- All cards cycle back regardless of type (including building cards)
- Playing a building card again places another copy of that building on a different hex

### Aside slot
- Push the current top card sideways into the aside slot
- The slot is blocked until you pay the aside card's cost and play it
- Played aside card goes to the bottom of the deck
- You cannot push another card into a blocked aside slot
- Starting aside slots: **1**
- Additional aside slots: purchasable via a building card from the shop (requires free hex)

### Deck lockout prevention
- If the top card cannot be played (insufficient bar, or no free hex for a building card), it auto-moves to the bottom after a fixed timer ❓ (value TBD — single constant, easy to tune)
- This ensures the deck keeps moving even in edge cases

### Deck size
- Starting deck: 7 cards
- No maximum size — buying cards from the shop grows the deck
- Larger deck = each card appears less frequently (trade-off: power vs cycle speed)
- Card removal available in shop (see Shop section)

---

## Starting Deck ✅ (7 cards)

| Card | Cost | Effect | Quantity |
|------|------|--------|----------|
| Hex Claim | 1 bar | Claim one adjacent unclaimed hex (player chooses which) | ×2 |
| +1 Bar Burst | 0 bar | Instantly add +1 to current bar | ×2 |
| +15% Bar Speed | 1 bar | Bar fill rate +15% for 15s (stacks with other active copies) | ×2 |
| Basic Soldier | 2 bar | Deploy a soldier unit toward enemy capital | ×1 |

**Notes:**
- The two 0-cost +1 Bar Burst cards are always worth playing whenever they appear — no decision required
- Multiple +15% Bar Speed cards active simultaneously stack (real-time game, fast deck cycling can keep multiple running)
- The single Basic Soldier creates early offensive interaction from the start of the game

---

## Map ✅

### Purpose
- **Distance** = reaction time. Units take time to cross the map — this is the defender's window to respond.
- **Claimed hexes** = building slots. No free claimed hex = building cards are unplayable (auto-cycle).

### Hex ownership
- Each player starts with their capital hex claimed
- Hex Claim cards expand territory to adjacent unclaimed hexes (player chooses which)
- Players keep hex ownership even if a building on that hex is destroyed
- A destroyed building frees the hex for a new building — the hex itself is never lost in MVP

### Unit movement
- Soldiers ignore non-combat buildings and march directly toward the enemy capital
- Soldiers fight enemy soldiers and towers they encounter en route
- Soldiers do not capture hexes they pass through (MVP simplification)

---

## Units ✅

### Basic Soldier

| Stat | Value |
|------|-------|
| Bar cost | 2 |
| HP | 2 |
| Attack power | 1 |
| Attack speed | 1 hit/sec |

### Combat resolution

**Soldier vs Soldier (1v1):**
- Both attack simultaneously, once per second
- Equal soldiers (2 HP, 1 atk): each deals 1 damage/sec — both die after 2 seconds

**Soldier vs Soldier (2 attackers vs 1 defender):**
- Defender takes 2 damage/sec (from both attackers) → dies in 1 second
- Defender attacks one attacker: 1 damage dealt before dying
- Result: one attacker at 1 HP, one at 2 HP — both survive and continue to capital

**Soldier vs Tower:**
- Tower fires at soldiers within its range (stats TBD — tower is a shop card)
- Soldier attacks tower on contact
- Higher-HP soldiers or multiple soldiers can clear a tower; fragile soldiers may die to tower fire before reaching it

### Capital damage
- Each soldier that reaches the enemy capital hex deals damage to it ❓ (amount TBD — likely 1 HP per soldier, subject to playtesting)
- Soldier combat with capitol like with tower, 
- when new unit spawn from capitol and its under attack, spawned unit is under attack instead of capitol

---

## Win Condition ✅

- **Capital HP:** 20 (subject to playtesting)
- **No built-in capital defense** — the capital is a pure HP counter
- When capital HP reaches 0, that player loses immediately
- Damage is permanent — no HP regeneration (potential card in shop to heal capitol)

---

## Shop ✅

### Structure
- Shared pool visible to both players
- Dominion-style: fixed set of available card types, each with a limited quantity
- Buy at any time by paying bar cost → new card goes to the bottom of your deck
- Both players draw from the same pool — first to buy gets the card; competing for limited copies is a strategic layer

### Buying trade-offs
- Buying grows your deck → each card appears less often
- Buying a powerful card that costs high bar = slower cycle until bar refills
- Removing cards (see below) is the counter to an overgrown deck

### Card removal
- Available as a shop option (not a playable card)
- Removes the current **top card** or current **aside card** from your deck permanently
- Cost: mid-to-high bar ❓ (value TBD)
- Strategic use: remove starting Hex Claim cards once territory is established; remove weak filler cards to speed up cycle to power cards

---

## Shop Card Pool ❓ (not yet designed)

Card types confirmed for the shop. Exact cards, stats, costs, and quantities are TBD.

### Card types

**Units:**
- More unit types (faster, tankier, special behaviors) — TBD

**Buildings — Economy:**
- **Bar Speed Building:** permanent bar fill rate boost on the hex it's placed on ❓ (bonus value TBD)

**Buildings — Military:**
- **Tower:** fights soldiers passing within range; has HP and attack stats ❓ (all stats TBD)
- **Spawner:** produces a free unit of a given type every X seconds without using the deck; high bar cost ❓ (unit type, interval TBD)

**Buildings — Utility:**
- **Aside Slot Building:** permanently adds 1 aside slot; requires a free hex to place; played as a building card

**Deck Management:**
- Card removal (shop service, not a card) — removes top or aside card

---

## Open Questions ❓

| # | Question | Notes |
|---|----------|-------|
| 1 | Auto-cycle timer value | Single constant; TBD during playtesting |
| 2 | Tower stats (HP, range, fire rate) | TBD when tower confirmed as shop card |
| 3 | Spawner unit type and spawn interval | TBD when spawner confirmed as shop card |
| 4 | Bar speed building bonus amount | TBD |
| 5 | All shop card bar costs | TBD — shop design session needed |
| 6 | Shop card quantities (how many copies of each) | TBD |
| 7 | Damage per soldier hit on capital | Likely 1, needs playtesting |
| 8 | Bar scale cap (is 0–10 correct?) | TBD |
| 9 | Map size and capital distance | TBD — affects reaction time feel |
| 10 | +15% bar stacking: additive or multiplicative? | TBD (additive simpler: 2 stacks = +30%) |

---

## What This Design Removes vs Previous System

Compared to the Action Bar + Roguelike (Variant E) design:

| Removed | Why |
|---------|-----|
| Research resource (separate from bar) | Single resource (bar) is simpler |
| Research Labs as building type | No separate research track needed |
| Roguelike pick system (stat bonuses) | Deck composition creates variance instead |
| Territory-based bar fill rate | Fixed rate removes hex-count math complexity |
| Unit cap (3 simultaneous) | Deck cycling rate naturally limits unit spam |
| Tower maintenance drain | Towers no longer need to cost ongoing bar |
| Sequential tower upgrade gating | Tower is just a shop card, no prerequisite chain |

The core player action is now always the same: manage your deck, manage your bar, buy cards that strengthen your deck. Map and units are the expression of that investment.
