package game

func RunTick(state *GameState, dt float64) {
	if state.Over {
		return
	}
	state.Elapsed += dt
}
