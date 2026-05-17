package game

func RunTick(state *GameState, dt float64, actions []Action) {
	if state.Over {
		return
	}
	ProcessActions(state, actions)
	if state.Over {
		return
	}
	RunEconomy(state, dt)
	runUpgradeTimers(state, dt)
	RunBattles(state, dt)
	RunAutoDropPhase(state, dt)
	runFortifyTimers(state, dt)
	runVanguardTimers(state, dt)
	state.Elapsed += dt
}

func runFortifyTimers(state *GameState, dt float64) {
	for _, hs := range state.Hexes {
		if hs.FortifyTimer > 0 {
			hs.FortifyTimer -= dt
			if hs.FortifyTimer < 0 {
				hs.FortifyTimer = 0
			}
		}
	}
}

func runUpgradeTimers(state *GameState, dt float64) {
	for _, hs := range state.Hexes {
		if hs.UpgradeTimer > 0 {
			hs.UpgradeTimer -= dt
			if hs.UpgradeTimer < 0 {
				hs.UpgradeTimer = 0
			}
			// When timer expires, increment the level
			if hs.UpgradeTimer == 0 {
				hs.Level++
			}
		}
	}
}

func runVanguardTimers(state *GameState, dt float64) {
	for _, player := range state.Players {
		if player.VanguardTimer > 0 {
			player.VanguardTimer -= dt
			if player.VanguardTimer < 0 {
				player.VanguardTimer = 0
			}
		}
	}
}
