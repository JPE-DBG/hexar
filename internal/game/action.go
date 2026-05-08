package game

import "errors"

type ActionType int

const (
	ActionClaim ActionType = iota
	ActionUpgrade
	ActionDemolish
	ActionAttack
	ActionCounterSpend
	ActionUnlockTech
	ActionDropHex
	ActionFortify
)

type Action struct {
	Type     ActionType
	Player   PlayerID
	Target   Hex
	Building BuildingType
	TechID   TechID
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
	ErrNotDefender       = errors.New("player is not the defender")
	ErrNoBattleOnHex     = errors.New("no active battle on hex")
	ErrCounterSpendCap   = errors.New("counter-spend cap reached")
	ErrInsufficientTP    = errors.New("insufficient tech points")
	ErrTechAlreadyOwned  = errors.New("tech already owned")
	ErrInvalidTechID     = errors.New("invalid tech ID")
	ErrCannotDropCapital = errors.New("cannot drop capital hex")
	ErrTechNotOwned      = errors.New("tech not owned")
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
		case ActionCounterSpend:
			if ValidateCounterSpend(state, a) == nil {
				ApplyCounterSpend(state, a)
			}
		case ActionUnlockTech:
			if ValidateUnlockTech(state, a) == nil {
				ApplyUnlockTech(state, a)
			}
		case ActionDropHex:
			if ValidateDropHex(state, a) == nil {
				ApplyDropHex(state, a)
			}
		case ActionFortify:
			if ValidateFortify(state, a) == nil {
				ApplyFortify(state, a)
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
	if !player.Tech[TechBlitz] && player.Gold < ClaimCost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyClaim(state *GameState, action Action) {
	player := state.Players[action.Player]
	if !player.Tech[TechBlitz] {
		player.Gold -= ClaimCost
	}
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

func findBattleByDefenderHex(state *GameState, target Hex) *Battle {
	for i := range state.Battles {
		if state.Battles[i].DefenderHex == target {
			return &state.Battles[i]
		}
	}
	return nil
}

func ValidateCounterSpend(state *GameState, action Action) error {
	b := findBattleByDefenderHex(state, action.Target)
	if b == nil {
		return ErrNoBattleOnHex
	}
	if b.Defender != action.Player {
		return ErrNotDefender
	}
	player := state.Players[action.Player]
	if player == nil {
		return ErrPlayerNotFound
	}
	if player.Gold < CounterSpendCostPerSec {
		return ErrInsufficientGold
	}
	timeCap := int(b.TimeLeft)
	garrisonBoost := garrisonBoostForDefender(state, action.Player, action.Target)
	effectiveCap := max(0, min(CounterSpendCap, timeCap)-garrisonBoost)
	if b.CounterBoost >= effectiveCap {
		return ErrCounterSpendCap
	}
	return nil
}

func ApplyCounterSpend(state *GameState, action Action) {
	b := findBattleByDefenderHex(state, action.Target)
	state.Players[action.Player].Gold -= CounterSpendCostPerSec
	b.CounterBoost++
}

func ValidateUnlockTech(state *GameState, action Action) error {
	if action.TechID < 0 || action.TechID >= TechCount {
		return ErrInvalidTechID
	}
	player := state.Players[action.Player]
	if player == nil {
		return ErrPlayerNotFound
	}
	if player.Tech[action.TechID] {
		return ErrTechAlreadyOwned
	}
	if player.TP < TechCost(action.TechID) {
		return ErrInsufficientTP
	}
	return nil
}

func ApplyUnlockTech(state *GameState, action Action) {
	player := state.Players[action.Player]
	player.TP -= TechCost(action.TechID)
	player.Tech[action.TechID] = true
}

func ValidateDropHex(state *GameState, action Action) error {
	hs, ok := state.Hexes[action.Target]
	if !ok {
		return ErrHexNotFound
	}
	if hs.Owner != action.Player {
		return ErrNotOwner
	}
	if hs.Capital {
		return ErrCannotDropCapital
	}
	player := state.Players[action.Player]
	if player == nil {
		return ErrPlayerNotFound
	}
	if hasBattleOnHex(state, action.Target) {
		return ErrBattleInProgress
	}
	return nil
}

func ApplyDropHex(state *GameState, action Action) {
	hs := state.Hexes[action.Target]
	player := state.Players[action.Player]
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
	player.AutoDropActive = false
}

func ValidateFortify(state *GameState, action Action) error {
	hs, ok := state.Hexes[action.Target]
	if !ok {
		return ErrHexNotFound
	}
	if hs.Owner != action.Player {
		return ErrNotOwner
	}
	player := state.Players[action.Player]
	if player == nil {
		return ErrPlayerNotFound
	}
	if !player.Tech[TechFortify] {
		return ErrTechNotOwned
	}
	if player.Gold < FortifyCost {
		return ErrInsufficientGold
	}
	return nil
}

func ApplyFortify(state *GameState, action Action) {
	state.Players[action.Player].Gold -= FortifyCost
	state.Hexes[action.Target].FortifyTimer = FortifyDuration
}
