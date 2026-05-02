package game

import "math"

func buildCost(building BuildingType) float64 {
	switch building {
	case BuildingGold:
		return GoldBuildCost
	case BuildingPower:
		return PowerBuildCost
	case BuildingResearch:
		return ResearchBuildCost
	default:
		return 0
	}
}

func UpgradeCost(building BuildingType, currentLevel int) float64 {
	return buildCost(building) * math.Pow(2, float64(currentLevel))
}

func TotalInvested(building BuildingType, level int) float64 {
	return buildCost(building) * (math.Pow(2, float64(level)) - 1)
}

func ValidateUpgrade(state *GameState, action Action) error {
	hs, ok := state.Hexes[action.Target]
	if !ok {
		return ErrHexNotFound
	}
	if hs.Owner != action.Player {
		return ErrNotOwner
	}
	if hs.Level == 0 {
		if action.Building == BuildingNone {
			return ErrNoBuilding
		}
	} else {
		if hs.Building == BuildingNone {
			return ErrNoBuilding
		}
	}
	building := hs.Building
	if hs.Level == 0 {
		building = action.Building
	}
	cost := UpgradeCost(building, hs.Level)
	if state.Players[action.Player] == nil {
		return ErrPlayerNotFound
	}
	if state.Players[action.Player].Gold < cost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyUpgrade(state *GameState, action Action) {
	hs := state.Hexes[action.Target]
	if hs.Level == 0 {
		hs.Building = action.Building
	}
	state.Players[action.Player].Gold -= UpgradeCost(hs.Building, hs.Level)
	hs.Level++
}

func ValidateDemolish(state *GameState, action Action) error {
	hs, ok := state.Hexes[action.Target]
	if !ok {
		return ErrHexNotFound
	}
	if hs.Owner != action.Player {
		return ErrNotOwner
	}
	if hs.Building == BuildingNone {
		return ErrNoBuilding
	}
	return nil
}

func ApplyDemolish(state *GameState, action Action) {
	hs := state.Hexes[action.Target]
	refund := TotalInvested(hs.Building, hs.Level) * DemolishRefund
	state.Players[action.Player].Gold += refund
	hs.Building = BuildingNone
	hs.Level = 0
}
