package game

func RunTick(state *GameState, dt float64, actions []Action) {
	if state.Over {
		return
	}
	ProcessActions(state, actions)
	if state.Over {
		return
	}
	RunBar(state, dt)
	RunUnits(state, dt)
	RunDeck(state, dt)
	CheckVictory(state)
	state.Elapsed += dt
}
