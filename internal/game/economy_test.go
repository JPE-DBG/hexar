package game

import (
	"math"
	"testing"
)

func TestEconomy10Hexes30Seconds(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid}

	for i := 0; i < 10; i++ {
		h := Hex{Q: i, R: 0}
		state.Hexes[h] = &HexState{Owner: pid}
	}

	for range 300 {
		RunTick(state, TickDt, nil)
	}

	// 10 hexes: income 20/sec, maintenance 10/sec, net +10/sec × 30s = 300
	expected := 300.0
	if math.Abs(state.Players[pid].Gold-expected) > 0.1 {
		t.Errorf("gold = %.2f, want %.2f", state.Players[pid].Gold, expected)
	}
}

func TestMaintenanceStepped(t *testing.T) {
	tests := []struct {
		hexes    int
		wantMnt  float64
	}{
		{5, 5.0},
		{10, 10.0},
		{15, 20.0},
		{20, 30.0},
		{25, 45.0},
	}
	for _, tt := range tests {
		got := CalcMaintenance(tt.hexes)
		if got != tt.wantMnt {
			t.Errorf("CalcMaintenance(%d) = %.1f, want %.1f", tt.hexes, got, tt.wantMnt)
		}
	}
}

func TestClaimValidation(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 100}

	capital := Hex{Q: 0, R: 0}
	target := Hex{Q: 1, R: 0}
	farHex := Hex{Q: 3, R: 3}

	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	state.Hexes[target] = &HexState{}
	state.Hexes[farHex] = &HexState{}

	// Valid claim
	err := ValidateClaim(state, ClaimAction{Player: pid, Target: target})
	if err != nil {
		t.Errorf("expected valid claim, got: %v", err)
	}

	// Not adjacent
	err = ValidateClaim(state, ClaimAction{Player: pid, Target: farHex})
	if err != ErrNotAdjacent {
		t.Errorf("expected ErrNotAdjacent, got: %v", err)
	}

	// Already owned
	state.Hexes[target].Owner = pid
	err = ValidateClaim(state, ClaimAction{Player: pid, Target: target})
	if err != ErrHexOwned {
		t.Errorf("expected ErrHexOwned, got: %v", err)
	}
}

func TestClaimDeductsGold(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 50}

	capital := Hex{Q: 0, R: 0}
	target := Hex{Q: 1, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	state.Hexes[target] = &HexState{}

	RunTick(state, TickDt, []ClaimAction{{Player: pid, Target: target}})

	if state.Players[pid].Gold < 39.0 || state.Players[pid].Gold > 41.0 {
		t.Errorf("gold = %.2f, want ~40 (50 - 10 + tick income)", state.Players[pid].Gold)
	}
	if state.Hexes[target].Owner != pid {
		t.Error("expected target to be owned by player")
	}
}

func TestInsufficientGold(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 5}

	capital := Hex{Q: 0, R: 0}
	target := Hex{Q: 1, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	state.Hexes[target] = &HexState{}

	err := ValidateClaim(state, ClaimAction{Player: pid, Target: target})
	if err != ErrInsufficientGold {
		t.Errorf("expected ErrInsufficientGold, got: %v", err)
	}
}

func TestGoldClampsAtZero(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 0}

	// 25 hexes: income 50/sec, maintenance 45/sec, net +5/sec — but with Gold=0 start it won't go negative
	// Use 31 hexes for negative: income 62/sec, maintenance 63/sec
	for i := 0; i < 31; i++ {
		h := Hex{Q: i, R: 0}
		state.Hexes[h] = &HexState{Owner: pid}
	}

	for range 100 {
		RunTick(state, TickDt, nil)
	}

	if state.Players[pid].Gold < 0 {
		t.Errorf("gold = %.2f, should never go negative", state.Players[pid].Gold)
	}
}
