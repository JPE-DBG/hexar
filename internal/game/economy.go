package game

import "sort"

func HexIncome(hs *HexState) float64 {
	if hs.Building == BuildingGold {
		return (BaseIncomePerSec + GoldPerLevel*float64(hs.Level)) * GoldBonusMultiplier
	}
	return BaseIncomePerSec
}

func hexIncomeForPlayer(hs *HexState, player *Player) float64 {
	if hs.Building != BuildingGold {
		return BaseIncomePerSec
	}
	income := (BaseIncomePerSec + GoldPerLevel*float64(hs.Level)) * GoldBonusMultiplier
	if player.Tech[TechCompoundGrowth] {
		income *= CompoundGrowthMultiplier
	}
	if player.Tech[TechProsperity] {
		income += ProsperityBonus
	}
	return income
}

func CalcMaintenance(hexCount int) float64 {
	tier1 := min(hexCount, MaintenanceTier1Cap)
	tier2 := min(max(hexCount-MaintenanceTier1Cap, 0), MaintenanceTier2Cap-MaintenanceTier1Cap)
	tier3 := max(hexCount-MaintenanceTier2Cap, 0)
	return float64(tier1)*MaintenanceTier1 + float64(tier2)*MaintenanceTier2 + float64(tier3)*MaintenanceTier3
}

func calcMaintenanceForPlayer(hexCount int, player *Player) float64 {
	if !player.Tech[TechSupplyLines] {
		return CalcMaintenance(hexCount)
	}
	tier1 := min(hexCount, MaintenanceTier1Cap)
	tier2 := min(max(hexCount-MaintenanceTier1Cap, 0), MaintenanceTier2Cap-MaintenanceTier1Cap)
	tier3 := max(hexCount-MaintenanceTier2Cap, 0)
	return float64(tier1)*SupplyLinesTier1 + float64(tier2)*SupplyLinesTier2 + float64(tier3)*SupplyLinesTier3
}

func calcIncomeStats(state *GameState) (income map[PlayerID]float64, hexCount map[PlayerID]int) {
	income = make(map[PlayerID]float64)
	hexCount = make(map[PlayerID]int)
	for _, hs := range state.Hexes {
		if hs.Owner != NoPlayer {
			hexCount[hs.Owner]++
			income[hs.Owner] += hexIncomeForPlayer(hs, state.Players[hs.Owner])
		}
	}
	return
}

func RunEconomy(state *GameState, dt float64) {
	playerIncome, playerHexCount := calcIncomeStats(state)
	playerTPIncome := make(map[PlayerID]float64)

	for _, hs := range state.Hexes {
		if hs.Owner != NoPlayer && hs.Building == BuildingResearch {
			playerTPIncome[hs.Owner] += ResearchPerLevel * float64(hs.Level)
		}
	}

	for pid, player := range state.Players {
		income := playerIncome[pid]
		maintenance := calcMaintenanceForPlayer(playerHexCount[pid], player)
		player.Gold += (income - maintenance) * dt
		if player.Gold < 0 {
			player.Gold = 0
		}
		player.TP += playerTPIncome[pid] * dt
	}
}

func RunAutoDropPhase(state *GameState, dt float64) {
	playerIncome, playerHexCount := calcIncomeStats(state)

	for pid, player := range state.Players {
		net := playerIncome[pid] - calcMaintenanceForPlayer(playerHexCount[pid], player)
		gracePeriod := AutoDropGracePeriod
		if player.Tech[TechResilience] {
			gracePeriod = ResilienceGracePeriod
		}
		if net < 0 {
			if !player.AutoDropActive {
				player.AutoDropActive = true
				player.AutoDropGrace = gracePeriod
			} else {
				player.AutoDropGrace -= dt
				if player.AutoDropGrace < TickDt {
					autoDropLowestHex(state, pid)
					player.AutoDropActive = false
				}
			}
		} else {
			player.AutoDropActive = false
			player.AutoDropGrace = 0
		}
	}
}

func autoDropLowestHex(state *GameState, pid PlayerID) {
	type entry struct {
		h  Hex
		hs *HexState
	}
	var owned []entry
	for h, hs := range state.Hexes {
		if hs.Owner == pid && !hs.Capital && !hasBattleOnHex(state, h) {
			owned = append(owned, entry{h, hs})
		}
	}
	if len(owned) == 0 {
		return
	}

	sort.Slice(owned, func(i, j int) bool {
		a, b := owned[i], owned[j]
		ia, ib := HexIncome(a.hs), HexIncome(b.hs)
		if ia != ib {
			return ia < ib
		}
		ta, tb := TotalInvested(a.hs.Building, a.hs.Level), TotalInvested(b.hs.Building, b.hs.Level)
		if ta != tb {
			return ta < tb
		}
		if a.h.Q != b.h.Q {
			return a.h.Q < b.h.Q
		}
		return a.h.R < b.h.R
	})

	target := owned[0]
	hs := target.hs
	player := state.Players[pid]

	if hs.Building != BuildingNone {
		refund := AutoDropRefund
		if player.Tech[TechResilience] {
			refund = ResilienceDropRefund
		}
		player.Gold += TotalInvested(hs.Building, hs.Level) * refund
	}
	hs.Owner = NoPlayer
	hs.Building = BuildingNone
	hs.Level = 0
	hs.Capital = false
}
