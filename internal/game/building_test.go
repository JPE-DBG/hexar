package game

import (
	"math"
	"testing"
)

func TestUpgradeFromEmptyEconomy(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 200}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid}

	err := ValidateUpgrade(state, Action{Type: ActionUpgrade, Player: pid, Target: hex, Building: BuildingGold})
	if err != nil {
		t.Fatalf("expected valid upgrade, got: %v", err)
	}

	ApplyUpgrade(state, Action{Type: ActionUpgrade, Player: pid, Target: hex, Building: BuildingGold})

	if state.Players[pid].Gold != 140 {
		t.Errorf("gold = %.1f, want 140 (200 - 60)", state.Players[pid].Gold)
	}
	if state.Hexes[hex].Building != BuildingGold {
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

func TestUpgradeFromEmptyDefense(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 200}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid}

	ApplyUpgrade(state, Action{Type: ActionUpgrade, Player: pid, Target: hex, Building: BuildingPower})

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
		// Economy: 60 × 2^level
		{BuildingGold, 0, 60},
		{BuildingGold, 1, 120},
		{BuildingGold, 2, 240},
		{BuildingGold, 3, 480},
		// Defense: 60 × 2^level
		{BuildingPower, 0, 60},
		{BuildingPower, 1, 120},
		{BuildingPower, 2, 240},
		// Research: 80 × 2^level
		{BuildingResearch, 0, 80},
		{BuildingResearch, 1, 160},
		{BuildingResearch, 2, 320},
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
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingPower, Level: 1}

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
	// Economy L2: TotalInvested = 60*(2^2-1) = 180, 50% refund = 90
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingGold, Level: 2}

	ApplyDemolish(state, Action{Type: ActionDemolish, Player: pid, Target: hex})

	expected := 90.0
	if math.Abs(state.Players[pid].Gold-expected) > 0.01 {
		t.Errorf("gold = %.2f, want %.2f", state.Players[pid].Gold, expected)
	}
	if state.Hexes[hex].Building != BuildingNone {
		t.Error("expected no building after demolish")
	}
}

func TestUpgradeEmptyHexRequiresBuilding(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 200}
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid}

	err := ValidateUpgrade(state, Action{Type: ActionUpgrade, Player: pid, Target: hex, Building: BuildingNone})
	if err != ErrNoBuilding {
		t.Errorf("expected ErrNoBuilding for empty hex without building type, got: %v", err)
	}
}

func TestEconomyBuildingIncome(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 0}

	// 1 hex with Economy building at level 2
	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: pid, Building: BuildingGold, Level: 2}

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
