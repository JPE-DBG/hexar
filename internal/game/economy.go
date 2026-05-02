package game

import "sort"

func HexIncome(hs *HexState) float64 {
	if hs.Building == BuildingGold {
		return (BaseIncomePerSec + GoldPerLevel*float64(hs.Level)) * GoldBonusMultiplier
	}
	return BaseIncomePerSec
}

func CalcMaintenance(hexCount int) float64 {
	tier1 := min(hexCount, MaintenanceTier1Cap)
	tier2 := min(max(hexCount-MaintenanceTier1Cap, 0), MaintenanceTier2Cap-MaintenanceTier1Cap)
	tier3 := max(hexCount-MaintenanceTier2Cap, 0)
	return float64(tier1)*MaintenanceTier1 + float64(tier2)*MaintenanceTier2 + float64(tier3)*MaintenanceTier3
}

func RunEconomy(state *GameState, dt float64) {
	playerIncome   := make(map[PlayerID]float64)
	playerHexCount := make(map[PlayerID]int)
	playerTPIncome := make(map[PlayerID]float64)

	for _, hs := range state.Hexes {
		if hs.Owner != NoPlayer {
			playerHexCount[hs.Owner]++
			playerIncome[hs.Owner] += HexIncome(hs)
			if hs.Building == BuildingResearch {
				playerTPIncome[hs.Owner] += ResearchPerLevel * float64(hs.Level)
			}
		}
	}

	for pid, player := range state.Players {
		income := playerIncome[pid]
		maintenance := CalcMaintenance(playerHexCount[pid])
		player.Gold += (income - maintenance) * dt
		if player.Gold < 0 {
			player.Gold = 0
		}
		player.TP += playerTPIncome[pid] * dt
	}
}

func RunAutoDropPhase(state *GameState, dt float64) {
	playerIncome   := make(map[PlayerID]float64)
	playerHexCount := make(map[PlayerID]int)

	for _, hs := range state.Hexes {
		if hs.Owner != NoPlayer {
			playerHexCount[hs.Owner]++
			playerIncome[hs.Owner] += HexIncome(hs)
		}
	}

	for pid, player := range state.Players {
		net := playerIncome[pid] - CalcMaintenance(playerHexCount[pid])
		if net < 0 {
			if !player.AutoDropActive {
				player.AutoDropActive = true
				player.AutoDropGrace = AutoDropGracePeriod
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
		if a.h.Q != b.h.Q {
			return a.h.Q < b.h.Q
		}
		return a.h.R < b.h.R
	})

	target := owned[0]
	hs := target.hs
	player := state.Players[pid]

	if hs.Building != BuildingNone {
		player.Gold += TotalInvested(hs.Building, hs.Level) * AutoDropRefund
	}
	hs.Owner = NoPlayer
	hs.Building = BuildingNone
	hs.Level = 0
	hs.Capital = false
	player.AutoDropActive = false
}
