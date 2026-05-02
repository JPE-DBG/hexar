package game

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
	defenderPower := hs.Power()
	if attackerPower <= defenderPower {
		return ErrInsufficientPower
	}

	if state.Players[action.Player] == nil {
		return ErrPlayerNotFound
	}
	if state.Players[action.Player].Gold < AttackCost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyAttack(state *GameState, action Action) {
	state.Players[action.Player].Gold -= AttackCost
	hs := state.Hexes[action.Target]
	attackerPower := bestAdjacentPower(state, action.Player, action.Target)
	defenderPower := hs.Power()
	diff := attackerPower - defenderPower

	if diff > InstantTakeoverMinDiff {
		transferHex(state, action.Target, action.Player)
	} else {
		duration := 5.0 + float64(attackerPower+defenderPower)/2.0
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

	if attackerHex.Power() > defenderHex.Power() {
		transferHex(state, b.DefenderHex, b.Attacker)
	}
}

func transferHex(state *GameState, target Hex, newOwner PlayerID) {
	hs := state.Hexes[target]
	hs.Owner = newOwner
	hs.Building = BuildingNone
	hs.Level = 0
	hs.Capital = false
}

func bestAdjacentPower(state *GameState, player PlayerID, target Hex) int {
	best := 0
	for _, n := range target.Neighbors() {
		hs, ok := state.Hexes[n]
		if !ok {
			continue
		}
		if hs.Owner == player {
			p := hs.Power()
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
			p := hs.Power()
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
