package game

import "errors"

type ActionType int

const (
	ActionPlayTopCard     ActionType = iota // play top card of deck
	ActionPlayAsideCard                     // play the aside slot card
	ActionPushAside                         // push top card to aside slot
	ActionBuyCard                           // buy a card from the shop
	ActionRemoveTopCard                     // permanently remove top card (shop service)
	ActionRemoveAsideCard                   // permanently remove aside card (shop service)
	ActionClaimHex                          // target hex to claim (follows ActionPlayTopCard for Hex Claim cards)
	ActionForfeit                           // Player field = loser; enqueued by disconnect timer
)

type Action struct {
	Type     ActionType
	Player   PlayerID
	Target   Hex      // hex to claim (ActionClaimHex)
	CardType CardType // card to buy (ActionBuyCard)
}

var (
	ErrInsufficientBar   = errors.New("insufficient bar")
	ErrDeckEmpty         = errors.New("deck is empty")
	ErrAsideSlotEmpty    = errors.New("aside slot is empty")
	ErrAsideSlotOccupied = errors.New("aside slot is occupied")
	ErrCardNotInShop     = errors.New("card not in shop")
	ErrShopEmpty         = errors.New("shop card sold out")
	ErrHexNotFound       = errors.New("hex not found")
	ErrHexOwned          = errors.New("hex already owned")
	ErrNotAdjacent       = errors.New("no adjacent owned hex")
	ErrNotOwner          = errors.New("hex not owned by player")
	ErrPlayerNotFound    = errors.New("player not found")
	ErrCannotDropCapital = errors.New("cannot target capital hex")
	ErrNoFreeHex         = errors.New("no free claimed hex for building")
)

func ProcessActions(state *GameState, actions []Action) {
	for _, a := range actions {
		player := state.Players[a.Player]
		if player == nil {
			continue
		}
		switch a.Type {
		case ActionPlayTopCard:
			applyPlayTopCard(state, player)
		case ActionPlayAsideCard:
			applyPlayAsideCard(state, player)
		case ActionPushAside:
			applyPushAside(player)
		case ActionBuyCard:
			applyBuyCard(state, player, a.CardType)
		case ActionRemoveTopCard:
			applyRemoveTopCard(state, player)
		case ActionRemoveAsideCard:
			applyRemoveAsideCard(state, player)
		case ActionClaimHex:
			applyClaimHex(state, player, a.Target)
		case ActionForfeit:
			forfeitPlayer(state, a.Player)
		}
	}
}

func forfeitPlayer(state *GameState, loser PlayerID) {
	var winner PlayerID
	for pid := range state.Players {
		if pid != loser {
			winner = pid
			break
		}
	}
	triggerVictory(state, winner, "forfeit")
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
