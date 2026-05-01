package game

import "errors"

type ClaimAction struct {
	Player PlayerID
	Target Hex
}

var (
	ErrHexNotFound      = errors.New("hex not found")
	ErrHexOwned         = errors.New("hex already owned")
	ErrNotAdjacent      = errors.New("no adjacent owned hex with power")
	ErrInsufficientGold = errors.New("insufficient gold")
)

func ValidateClaim(state *GameState, action ClaimAction) error {
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
	if player.Gold < ClaimCost {
		return ErrInsufficientGold
	}

	return nil
}

func ApplyClaim(state *GameState, action ClaimAction) {
	state.Players[action.Player].Gold -= ClaimCost
	state.Hexes[action.Target].Owner = action.Player
}

func ProcessActions(state *GameState, actions []ClaimAction) {
	for _, a := range actions {
		if ValidateClaim(state, a) == nil {
			ApplyClaim(state, a)
		}
	}
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
