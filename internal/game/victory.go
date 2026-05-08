package game

func TriggerVictory(state *GameState, winner PlayerID) {
	if state.Over {
		return
	}
	state.Over = true
	state.Winner = winner
	for pid := range state.Players {
		if pid == winner {
			continue
		}
		for _, hs := range state.Hexes {
			if hs.Owner == pid {
				hs.Owner = NoPlayer
				hs.Building = BuildingNone
				hs.Level = 0
				hs.Capital = false
			}
		}
		remaining := state.Battles[:0]
		for _, b := range state.Battles {
			if b.Attacker != pid && b.Defender != pid {
				remaining = append(remaining, b)
			}
		}
		state.Battles = remaining
	}
}
