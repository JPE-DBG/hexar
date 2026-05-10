package game

func hexEffectivePower(state *GameState, playerID PlayerID, hex Hex) int {
	hs, ok := state.Hexes[hex]
	if !ok {
		return 0
	}
	power := hs.Power()
	if pl := state.Players[playerID]; pl != nil && pl.Tech[TechIronGrip] {
		power++
	}
	return power
}

func effectiveAttackCost(state *GameState, attacker PlayerID, target Hex) float64 {
	player := state.Players[attacker]
	if player == nil {
		return AttackCost
	}
	hs := state.Hexes[target]
	if hs == nil {
		return AttackCost
	}

	cost := AttackCost

	// Reclamation: -50g if previously owned by attacker
	if player.Tech[TechReclamation] && hs.PreviousOwner == attacker {
		cost -= ReclamationCost
	}

	// Vanguard: -50g if timer active (within 12s of last capture)
	if player.Tech[TechVanguard] && player.VanguardTimer > 0 {
		cost -= VanguardAttackCost
	}

	// Both techs stack: 100 - 50 - 50 = 0 (free attack when reclaiming during Vanguard)
	if cost < 0 {
		cost = 0
	}

	return cost
}

func garrisonBoostForDefender(state *GameState, defender PlayerID, defenderHex Hex) int {
	player := state.Players[defender]
	if player == nil || !player.Tech[TechGarrison] {
		return 0
	}
	boost := 0
	for _, n := range defenderHex.Neighbors() {
		hs, ok := state.Hexes[n]
		if !ok {
			continue
		}
		if hs.Owner == defender {
			boost++
		}
	}
	return min(boost, GarrisonMaxBoost)
}

func ValidateAttack(state *GameState, action Action) error {
	hs, ok := state.Hexes[action.Target]
	if !ok {
		return ErrHexNotFound
	}
	if hs.Owner == NoPlayer || hs.Owner == action.Player {
		return ErrNotEnemy
	}
	if hasBattleOnHex(state, action.Target) {
		return ErrBattleInProgress
	}

	attackerPower := bestAdjacentPower(state, action.Player, action.Target)
	defenderPower := hexEffectivePower(state, hs.Owner, action.Target)
	// Garrison is NOT validated here — it only activates during resolveBattle.
	// The UI prevents committing 100g to a guaranteed-loss battle by including
	// the Garrison bonus in its own threshold check (effectiveDefPower = defPower + garrisonBonus).
	if attackerPower <= defenderPower {
		return ErrInsufficientPower
	}

	if state.Players[action.Player] == nil {
		return ErrPlayerNotFound
	}
	cost := effectiveAttackCost(state, action.Player, action.Target)
	if state.Players[action.Player].Gold < cost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyAttack(state *GameState, action Action) {
	player := state.Players[action.Player]
	cost := effectiveAttackCost(state, action.Player, action.Target)
	player.Gold -= cost

	hs := state.Hexes[action.Target]
	attackerPower := bestAdjacentPower(state, action.Player, action.Target)
	defenderPower := hexEffectivePower(state, hs.Owner, action.Target)
	diff := attackerPower - defenderPower

	if diff > InstantTakeoverMinDiff && hs.FortifyTimer <= 0 {
		transferHex(state, action.Target, action.Player)
	} else {
		duration := BaseBattleDuration + float64(attackerPower+defenderPower)/2.0
		if player.Tech[TechSiegeMastery] {
			duration = max(SiegeMasteryMinDur, duration*SiegeMasteryMult)
		}
		state.Battles = append(state.Battles, Battle{
			AttackerHex: bestAdjacentHex(state, action.Player, action.Target),
			DefenderHex: action.Target,
			Attacker:    action.Player,
			Defender:    hs.Owner,
			TimeLeft:    duration,
		})
	}
}

func RunBattles(state *GameState, dt float64) {
	remaining := state.Battles[:0]
	for i := range state.Battles {
		b := &state.Battles[i]
		b.TimeLeft -= dt
		if b.TimeLeft <= 0 {
			resolveBattle(state, b)
		} else {
			remaining = append(remaining, *b)
		}
	}
	state.Battles = remaining
}

func resolveBattle(state *GameState, b *Battle) {
	attackerHex := state.Hexes[b.AttackerHex]
	defenderHex := state.Hexes[b.DefenderHex]
	if attackerHex == nil || defenderHex == nil {
		return
	}

	garrisonBoost := garrisonBoostForDefender(state, b.Defender, b.DefenderHex)
	totalBoost := min(CounterSpendCap, b.CounterBoost+garrisonBoost)

	attackerPow := hexEffectivePower(state, b.Attacker, b.AttackerHex)
	defenderPow := hexEffectivePower(state, b.Defender, b.DefenderHex)

	if attackerPow > defenderPow+totalBoost {
		transferHex(state, b.DefenderHex, b.Attacker)
	} else if attackerPow == defenderPow+totalBoost {
		attacker := state.Players[b.Attacker]
		if attacker != nil && attacker.Tech[TechSiegeMastery] {
			transferHex(state, b.DefenderHex, b.Attacker)
		}
	}
}

func transferHex(state *GameState, target Hex, newOwner PlayerID) {
	hs := state.Hexes[target]
	wasCapital := hs.Capital
	oldOwner := hs.Owner
	hs.PreviousOwner = oldOwner
	hs.Owner = newOwner
	hs.Building = BuildingNone
	hs.Level = 0
	hs.Capital = false
	hs.FortifyTimer = 0

	if pl := state.Players[newOwner]; pl != nil && pl.Tech[TechWarChest] && oldOwner != NoPlayer {
		pl.Gold += WarChestRefund
	}
	if pl := state.Players[newOwner]; pl != nil && pl.Tech[TechVanguard] {
		pl.VanguardTimer = VanguardWindow
	}
	if wasCapital {
		TriggerVictory(state, newOwner)
	}
}

func bestAdjacentPower(state *GameState, player PlayerID, target Hex) int {
	best := 0
	for _, n := range target.Neighbors() {
		hs, ok := state.Hexes[n]
		if !ok {
			continue
		}
		if hs.Owner == player {
			p := hexEffectivePower(state, player, n)
			if p > best {
				best = p
			}
		}
	}
	return best
}

func bestAdjacentHex(state *GameState, player PlayerID, target Hex) Hex {
	bestPower := 0
	bestHex := target
	for _, n := range target.Neighbors() {
		hs, ok := state.Hexes[n]
		if !ok {
			continue
		}
		if hs.Owner == player {
			p := hexEffectivePower(state, player, n)
			if p > bestPower {
				bestPower = p
				bestHex = n
			}
		}
	}
	return bestHex
}

func hasBattleOnHex(state *GameState, target Hex) bool {
	for _, b := range state.Battles {
		if b.DefenderHex == target {
			return true
		}
	}
	return false
}
