package game

func RunTick(state *GameState, dt float64, actions []Action) {
	if state.Over {
		return
	}
	ProcessActions(state, actions)
	RunEconomy(state, dt)
	RunBattles(state, dt)
	RunAutoDropPhase(state, dt)
	state.Elapsed += dt
}
