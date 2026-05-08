package game

import "testing"

func TestVictoryCapitalCapture(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	cap2 := Hex{Q: 5, R: 0}
	hex2a := Hex{Q: 4, R: 0}
	hex2b := Hex{Q: 6, R: 0}
	cap1 := Hex{Q: 0, R: 0}
	state.Hexes[cap2] = &HexState{Owner: p2, Capital: true}
	state.Hexes[hex2a] = &HexState{Owner: p2, Building: BuildingPower, Level: 2}
	state.Hexes[hex2b] = &HexState{Owner: p2}
	state.Hexes[cap1] = &HexState{Owner: p1, Capital: true}

	// Also add an ongoing battle involving p2 to test cancellation
	state.Battles = []Battle{{
		AttackerHex: hex2a, DefenderHex: cap1,
		Attacker: p2, Defender: p1, TimeLeft: 5.0,
	}}

	TriggerVictory(state, p1)

	if !state.Over {
		t.Error("state.Over should be true after TriggerVictory")
	}
	if state.Winner != p1 {
		t.Errorf("Winner should be p1, got %v", state.Winner)
	}

	// All loser (p2) hexes should be unclaimed
	for h, hs := range state.Hexes {
		if hs.Owner == p2 {
			t.Errorf("Hex %v should be unclaimed after victory, still owned by p2", h)
		}
	}

	// Battle involving loser should be cancelled
	for _, b := range state.Battles {
		if b.Attacker == p2 || b.Defender == p2 {
			t.Error("Battles involving loser should be cancelled")
		}
	}
}

func TestVictoryIdempotent(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p2] = &Player{ID: p2, Gold: 0}
	state.Hexes[Hex{Q: 0, R: 0}] = &HexState{Owner: p1, Capital: true}
	state.Hexes[Hex{Q: 5, R: 0}] = &HexState{Owner: p2, Capital: true}

	TriggerVictory(state, p1)
	TriggerVictory(state, p2) // second call should be no-op

	if state.Winner != p1 {
		t.Error("Second TriggerVictory should not overwrite winner")
	}
}

func TestVictoryViaCapitalAttack(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	adjHex := Hex{Q: -1, R: 0}
	enemyCap := Hex{Q: 0, R: 0}
	ownHex := Hex{Q: 0, R: 1}

	// p1 has adj hex with Power 10; p2 capital has Power 0 → diff > 3, instant takeover
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 10}
	state.Hexes[enemyCap] = &HexState{Owner: p2, Capital: true}
	state.Hexes[ownHex] = &HexState{Owner: p2}

	a := Action{Type: ActionAttack, Player: p1, Target: enemyCap}
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("attack on capital should validate: %v", err)
	}
	ApplyAttack(state, a)

	if !state.Over {
		t.Error("Game should be over after capturing enemy capital")
	}
	if state.Winner != p1 {
		t.Errorf("Winner should be p1, got %v", state.Winner)
	}
	// Enemy's remaining hexes should be unclaimed
	if state.Hexes[ownHex].Owner != NoPlayer {
		t.Error("Loser's non-capital hex should become unclaimed")
	}
}

func TestRunTickStopsWhenOver(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 100}
	state.Hexes[Hex{Q: 0, R: 0}] = &HexState{Owner: p1, Capital: true}
	state.Over = true

	goldBefore := state.Players[p1].Gold
	RunTick(state, TickDt, nil)
	// Economy should not run (gold unchanged)
	if state.Players[p1].Gold != goldBefore {
		t.Error("RunTick should not process economy when game is over")
	}
}
