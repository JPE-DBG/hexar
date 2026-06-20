# Combat Variants — Charge + Card Deck + Roguelike
## 10 Detailed Design Breakdowns

Generated: 2026-06-13. Continuation of `combat-systems-brainstorm.md` (theme 4 × theme 9).
See `combat-design-constraints.md` for hard rules all variants must satisfy.

---

## Variant 1: Charge Deck — Draw to Attack

### Core Mechanic
Economy hexes fill a shared **charge pool** passively (base 0.5/sec per hex; Economy buildings boost this). Max pool: 200 charge. Players hold a **hand of up to 4 cards**; one new card is auto-drawn from their personal deck every 8 seconds. Cards cost charge to play and target specific hexes.

Default starting cards (2 copies each in starting deck):
- **Strike** (20 charge): capture one adjacent enemy hex; 3-second contest (defender can counter-spend gold)
- **Barrier** (15 charge): protect one of your hexes from being captured for 12 seconds

Claiming unclaimed hexes costs 5 charge, no card required (allows early expansion).

Roguelike research adds one card of player's choice to their deck from a pool of ~20 card types.

### Early Game (0–8 min)
Both players race to claim unclaimed hexes (5 charge each, instant). Charge generation is low — only 1–2 Strikes possible in the first few minutes. Border forms when both players meet. First Strike cards trade hexes back and forth. Barrier on a single hex is weak (opponent ignores it and attacks adjacent) — players learn this quickly. First research pick around minute 3 adds a stronger or more specialized card to the deck.

### Mid Game (8–20 min)
Decks enriched with 2–4 researched cards. Charge generation higher from larger territory and Economy buildings. Build identity emerges: aggressive decks (Surge + Fast Strike + Chain) versus disruptive decks (Drain + Raze + Strike). Opponent's charge bar is visible — high charge means incoming Strike, prompting either a preemptive attack or a Barrier play. Card draw timing becomes tactical: hold a Strike in hand to play immediately after opponent wastes their Barrier.

### Late Game (20–30 min)
Decks have 8–12 total cards; richer hands. High charge enables multi-Strike pushes within seconds. Capital approaches: Strikes chain toward the capital while opponent tries to respond. Drain cards create windows by emptying opponent's charge mid-combo.

### Win Condition
Enemy capital captured via Strike. Capital counts as an adjacent enemy hex once your territory reaches it.

### Player Interactions
High frequency. Every Strike is a direct attack; every Barrier is a direct response or anticipation. Drain disrupts opponent's timing. Players react to each other's visible charge levels constantly.

