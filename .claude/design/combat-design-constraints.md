# Combat Design Constraints — Discovered During Brainstorming

Last updated: 2026-06-13

These constraints must be satisfied by any new combat system for Hexar.
They were discovered through playtesting and design discussion, not upfront assumptions.

---

## Hard Constraints (non-negotiable)

- **Real-time** — not turn-based. Players act simultaneously.
- **No chaos** — no system where many battles fire simultaneously and the player loses track of what's happening. Few meaningful events visible at once. 
- **Battles cannot be crazy fast** — outcomes must be readable in real-time. No instant multi-hex flips from a single action without warning.
- **Limited simultaneous attack/defense spots** — hard cap on how many contests can be active at once per player.
- **Capital capture = win condition** — game ends when capital is taken.
- **1v1 first, more players fightning dead match later** — design must scale to more players without breaking.

## Soft Constraints (strongly preferred)

- **Economic advantage should translate to military advantage** — the player with better economy should be able to win fights, not just build faster.
- **No Power buildings in current form** — the permanent/symmetric arms race that caused the original stalemate must not return.
- **Roguelike research** — research grants the player a pick-1-of-3 improvement after gaining certain reaserch points. Each game, a player's combat and resource kit should feel different from the last.
- **Meaningful early game** — players should interact (or at least make meaningful decisions) before minute 5.
- **Defense should matter but not be a permanent wall** — defense raises the cost/difficulty of attacking; it never makes attack impossible.

## Anti-patterns to Avoid

- Single-hex fortify/shield: protecting one hex is almost always useless because opponent attacks adjacent. Defensive abilities must protect areas or the player's whole kit.
- Symmetric arms races: if both players can always match each other's investment with the same investment, the game deadlocks.
- Interaction-free early game: if neither player can affect the other for the first 3–5 minutes, the game feels like parallel solitaire.
- Randomness that feels unfair: card draw RNG is acceptable if the pool is small and predictable; pure dice rolls are not.
- Cooldown-only systems with 3 abilities: too few options for a 30-minute game; gaps between ability fires feel empty.
