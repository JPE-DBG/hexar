package game

import (
	"testing"
)

func TestAttackInstantTakeover(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	p2 := PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 200}
	state.Players[p2] = &Player{ID: p2, Gold: 200}

	// Attacker hex with Defense level 4 = Power 5
	attHex := Hex{Q: 0, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Building: BuildingDefense, Level: 4}

	// Defender hex with no building = Power 0, adjacent to attacker
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[defHex] = &HexState{Owner: p2}

	// Power diff = 5 - 0 = 5 > 3, should be instant
	action := Action{Type: ActionAttack, Player: p1, Target: defHex}
	err := ValidateAttack(state, action)
	if err != nil {
		t.Fatalf("expected valid attack, got: %v", err)
	}

	ApplyAttack(state, action)

	if state.Hexes[defHex].Owner != p1 {
		t.Error("expected instant takeover")
	}
	if len(state.Battles) != 0 {
		t.Error("expected no battle for instant takeover")
	}
}

func TestAttackStartsBattle(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	p2 := PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 200}
	state.Players[p2] = &Player{ID: p2, Gold: 200}

	// Attacker: Defense level 1 = Power 2
	attHex := Hex{Q: 0, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Building: BuildingDefense, Level: 1}

	// Defender: Capital = Power 1
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[defHex] = &HexState{Owner: p2, Capital: true}

	// Power diff = 2 - 1 = 1, starts battle
	ApplyAttack(state, Action{Type: ActionAttack, Player: p1, Target: defHex})

	if len(state.Battles) != 1 {
		t.Fatalf("expected 1 battle, got %d", len(state.Battles))
	}
	// Duration = 5 + (2+1)/2 = 6.5
	if state.Battles[0].TimeLeft != 6.5 {
		t.Errorf("battle time = %.1f, want 6.5", state.Battles[0].TimeLeft)
	}
	if state.Players[p1].Gold != 100 {
		t.Errorf("gold = %.1f, want 100 (200 - 100)", state.Players[p1].Gold)
	}
}

func TestBattleResolution(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	p2 := PlayerID(2)
	state.Players[p1] = &Player{ID: p1}
	state.Players[p2] = &Player{ID: p2}

	attHex := Hex{Q: 0, R: 0}
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Building: BuildingDefense, Level: 1} // Power 2
	state.Hexes[defHex] = &HexState{Owner: p2, Capital: true}                       // Power 1

	state.Battles = []Battle{{
		AttackerHex: attHex,
		DefenderHex: defHex,
		Attacker:    p1,
		Defender:    p2,
		TimeLeft:    0.1,
	}}

	RunBattles(state, 0.1)

	if len(state.Battles) != 0 {
		t.Error("expected battle to resolve")
	}
	if state.Hexes[defHex].Owner != p1 {
		t.Error("expected attacker to win (power 2 > 1)")
	}
}

func TestBattleDefenderWins(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	p2 := PlayerID(2)
	state.Players[p1] = &Player{ID: p1}
	state.Players[p2] = &Player{ID: p2}

	attHex := Hex{Q: 0, R: 0}
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Capital: true}                       // Power 1
	state.Hexes[defHex] = &HexState{Owner: p2, Building: BuildingDefense, Level: 1} // Power 2

	state.Battles = []Battle{{
		AttackerHex: attHex,
		DefenderHex: defHex,
		Attacker:    p1,
		Defender:    p2,
		TimeLeft:    0.1,
	}}

	RunBattles(state, 0.1)

	if state.Hexes[defHex].Owner != p2 {
		t.Error("expected defender to keep hex (power 2 > 1)")
	}
}

func TestAttackInsufficientPower(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	p2 := PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 200}
	state.Players[p2] = &Player{ID: p2}

	attHex := Hex{Q: 0, R: 0}
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Capital: true}                       // Power 1
	state.Hexes[defHex] = &HexState{Owner: p2, Building: BuildingDefense, Level: 0} // Power 1

	// Power 1 vs 1: not strictly greater
	err := ValidateAttack(state, Action{Type: ActionAttack, Player: p1, Target: defHex})
	if err != ErrInsufficientPower {
		t.Errorf("expected ErrInsufficientPower, got: %v", err)
	}
}

func TestAttackCost(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	p2 := PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 150}
	state.Players[p2] = &Player{ID: p2}

	attHex := Hex{Q: 0, R: 0}
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Building: BuildingDefense, Level: 4} // Power 5
	state.Hexes[defHex] = &HexState{Owner: p2}                                      // Power 0

	ApplyAttack(state, Action{Type: ActionAttack, Player: p1, Target: defHex})

	if state.Players[p1].Gold != 50 {
		t.Errorf("gold = %.1f, want 50 (150 - 100)", state.Players[p1].Gold)
	}
}
