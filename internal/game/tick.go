package game

func RunTick(state *GameState, dt float64, actions []ClaimAction) {
	if state.Over {
		return
	}
	ProcessActions(state, actions)
	RunEconomy(state, dt)
	state.Elapsed += dt
}
