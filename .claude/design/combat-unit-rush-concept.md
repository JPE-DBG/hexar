# Combat Direction: Action Bar + Auto-Units + Hex Grid

Generated: 2026-06-13. User's breakthrough concept — most promising direction found so far.
Inspired by: Antiyoy (hex unit movement) + Clash Royale (action bar, card play, auto-march).

---

## Core Concept

**Action bar** (0–10) fills automatically over time. Player spends bar points to play **cards** that spawn **units**. Units auto-move hex by hex toward their target. Units fight automatically when they meet enemy units. **Buildings** are static structures that modify bar fill rate (economy) or provide defense (towers).

This solves:
- ✅ No arms race — units are consumed when fighting, not permanent
- ✅ Economic advantage is direct — better economy = faster bar = more cards = more units
- ✅ Readable — units move one hex at a time, nothing is instant
- ✅ Natural chaos cap — bar must refill, limits simultaneous unit density
- ✅ Area defense — towers cover a radius, not a single hex
- ✅ Roguelike — unit types, card pool, building types all vary per run
- ✅ Meaningful early game — first Settler card plays within seconds of game start
- ✅ Scales to multiplayer — each player has their own bar

---

## Established Cards (Starting Kit)

| Card | Cost | Effect |
|------|------|--------|
| Settler | 2 | Spawns unit that marches from capital toward nearest unclaimed hex and claims it |
| Soldier | 4 | Spawns unit that marches toward enemy capital, capturing owned hexes along the path |
| Builder | 3 | Places an economy building on a target owned hex (improves bar fill rate) |

---

## Open Design Questions (Must Resolve Per Variant)

1. **Spawn point**: units always from capital, OR from any owned border hex?
2. **Unit count cap**: max simultaneous units per player? (chaos prevention)
3. **Unit health**: 1-hit or multi-hit? Do stronger units beat weaker ones or destroy each other?
4. **Targeting**: units auto-target (dumb rush) or player sets target on card play?
5. **Tower placement**: any hex, border hexes only, or purchased separately?
6. **Bar fill rate**: fixed, OR depends on economy buildings, OR scales with hex count?
7. **What happens to a hex's buildings when a unit captures it?** Destroyed? Kept?
8. **Can units be recalled or redirected after spawning?**
9. **Multiple unit types from the start, or unlocked via roguelike?**

---

## Variants to Explore

See `combat-unit-rush-variants.md` (next file).
