package game

func CheckVictory(state *GameState) {
	for pid, p := range state.Players {
		if p.CapitalHP <= 0 {
			winner := opponentOf(pid)
			triggerVictory(state, winner, "capital")
			return
		}
	}
}

func triggerVictory(state *GameState, winner PlayerID, reason string) {
	if state.Over {
		return
	}
	state.Over = true
	state.Winner = winner
	state.WinReason = reason
}
