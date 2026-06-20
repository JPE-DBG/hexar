package game

func RunBar(state *GameState, dt float64) {
	for _, p := range state.Players {
		rate := BarFillRate
		for i := range p.BarBoosts {
			rate += p.BarBoosts[i].Bonus
		}
		p.Bar += rate * dt
		if p.Bar > BarMax {
			p.Bar = BarMax
		}

		// tick down active boosts
		active := p.BarBoosts[:0]
		for _, b := range p.BarBoosts {
			b.TimeLeft -= dt
			if b.TimeLeft > 0 {
				active = append(active, b)
			}
		}
		p.BarBoosts = active
	}
}
