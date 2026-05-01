package game

func CalcIncome(hexCount int) float64 {
	return float64(hexCount) * BaseIncomePerSec
}

func CalcMaintenance(hexCount int) float64 {
	tier1 := min(hexCount, MaintenanceTier1Cap)
	tier2 := min(max(hexCount-MaintenanceTier1Cap, 0), MaintenanceTier2Cap-MaintenanceTier1Cap)
	tier3 := max(hexCount-MaintenanceTier2Cap, 0)
	return float64(tier1)*MaintenanceTier1 + float64(tier2)*MaintenanceTier2 + float64(tier3)*MaintenanceTier3
}

func RunEconomy(state *GameState, dt float64) {
	hexCounts := make(map[PlayerID]int)
	for _, hs := range state.Hexes {
		if hs.Owner != NoPlayer {
			hexCounts[hs.Owner]++
		}
	}

	for pid, player := range state.Players {
		count := hexCounts[pid]
		income := CalcIncome(count)
		maintenance := CalcMaintenance(count)
		player.Gold += (income - maintenance) * dt
		if player.Gold < 0 {
			player.Gold = 0
		}
	}
}
