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

	if state.Players[pid].Gold != 120 {
		t.Errorf("gold = %.1f, want 120 (200 - 80)", state.Players[pid].Gold)
	}
	if state.Hexes[hex].Building != BuildingEconomy {
		t.Error("expected Economy building")
	}

	// Verify income: (2 + 0*0.5) * 1.5 = 3.0/sec
	income := HexIncome(state.Hexes[hex])
	if income != 3.0 {
		t.Errorf("hex income = %.2f, want 3.0", income)
	}
}

func TestBuildDefense(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 200}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid}

	ApplyBuild(state, Action{Type: ActionBuild, Player: pid, Target: hex, Building: BuildingDefense})

	if state.Hexes[hex].Power() != 1 {
		t.Errorf("power = %d, want 1 (defense level 0 = +1)", state.Hexes[hex].Power())
	}
}

func TestUpgradeCost(t *testing.T) {
	tests := []struct {
		building BuildingType
		level    int
		want     float64
	}{
		{BuildingEconomy, 0, 40},
		{BuildingEconomy, 1, 80},
		{BuildingEconomy, 2, 160},
		{BuildingEconomy, 3, 320},
		{BuildingDefense, 0, 30},
		{BuildingDefense, 1, 60},
		{BuildingDefense, 2, 120},
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
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingDefense, Level: 0}

	ApplyUpgrade(state, Action{Type: ActionUpgrade, Player: pid, Target: hex})

	if state.Hexes[hex].Level != 1 {
		t.Errorf("level = %d, want 1", state.Hexes[hex].Level)
	}
	if state.Players[pid].Gold != 470 {
		t.Errorf("gold = %.1f, want 470 (500 - 30)", state.Players[pid].Gold)
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
	// Economy building at level 2: cost = 80 + 40 + 80 = 200 total invested, 50% refund = 100
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingEconomy, Level: 2}

	ApplyDemolish(state, Action{Type: ActionDemolish, Player: pid, Target: hex})

	expected := 100.0 // 50% of (80 + 40 + 80)
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

	// Per CLAUDE.md: (2 + 0.5*2) * 1.5 = 4.5/sec, maintenance 1/sec, net 3.5/sec
	for range 100 {
		RunTick(state, TickDt, nil)
	}

	// 10 seconds * 3.5/sec = 35
	expected := 35.0
	if math.Abs(state.Players[pid].Gold-expected) > 0.1 {
		t.Errorf("gold = %.2f, want %.2f", state.Players[pid].Gold, expected)
	}
}