### Roguelike Integration
Each research pick adds 1 card from a themed pool:
- **Aggressor pool**: Chain Strike (Strike enables free follow-up if target was unclaimed), Blitz Strike (3x charge cost but instant, no contest window), Surge (5-second triple charge gen)
- **Disruptor pool**: Drain (steal 40 charge from opponent), Raze (destroy one enemy building without capturing), Sabotage (opponent's next card costs double)
- **Economic pool**: Efficiency (reduce all card costs by 20% for 30 seconds), Harvest (gain charge from captured hexes), Reinforce (your hexes generate double charge for 15 seconds)

Strong build: Surge + Chain Strike + Drain (charge up, steal theirs, chain attacks).
Weak build: too many Barriers — can never push.

### Known Weaknesses
- **Auto-draw creates randomness**: drawing 3 Barriers when you need to attack feels bad and unfair.
- **Hand limit of 4 can block**: if you can't play cards fast enough, new draws stop; deck feels clogged.
- **Barrier on one hex is nearly useless** — must add area-defense cards (Shield Zone: protect a 3-hex cluster) or Barrier becomes a dead card.
- 8-second draw interval may need tuning to 5 seconds for real-time feel.

---

## Variant 2: Signature Ability + Charge Gate

### Core Mechanic
Each player has exactly **one signature ability**, chosen from 4 archetypes at game start. A **global charge bar** fills based on territory size and Economy buildings. First fill: ~8 minutes at base economy; faster with investment. At 100% the player can trigger their signature at a chosen moment (it does NOT auto-fire). After firing, bar resets; recharge takes ~5 minutes.

Archetypes:
- **Conqueror**: instantly capture all enemy hexes within radius 2 of a target point
- **Bastion**: all your hexes immune to capture for 20 seconds
- **Scavenger**: claim all unclaimed hexes on the entire map + 30-second double income spike
- **Saboteur**: destroy all buildings in a radius-2 zone of enemy territory (no capture)

Roguelike research upgrades your signature (wider radius, faster recharge, added secondary effect) — you never change archetype mid-game.

### Early Game (0–8 min)
Pure economy phase. No combat at all. Both players build territory (unclaimed hexes free to claim) and Economy buildings. Completely interaction-free. Players watch each other's charge bar fill. Only decision: build wide (more hexes = faster charge from territory count) vs build deep (Economy buildings on fewer hexes = more charge per hex). This early game is a real weakness of this variant.

### Mid Game (8–20 min)
First signatures fire around minute 8–10. Dramatic territorial swings: Conqueror grabs 7–10 hexes instantly. Counter-strategies unlock: Bastion blocks Conqueror (all hexes immune during their charge window). Strategic timing becomes the core skill: fire defensively to save territory, or offensively to grab the capital region? After firing, 5-minute recharge begins — the firer is vulnerable for this window.

### Late Game (20–30 min)
2–3 fires per player possible. Research upgrades have stacked; signatures are stronger. Capital becomes the primary Conqueror target. The decisive question: does the Conqueror player fire when the Bastion player's immunity is on cooldown? Whoever times the capital attempt correctly wins.

### Win Condition
Capture enemy capital via Conqueror (the only archetype that directly takes hexes). Saboteur + Scavenger can win indirectly by collapsing the opponent's economy until they can't defend.

### Player Interactions
**Rare but massive**. Each signature fire is a 5-second spectacle. Between fires: watching charge bars is the only interaction — slow and tense but not action-driven. Low interaction frequency.

### Roguelike Integration
Very linear: each research pick upgrades one signature trait. Conqueror progression: radius 2 → radius 3 → "Grand Conquest" (includes capital attack modifier). Research builds are predetermined by archetype choice — no divergence within a run. This is a weakness.

### Known Weaknesses
- **0–8 min is entirely interaction-free** — serious problem for a real-time game.
- **Only 2–3 decision moments per game** (when to fire) — insufficient for 30 minutes of play.
- **Bastion on individual hexes**: with a global immunity, Bastion is the ONE variant where area-defense works. Other variants can't borrow this.
- **Conqueror dominates**: other archetypes can't capture the capital directly, making them feel like inferior choices.
- Roguelike feels cosmetic rather than transformative.

---

## Variant 3: Three Slots, Slow Cooldowns

### Core Mechanic
Each player has **3 ability slots**, each on independent cooldown timers. No charge resource — abilities recharge on real-time clocks. All combat comes from these 3 abilities. Unclaimed hexes are free to claim (no ability needed), allowing early expansion.

Default starting abilities:
- **Advance** (20s CD): capture 1 adjacent enemy hex; 3-second contest window
- **Shield** (35s CD): one of your hexes immune to capture for 25 seconds  
- **Raze** (30s CD): destroy the building on 1 adjacent enemy hex (no capture)

Roguelike research: each pick replaces one slot with a stronger ability OR upgrades an existing one (shorter CD, wider effect). Player sees what their current 3 abilities are before picking.

### Early Game (0–8 min)
Rapid expansion into unclaimed hexes. First Advance available at second 20. Border forms and first Advances trade hexes. Players discover Shield's weakness: it protects one hex but opponent Advances into the adjacent one. Raze is strong early: destroy an opponent's Economy building, then Advance into the weakened hex next turn. Research at ~3 minutes gives first slot upgrade choice.

### Mid Game (8–20 min)
1–2 research picks have transformed one or two slots. Diverging kits: aggressive player swaps Shield for **Chain Advance** (captures 2 hexes in a line, 30s CD); defensive player upgrades Shield to **Fortress** (cluster of 3 hexes immune, 50s CD). Players learn each other's cooldown timings — "their Advance just fired, I have 20 seconds to push." Raze + Advance combo becomes a core tactic: Raze depletes the building refund, then Advance takes the weakened hex before opponent rebuilds.

### Late Game (20–30 min)
Fully custom 3-ability kits. Long CDs make every use consequential: waste Advance = 20-second vulnerability. Capital rush requires timing: chain Advances toward capital while opponent's Shield/Fortress is on cooldown. If Shield is unavailable when Advance hits capital — game over.

### Win Condition
Capture enemy capital via Advance (or Chain Advance). Requires positioning your territory adjacent to capital first.

### Player Interactions
**Moderate frequency, high stakes**. One ability fires every 20–40 seconds. Each use is a meaningful commitment. Opponents can read cooldown timers (transparent strategic information). Raze + Advance creates clear 2-step interactions with visible decision windows.

### Roguelike Integration
Strong ability divergence possible. Example builds:
- **Full Aggressor**: Chain Advance + Blitz Advance (instant, no contest, 45s CD) + Suppression (Raze entire border hex cluster)
- **Disruptor**: Raze + Sabotage (opponent ability on random CD penalty) + Chain Raze
- **Fortifier**: Fortress + Counter (when opponent Advances into your Fortress, auto-counter-Advance) + Resilience (immune hex reverts to you if captured while Shield active)

### Known Weaknesses
- **Shield on one hex is almost always wasted** — must upgrade to Fortress (area) to matter. Shield as a starting ability misleads players into thinking single-hex defense is viable.
- **3 abilities for 30 minutes is too few**: late game feels repetitive — same 3 abilities in rotation with no new stimulation.
- **Dead time between cooldowns**: with 20–40s CDs and only 3 abilities, there are ~10-second gaps where nothing is happening.
- **Fortify variant is meaningless at large territory scale**: same anti-pattern flagged across variants.

---

## Variant 4: Combo Cards

### Core Mechanic
Cards are drawn from a deck (same draw mechanic as Variant 1, one every 8 seconds, hand of 4). Playing **two cards in sequence within 3 seconds** creates a combo with an enhanced effect. Single cards still work (weaker version). Charge is the resource for card plays. Combo recognition is shown in the UI (cards glow when a valid combo partner is in hand).

Example combos:
- **Strike → Rally** (within 3s): capture hex AND it generates double income for 60 seconds
- **Barrier → Strike** (within 3s): protected rush — the capture cannot be countered during its contest window
- **Drain → Strike** (within 3s): drain opponent's charge, then immediately take hex while they can't respond
- **Surge → Drain** (within 3s): generate charge then immediately convert it to a massive Drain

Roguelike research adds cards that unlock new combo recipes, or adds "Catalyst" cards that combo with everything.

### Early Game (0–8 min)
Few cards, no combos yet (need 2 compatible cards simultaneously in hand). Expansion into unclaimed hexes continues freely. Single Strikes trade border hexes. First research pick adds the first "combo ingredient" and reveals the combo recipe to the player.

### Mid Game (8–20 min)
Combos start firing as decks deepen. Core skill: hold a Strike in hand until the matching combo card arrives, then execute within 3 seconds. Opponent watches hand size and tries to predict what combo is being assembled. Drain + Strike is the first reliable combo — disrupts opponent's timing, then capitalizes immediately. Research adds higher-tier combos and multi-card chains.

### Late Game (20–30 min)
3-card chains possible: Surge → Drain → Strike (build charge, steal theirs, attack in one fluid sequence). Capital assault requires assembling and holding a complete combo while under pressure. Opponent must either disrupt the assembly (Drain them first) or accept the hit and respond.

### Win Condition
Capital captured via Strike (or combo that includes a Strike variant). The decisive combo is usually a protected rush into the capital.

### Player Interactions
High with psychological depth. Watching hand sizes: "they've held a Strike for 10 seconds — what are they waiting for?" Drain combos directly punish the opponent for building charge. Counter-play: disrupt their charge before they execute the combo.

### Roguelike Integration
Each research pick adds a specific card (and reveals the combo recipe if it creates a new one). Strong builds: deck with 2–3 complete reliable combos in regular rotation. Weak builds: mismatched cards with no combo synergy — all individual card plays.

### Known Weaknesses
- **3-second combo window is very tight in real-time** — may need to extend to 5 seconds, or have UI assistance that highlights valid combos.
- **High new-player barrier**: learning which cards combo requires either a guide or experimentation that costs them the game.
- **Randomness of draw matters more here** than in V1: if the second half of your combo never draws, your whole strategy collapses.
- **Barrier is useless in combos unless Barrier → Strike is one of the first combos unlocked** — otherwise it's a wasted deck slot.

---

## Variant 5: Charge Lanes

### Core Mechanic
The map has **2 defined corridors** between the two capitals (guaranteed by map generation). Each corridor has its own **lane charge bar**, filled by hexes the player owns within that lane. When a player's lane charge exceeds the opponent's lane charge in the same corridor, the border in that lane **auto-advances** toward the opponent at a speed proportional to the charge differential. No manual attack button — the lane pushes automatically.

Abilities (on cooldowns, not charge-gated):
- **Push** (spend all accumulated lane charge for an instant burst advance in that lane)
- **Hold** (burn your lane charge to temporarily neutralize the opponent's advance in that lane)
- **Redirect** (shift charge from one lane to the other — costs a 10% charge penalty)

Roguelike research unlocks lane specializations and cross-lane effects.

### Early Game (0–8 min)
Players expand into both lanes. Whoever controls more hexes in a lane has the charge advantage there. Immediate low-level interaction: both players pushing into lane hexes. Decision: split evenly between lanes, or commit heavily to one?

### Mid Game (8–20 min)
Lane advantage/disadvantage becomes clear. Redirect ability is the key skill: transfer charge from the lane you're winning to the lane you're losing. Research picks lock in a lane strategy — specialize for burst (faster auto-advance in one lane) or balance (Redirect is cheaper). If opponent dominates both lanes, the only path back is to break their lane economy (take their core income hexes in a lane via Hold → Push combo).

### Late Game (20–30 min)
One lane's auto-advance reaches the capital region. Final Push burns all remaining lane charge for a decisive burst. Opponent must sacrifice the other lane to Hold this one — creating an opening elsewhere.

### Win Condition
One lane's auto-advance reaches and enters the enemy capital hex.

### Player Interactions
**Constant low-level** (both lanes ticking) with **periodic high-level** events (Push/Hold fires). The auto-advance creates continuous visible pressure — not chaotic because it's predictable and lane-bound.

### Roguelike Integration
Lane specializations: **Northern Aggressor** (lane 1 Push costs halved, lane 1 auto-advance 40% faster), **Balanced Commander** (Redirect costs nothing), **Counterpusher** (Hold generates charge instead of consuming it). Each archetype completely changes how the lanes are managed.

### Known Weaknesses
- **Requires specific map design** — map must have exactly 2 clear corridors between capitals. Constrains map generation significantly.
- **Does not scale to 3–4 players** — "lane" concept becomes unclear with multiple capitals.
- **If opponent concedes one lane**, the winning player advances in it unopposed — the opponent's concession should be punished but the auto-advance might be too slow to matter.
- **No individual hex targeting** — feels less strategic, more like watching bars fill.
- No fortification meaningful at individual hex level.

---

## Variant 6: Card Draft Waves

### Core Mechanic
Every **90 seconds**, a "draft wave" occurs: both players are shown the **same 3 random cards** from a shared pool. Each player secretly picks 1; picks are revealed simultaneously. The unchosen cards are discarded (lost to both). Picked cards go directly to hand and are immediately playable using charge. Cards **expire** after 3 minutes (forcing active play). No base combat ability exists — all combat comes from drafted cards. Claiming unclaimed hexes remains free.

Draft pool contains ~30 card types across combat, economy, and disruption themes. Late-wave pools include more powerful cards.

### Early Game (0–8 min)
Waves at 90s, 3:00, 4:30, 6:00. First 90 seconds: pure expansion, zero interaction. **Draft 1 is critical** — your first combat card. The opponent's pick is revealed — you know what they took. Hate-drafting is legal (take a card you don't want to deny opponent). This creates immediate strategic tension even before any combat.

### Mid Game (8–20 min)
4–6 draft waves have fired. Hands are richer. Card expiry forces active use — sitting on cards is penalized. Meta-game emerges: "they've drafted 3 Strike variants — they're going aggressive, I should pick defensive or disruption cards in future waves." Contested draft moments: both players want the same card; who adapts? Economy investment in hexes increases charge generation, enabling more frequent card plays from the growing hand.

### Late Game (20–30 min)
Waves 10–12 include late-game powerful cards (area captures, global effects). Full hands with 6–8 cards available. Decisive combo plays from accumulated hand. Draft selection from late pool is a major strategic moment — a single late-game card can shift the endgame dramatically.

### Win Condition
Capture enemy capital via drafted Strike/Capture cards. The player who built the most capital-targeting card combination wins the endgame draft race.

### Player Interactions
**Two layers**: during battles (card plays are direct interactions); between battles (watching opponent's hand size and draft history to infer strategy). The reveal of each player's draft pick is its own interactive moment — instant strategic information.

### Roguelike Integration
The draft **is** the roguelike — no separate research system needed. Card pool varies each game (not all 30 cards appear). Strong runs: discovering a powerful synergy and drafting both halves before the opponent realizes it. Replayability comes entirely from the random pool and draft decisions.

### Known Weaknesses
- **Hate-drafting can feel frustrating** — opponent takes your best card "just because." May need to add a rule: "both players see different 3-card options" (removes hate-drafting but loses the contested-draft tension).
- **First 90 seconds has zero player interaction** — early game void.
- **Card expiry adds UI complexity** — each card needs a visible 3-minute countdown.
- **Late-game powerful cards can auto-win** — pool must be carefully tuned so late cards don't make early game irrelevant.
- Wave timing (90s) needs careful testing — may feel too slow early and too fast late.

---

## Variant 7: Charge + Reaction Window

### Core Mechanic
Charge fills based on territory and Economy buildings. Attacker spends **30 charge** to issue a **Challenge** to any adjacent enemy hex. A 5-second countdown appears on that hex — visible to both players. The defender has 5 seconds to spend **30 charge** to Block. If blocked: attacker's charge is lost, hex stays. If not blocked in 5 seconds: hex flips to attacker. **Only 1 active challenge per player at a time.** Abilities (from roguelike) modify the challenge mechanic rather than being separate attack types.

Core ability cards:
- **Blitz** (modify next Challenge): reaction window reduced from 5s to 2s
- **Feint** (new Challenge type): fake challenge costs attacker nothing — purely to force opponent to waste Block charge
- **Reflect** (defender ability): successful Block converts to a free counter-Challenge on the attacker's hex
- **Reinforce** (defender ability): Block costs 15 charge instead of 30 for the next 60 seconds

### Early Game (0–8 min)
Expansion via unclaimed hexes (free, no challenge). Border forms; first challenges begin. Players learn the mechanic: see the countdown, decide whether to Block or accept the loss. Early challenges test each other's reaction habits. Low charge = limited challenge frequency. First research pick around minute 3 adds a challenge modifier.

### Mid Game (8–20 min)
Charge generation increases. Challenges fire every 30–45 seconds. Feint becomes a key disruptor: force the opponent to Block a fake, then immediately fire a real Challenge while they're charge-depleted. Blitz demands split-second reactions — 2 seconds is genuinely stressful. Counter-play: hold Block charge in reserve specifically to absorb a Blitz. Research shapes whether you're a Challenger (offensive modifiers) or a Blocker (defensive modifiers).

### Late Game (20–30 min)
Fast charge generation enables challenge chains. Decisive capital challenge: attacker fires Blitz Challenge on capital — defender has 2 seconds to burn 30 charge or lose the game. Feint before the real capital challenge empties the defender's reserve.

### Win Condition
Successful unblocked Challenge on enemy capital hex.

### Player Interactions
**High and direct.** Every Challenge is a real-time call-and-response. The 5-second countdown creates visible, readable tension. Feint adds a psychological bluffing layer. Most interactive variant in the list.

### Roguelike Integration
Each research pick adds one Challenge modifier to hand. Strong build: Feint + Blitz + Chain Challenge (successful challenge enables free immediate follow-up). Weak build: too many defensive Reinforce cards — can hold territory but never push.

### Known Weaknesses
- **5-second window punishes latency** — players on slow connections or mobile will consistently fail to Block even if they intend to. Needs server-side lag compensation.
- **Blitz (2s window) may be dominant and unfun** — if it becomes a must-pick, every interaction becomes unpleasant.
- **Block-or-don't is binary** — limited decision depth beyond timing. May feel shallow after many games.
- **If attacker always has charge advantage** (from better economy), they can challenge faster than defender can Block — reduces to economic dominance again.

---

## Variant 8: Resource Conversion Pipeline

### Core Mechanic
Three resources: **Gold** (generated by Economy hexes/buildings), **Charge** (converted from gold: costs 1 gold/sec to generate 2 charge/sec; conversion is an active allocation choice, not automatic), **Cards** (permanent abilities from roguelike research only). Combat requires playing a Card, which costs Charge. The pipeline creates deliberate planning: you must decide in advance to start converting gold to charge, then play the card when charge is ready.

No base attack exists. Cards are the only way to affect enemy hexes. First card arrives via first research pick.

### Early Game (0–8 min)
Build gold pipeline (Economy buildings). First research pick (~3 min) grants first Card and a combat option. Players learn the conversion trade-off: every gold converted to charge is gold not spent on buildings. Zero combat before first research — parallel solitaire problem. Once first card is acquired, first combat possible.

### Mid Game (8–20 min)
Pipeline established. Both players manage gold allocation: build vs convert. Targeting enemy Economy hexes is now strategically valuable — disrupting their gold supply disrupts their combat pipeline, not just their income. Research unlocks more efficient cards or pipeline improvements (faster conversion, higher charge cap). Players start see-sawing: high conversion phase (charge up for attacks) vs high building phase (recover economy after attack).

### Late Game (20–30 min)
Rich pipeline supports frequent card plays. Both players have full card kits from research. Disrupting the pipeline (Raze + Advance on Economy hexes) remains viable as a late-game strategy.

### Win Condition
Capital captured via card play. Requires territory adjacent to capital.

### Player Interactions
**Indirect and direct.** Indirect: targeting Economy hexes disrupts the pipeline. Direct: card plays are immediate hex-level attacks. Strategic layer: observe opponent's conversion rate (visible via charge bar filling speed) to predict incoming card play.

### Roguelike Integration
Research adds cards AND pipeline upgrades (batch conversion, higher charge cap). Strong build: efficient conversion + aggressive cards. Weak build: high economy but wrong cards for stage of game.

### Known Weaknesses
- **Three-resource system is complex to display** — UI must clearly show Gold, Charge, Cards, and conversion allocation simultaneously.
- **Early game has zero combat** until first research pick — same parallel solitaire problem as V2/V9.
- **Conversion wait feels awkward** — "I want to attack but I have to wait for charge to build" creates frustrating idle moments.
- **Most complex system in the list** — may be too much cognitive overhead for a fast real-time game.

---

## Variant 9: Permanent Card Upgrades Only (Zero Start)

### Core Mechanic
Players start with **zero combat abilities**. Economy generates a charge pool passively. All combat tools come exclusively from roguelike research picks — every pick adds a permanent card to your kit. Cards are on individual cooldowns (20–40s each). Claiming unclaimed hexes is free and instant (keeping early game active). Capital is capturable only via a specific card type unlocked later in the research tree.

Research tier structure:
- **Tier 1** (first 3 picks): basic abilities — Strike, Raze, Surge, Rush, Fortify Zone
- **Tier 2** (next 3 picks): intermediate — Chain Strike, Blitz, Drain, Sabotage
- **Tier 3** (final picks): powerful — Capital Strike, Grand Advance, Nova

### Early Game (0–8 min)
Pure expansion. Zero player-vs-player interaction for the first ~3 minutes. Both players claim unclaimed hexes; map fills up. **First research pick is the most critical decision in the game**: Strike (go offensive now), Raze (be a disruptor), or Rush (expand faster and delay fighting). This first pick signals your archetype and the opponent will adapt their next picks accordingly.

### Mid Game (8–20 min)
2–4 picks have built a small kit. Combat begins once both players have Strike-equivalent cards. Asymmetric fights: one player went all-aggressor, opponent went Raze-heavy (economically disruptive). Aggressor pushes borders while disruptor undermines their Economy. Research picks become reactive: see opponent's cards in use, counter-pick accordingly.

### Late Game (20–30 min)
Kits have 6–8 abilities. Tier 3 capital-capable cards become available. The player who reached Tier 3 research first (by investing in Research buildings) gets first access to capital-targeting cards — a decisive advantage.

### Win Condition
Enemy capital captured via Tier 3 **Capital Strike** card (or equivalent). Requires having unlocked it through research progression.

### Player Interactions
**Zero in early game; high in mid/late.** Once kits develop: complex ability rotations, counter-plays based on opponent's visible abilities. The opponent's research choices are visible from what abilities they use — strategic reading is a key skill.

### Roguelike Integration
**The most pure roguelike feel** of all variants. Your entire combat kit is built from scratch each game. No two games have the same build. Archetype paths emerge: Aggressor, Disruptor, Defender, Economist. Research building investment determines how fast you access higher tiers — a genuine build tradeoff.

### Known Weaknesses
- **0–3 minute dead zone** is serious — the first real interaction requires the first research pick, which is ~3 minutes in.
- **Fix**: give all players a free, weak default Strike from turn 1 (reduces zero-start purity but fixes the dead zone).
- **Late Tier 3 gating**: if only Tier 3 has capital-capture, and research is slow, games may go very long before anyone can win.
- **Very complex late game** (6–8 ability rotations) — may be too much for casual players.
- **New player experience is brutal**: picking the wrong first card can cripple the entire game.

---

## Variant 10: Charge Shapes (Area Releases)

### Core Mechanic
Charge fills from Economy. Instead of single-hex attacks, releasing charge applies a **shape** — a multi-hex pattern on the grid that is captured or affected simultaneously. All hexes in the shape resolve in one instant — no 5-second contest window. The shape must be positioned with the player's territory as the "anchor" point (at least one hex in the shape must be adjacent to your current territory).

Default shapes (always available):
- **Spike** (15 charge): 1 hex in a direction — captures that 1 hex
- **Nova** (50 charge): ring around one of your hexes — reclaims all adjacent unclaimed hexes; cannot capture enemy hexes

Roguelike research unlocks:
- **Wedge** (30 charge): V-shape of 3 hexes pointing at a target
- **Spear** (45 charge): 5-hex line in one direction
- **Flank** (60 charge): two simultaneous Spikes in different directions
- **Wall** (35 charge): fortify 4 of your hexes simultaneously (25-second immunity)
- **Harvest** (20 charge): claim all unclaimed hexes in a radius-1 zone (economy expansion)

Roguelike also provides shape modifiers: cost reduction, larger radius, added effects (bonus income on captured hexes).

### Early Game (0–8 min)
Nova shape is the dominant early tool: cheaply claims clusters of unclaimed hexes around your capital. Spike used for first border attacks. Low charge means only cheap shapes available. Research unlocks first combat-capable shape (Wedge or Spear), enabling real offensive threat.

### Mid Game (8–20 min)
More shapes unlocked; strategic layering. Wedge creates multi-hex attacks that force multi-hex responses — defense can't focus on a single point. **Wall shape is actually useful here** (protects 4 hexes simultaneously — addresses the single-hex defense anti-pattern). Flank creates two simultaneous pressure points, splitting the opponent's attention. Charge generation dictates how frequently shapes can fire.

### Late Game (20–30 min)
High charge enables shape chains. Spear + Wedge combo pushes a narrow corridor deep into enemy territory toward the capital. Harvest used to quickly reclaim lost hexes. Nova used for dramatic economic recovery if territory has been reduced.

### Win Condition
Any shape that includes the enemy capital hex captures it. Spear pointed directly at capital is the archetypal endgame threat.

### Player Interactions
**Periodic high-impact events.** Each shape fires simultaneously on 1–5 hexes — dramatic and readable. Charge bar is visible — both players see when a shape is imminent. Wall gives meaningful area defense (finally a defensive ability that works at scale). Timing: use Wall when Spear is incoming, then counter-Spear while opponent's shape is on charge cooldown.

### Roguelike Integration
Research unlocks new shapes and modifiers — radically changes the threat profile each game. Strong builds: Spear-focused (repeated deep strikes), Flank-focused (split pressure), Nova-focused (economic dominance through rapid unclaimed capture). Each shape has a natural "counter" shape, enabling meta-play.

### Known Weaknesses
- **UI challenge**: showing shape previews on the hex grid clearly is non-trivial — player needs to see the shape before committing charge to it.
- **No individual hex targeting**: if you want to take exactly hex X but the shape forces you to also take hexes Y and Z (enemy territory you're not ready to hold), shape releases can feel clumsy.
- **Charge-gating creates idle time**: like V1/V8, waiting for charge to fill for an expensive Spear creates gaps.
- **Nova is strictly economic** (can't capture enemy hexes) — need careful balance so it doesn't dominate early game while being useless late.

---

## Cross-Variant Notes

**Defensive ability problem** (flagged across all 10 variants):
Single-hex shield/fortify abilities are nearly always useless in a territory game because the opponent attacks adjacent hexes. Defensive abilities that work:
- Area shields (Fortress: 3-hex cluster; Wall shape: 4 hexes simultaneously)
- Temporal disruption (Block charge in V7, Reflect counter-attack)
- Economic disruption defense (Drain that penalizes attacker)
- Lane-level control (Hold in V5)
Never design a defensive ability that affects exactly 1 hex unless it's attached to a mechanic where that hex is guaranteed to be the only valid attack target.

**Early game void problem** (V2, V8, V9):
Any variant where no combat is possible for the first 3+ minutes fails the real-time feel. Fix: always provide a default weak attack option from turn 1, even if it's just "Spike" in V10 or a "Basic Strike" in V9.

**Most promising for Hexar's constraints**: V7 (Challenge + Reaction) and V9 (Zero Start + Roguelike) have the strongest roguelike integration. V10 (Charge Shapes) solves the defensive ability problem most elegantly. V6 (Card Draft Waves) has the most natural roguelike feel built in.
