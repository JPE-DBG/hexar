# Combat System Brainstorm — 10 Candidate Architectures

Generated: 2026-06-13. These are raw directions, not finalized designs.
Current system (Power gate) is being replaced. None of these assume Power buildings exist.

---

## 1. Front Line Push
One shared front line between territories. Both players spend gold to push it toward the enemy capital. Higher income = faster push. One contest, no chaos. Roguelike: burst push abilities, line anchors, income multipliers during pushes.

## 2. General Tokens
Each player has 1–2 named generals (tokens on map). Generals move hex by hex (~3s per step). Entering enemy territory captures the hex and continues marching until stopped by an enemy general (instant duel). No battles elsewhere. Roguelike: each research pick improves one general (speed, capture chain, duel strength).

## 3. Siege Clock
Declare a siege on adjacent enemy hex — free to start. Visible countdown (~20s). Defender must spend gold to "hold" before it expires or hex flips. Holding resets clock. Each active siege costs upkeep per second. Roguelike: extend/shorten clock, siege area effects, free first siege.

## 4. Charge and Release
No combat buildings. Each owned hex slowly charges a global "power meter." Release charge at adjacent enemy hex → 3-second flash contest, whoever released more charge wins. Defender can counter-release. Depletes a shared charge pool. Roguelike: charge speed, release patterns, chain reactions.

## 5. Mercenary Waves
Spend gold to hire a wave that spawns on a border hex and auto-marches toward enemy capital capturing hexes. Strength = gold spent, lasts ~15s, then dissolves. Waves cancel each other if equal strength. One active wave per player. Roguelike: wave types (fast/weak, slow/strong, splitting, flanking).

## 6. Ritual Challenge (Duel)
Challenge any adjacent enemy hex: 5-second countdown visible to both. Both players commit gold during those 5 seconds. Whoever committed more wins — both lose what they committed. Cap: 1 active challenge per player. Roguelike: reduce window, get gold back on win, hidden challenge reveal.

## 7. Resonance Spread
No attack action. Each hex radiates "presence" outward. When your presence exceeds the enemy's in a border hex for 10 consecutive seconds, it flips. Defense = spending gold to anchor key hexes. Roguelike: presence radius, anchor strength, burst presence events.

## 8. Supply Starvation
Surround enemy hexes — a hex with no path to enemy capital starts a 15s starvation timer then goes unclaimed. Claiming unclaimed hexes is free and instant. Game = cutting corridors, not direct combat. Roguelike: starvation speed, claim bonus, corridor blocking abilities.

## 9. Hero Abilities (Cooldown Combat)
No combat buildings. Hero on capital with 4 ability slots on cooldowns. Example: Claim (take hex, 15s CD), Shield (hex immune 10s, 25s CD), Raze (destroy building, 20s CD), Rally (2× gold 8s, 30s CD). Roguelike: research picks unlock new abilities or upgrade existing slots — different builds each game.

## 10. Momentum Streak
Capture is instant (50g, click adjacent enemy hex). Captured hex is "unstable" for 8s (opponent reclaims free). After each capture: 4s global cooldown before next capture. Rushing is self-defeating — can't defend a chain of fresh captures. Roguelike: shorten cooldown, extend instability window, make certain captures permanent.

---

## Exploration History
→ **Charge + Release × Hero/Card Deck × Roguelike** — 10 variants detailed in `combat-charge-card-variants.md`
→ **Action Bar + Auto-Units + Hex Grid** — core concept in `combat-unit-rush-concept.md`
→ **Variant E: Roguelike-First (Zero Base Kit)** — full design in `combat-variant-e-roguelike-first.md` ← current focus
