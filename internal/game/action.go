package game

import "errors"

type ActionType int

const (
	ActionClaim ActionType = iota
	ActionUpgrade
	ActionDemolish
	ActionAttack
)

type Action struct {
	Type     ActionType
	Player   PlayerID
	Target   Hex
	Building BuildingType
}

var (
	ErrHexNotFound       = errors.New("hex not found")
	ErrHexOwned          = errors.New("hex already owned")
	ErrNotAdjacent       = errors.New("no adjacent owned hex")
	ErrInsufficientGold  = errors.New("insufficient gold")
	ErrNotOwner          = errors.New("hex not owned by player")
	ErrNoBuilding        = errors.New("hex has no building")
	ErrNotEnemy          = errors.New("target is not enemy hex")
	ErrInsufficientPower = errors.New("insufficient power to attack")
	ErrBattleInProgress  = errors.New("battle already in progress on hex")
	ErrPlayerNotFound    = errors.New("player not found")
)

func ProcessActions(state *GameState, actions []Action) {
	for _, a := range actions {
		switch a.Type {
		case ActionClaim:
			if ValidateClaim(state, a) == nil {
				ApplyClaim(state, a)
			}
		case ActionUpgrade:
			if ValidateUpgrade(state, a) == nil {
				ApplyUpgrade(state, a)
			}
		case ActionDemolish:
			if ValidateDemolish(state, a) == nil {
				ApplyDemolish(state, a)
			}
		case ActionAttack:
			if ValidateAttack(state, a) == nil {
				ApplyAttack(state, a)
			}
		}
	}
}

func ValidateClaim(state *GameState, action Action) error {
	hs, ok := state.Hexes[action.Target]
	if !ok {
		return ErrHexNotFound
	}
	if hs.Owner != NoPlayer {
		return ErrHexOwned
	}
	if !hasAdjacentOwned(state, action.Player, action.Target) {
		return ErrNotAdjacent
	}
	player := state.Players[action.Player]
	if player == nil {
		return ErrPlayerNotFound
	}
	if player.Gold < ClaimCost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyClaim(state *GameState, action Action) {
	state.Players[action.Player].Gold -= ClaimCost
	state.Hexes[action.Target].Owner = action.Player
}

func hasAdjacentOwned(state *GameState, player PlayerID, target Hex) bool {
	for _, n := range target.Neighbors() {
		hs, ok := state.Hexes[n]
		if !ok {
			continue
		}
		if hs.Owner == player {
			return true
		}
	}
	return false
}
