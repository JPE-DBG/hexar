package game

import "math"

func BuildCost(building BuildingType) float64 {
	switch building {
	case BuildingEconomy:
		return EconomyBuildCost
	case BuildingDefense:
		return DefenseBuildCost
	default:
		return 0
	}
}

func UpgradeCost(building BuildingType, currentLevel int) float64 {
	base := BuildCost(building)
	return base * math.Pow(2, float64(currentLevel))
}

func TotalInvested(building BuildingType, level int) float64 {
	base := BuildCost(building)
	return base * (math.Pow(2, float64(level)) - 1)
}

func ValidateBuild(state *GameState, action Action) error {
	hs, ok := state.Hexes[action.Target]
	if !ok {
		return ErrHexNotFound
	}
	if hs.Owner != action.Player {
		return ErrNotOwner
	}
	if hs.Building != BuildingNone {
		return ErrHasBuilding
	}
	cost := BuildCost(action.Building)
	if state.Players[action.Player].Gold < cost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyBuild(state *GameState, action Action) {
	state.Players[action.Player].Gold -= BuildCost(action.Building)
	hs := state.Hexes[action.Target]
	hs.Building = action.Building
	hs.Level = 1
}

func ValidateUpgrade(state *GameState, action Action) error {
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
	cost := UpgradeCost(hs.Building, hs.Level)
	if state.Players[action.Player].Gold < cost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyUpgrade(state *GameState, action Action) {
	hs := state.Hexes[action.Target]
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
