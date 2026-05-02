package game

import (
	"math"
	"testing"
)

func TestBuildEconomy(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 200}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid}

	err := ValidateBuild(state, Action{Type: ActionBuild, Player: pid, Target: hex, Building: BuildingEconomy})
	if err != nil {
		t.Fatalf("expected valid build, got: %v", err)
	}

	ApplyBuild(state, Action{Type: ActionBuild, Player: pid, Target: hex, Building: BuildingEconomy})

	if state.Players[pid].Gold != 140 {
		t.Errorf("gold = %.1f, want 140 (200 - 60)", state.Players[pid].Gold)
	}
	if state.Hexes[hex].Building != BuildingEconomy {
		t.Error("expected Economy building")
	}
	if state.Hexes[hex].Level != 1 {
		t.Errorf("level = %d, want 1", state.Hexes[hex].Level)
	}

	// Verify income: (2 + 1*0.6) * 1.5 = 3.9/sec
	income := HexIncome(state.Hexes[hex])
	if math.Abs(income-3.9) > 0.001 {
		t.Errorf("hex income = %.4f, want 3.9", income)
	}
}

func TestBuildDefense(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 200}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid}

	ApplyBuild(state, Action{Type: ActionBuild, Player: pid, Target: hex, Building: BuildingDefense})

	if state.Hexes[hex].Level != 1 {
		t.Errorf("level = %d, want 1", state.Hexes[hex].Level)
	}
	if state.Hexes[hex].Power() != 1 {
		t.Errorf("power = %d, want 1 (defense level 1)", state.Hexes[hex].Power())
	}
}

func TestUpgradeCost(t *testing.T) {
	tests := []struct {
		building BuildingType
		level    int
		want     float64
	}{
		{BuildingEconomy, 1, 60},
		{BuildingEconomy, 2, 120},
		{BuildingEconomy, 3, 240},
		{BuildingDefense, 1, 120},
		{BuildingDefense, 2, 240},
		{BuildingDefense, 3, 480},
	}
	for _, tt := range tests {
		got := UpgradeCost(tt.building, tt.level)
		if got != tt.want {
			t.Errorf("UpgradeCost(%d, %d) = %.0f, want %.0f", tt.building, tt.level, got, tt.want)
		}
	}
}

func TestUpgradeIncreasesLevel(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 500}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingDefense, Level: 1}

	ApplyUpgrade(state, Action{Type: ActionUpgrade, Player: pid, Target: hex})

	if state.Hexes[hex].Level != 2 {
		t.Errorf("level = %d, want 2", state.Hexes[hex].Level)
	}
	if state.Players[pid].Gold != 380 {
		t.Errorf("gold = %.1f, want 380 (500 - 120)", state.Players[pid].Gold)
	}
	if state.Hexes[hex].Power() != 2 {
		t.Errorf("power = %d, want 2", state.Hexes[hex].Power())
	}
}

func TestDemolishRefund(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 0}
	hex := Hex{Q: 0, R: 0}
	// Economy building at level 2: TotalInvested = 60*2^(2-1) = 120, 50% refund = 60
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingEconomy, Level: 2}

	ApplyDemolish(state, Action{Type: ActionDemolish, Player: pid, Target: hex})

	expected := 60.0
	if math.Abs(state.Players[pid].Gold-expected) > 0.01 {
		t.Errorf("gold = %.2f, want %.2f", state.Players[pid].Gold, expected)
	}
	if state.Hexes[hex].Building != BuildingNone {
		t.Error("expected no building after demolish")
	}
}

func TestBuildOnOccupiedHex(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 200}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingDefense}

	err := ValidateBuild(state, Action{Type: ActionBuild, Player: pid, Target: hex, Building: BuildingEconomy})
	if err != ErrHasBuilding {
		t.Errorf("expected ErrHasBuilding, got: %v", err)
	}
}

func TestEconomyBuildingIncome(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 0}

	// 1 hex with Economy building at level 2
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingEconomy, Level: 2}

	// Per CLAUDE.md: (2 + 0.6*2) * 1.5 = 4.8/sec, maintenance 1/sec, net 3.8/sec
	for range 100 {
		RunTick(state, TickDt, nil)
	}

	// 10 seconds * 3.8/sec = 38
	expected := 38.0
	if math.Abs(state.Players[pid].Gold-expected) > 0.1 {
		t.Errorf("gold = %.2f, want %.2f", state.Players[pid].Gold, expected)
	}
}
