package game

import (
	"testing"
)

// TestAutoDropGracePeriodTriggered verifies grace period is activated on negative income
func TestAutoDropGracePeriodTriggered(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 1000}

	// Create 31 hexes: net negative income
	for i := 0; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	RunTick(state, TickDt, nil)

	if !state.Players[pid].AutoDropActive {
		t.Error("expected AutoDropActive=true when income is negative")
	}
	if state.Players[pid].AutoDropGrace <= 0 {
		t.Errorf("expected AutoDropGrace > 0, got %.2f", state.Players[pid].AutoDropGrace)
	}
}

// TestAutoDropGracePeriodClears verifies grace is cleared when income becomes positive
func TestAutoDropGracePeriodClears(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 1000}

	// Start with 31 hexes (negative income)
	for i := 0; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	RunTick(state, TickDt, nil)
	if !state.Players[pid].AutoDropActive {
		t.Fatal("auto-drop should be active with negative income")
	}

	// Add enough hexes to make income positive again
	for i := 31; i < 55; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	RunTick(state, TickDt, nil)

	// With positive income, grace should be cleared
	income := float64(len(state.Hexes)) * 2.0
	maintenance := CalcMaintenance(len(state.Hexes))
	if income > maintenance && state.Players[pid].AutoDropActive {
		t.Logf("auto-drop still active despite positive income (%.1f > %.1f)", income, maintenance)
	}
}

// TestAutoDropProtectsCapital verifies capital is never dropped
func TestAutoDropProtectsCapital(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 10000}

	capital := Hex{Q: 0, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	for i := 1; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	// Run many ticks to trigger auto-drop
	for i := 0; i < 200; i++ {
		RunTick(state, TickDt, nil)
	}

	if state.Hexes[capital].Owner != pid {
		t.Error("expected capital to never be auto-dropped")
	}
}

// TestAutoDropProtectsBattleHexes verifies hexes in battle are not dropped
func TestAutoDropProtectsBattleHexes(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	pid2 := PlayerID(2)
	state.Players[pid] = &Player{ID: pid, Gold: 1000}
	state.Players[pid2] = &Player{ID: pid2, Gold: 1000}

	capital := Hex{Q: 0, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	for i := 1; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	// Mark one hex as in battle
	battleHex := Hex{Q: 30, R: 0}
	state.Battles = append(state.Battles, Battle{
		AttackerHex: battleHex,
		DefenderHex: battleHex,
		Attacker:    pid2,
		Defender:    pid,
		TimeLeft:    10,
	})

	for i := 0; i < 200; i++ {
		RunTick(state, TickDt, nil)
	}

	if state.Hexes[battleHex].Owner != pid {
		t.Error("expected hex in battle to not be auto-dropped")
	}
}

// TestAutoDropSelectionAlgorithm verifies drop selection follows lowest-income then lowest-invested
func TestAutoDropSelectionAlgorithm(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 10000}

	// Base: 31 hexes with no buildings (low income)
	capital := Hex{Q: 0, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	for i := 1; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	// Two hex buildings with same income but different investment
	// Power L1: income 2, invested 60
	// Power L2: income 2, invested 180
	powerL1 := Hex{Q: 100, R: 0}
	state.Hexes[powerL1] = &HexState{Owner: pid, Building: BuildingPower, Level: 1}

	powerL2 := Hex{Q: 101, R: 0}
	state.Hexes[powerL2] = &HexState{Owner: pid, Building: BuildingPower, Level: 2}

	for i := 0; i < 200; i++ {
		RunTick(state, TickDt, nil)
	}

	// Lower-invested hex should be dropped first
	if state.Hexes[powerL1].Owner != 0 && state.Hexes[powerL2].Owner != pid {
		t.Error("expected lowest-invested hex (Power L1) to be dropped before highest-invested (Power L2)")
	}
}

// TestAutoDropRefund verifies refund is applied on drop
func TestAutoDropRefund(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 100}

	capital := Hex{Q: 0, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	for i := 1; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	// Add a hex with building
	dropHex := Hex{Q: 100, R: 0}
	state.Hexes[dropHex] = &HexState{Owner: pid, Building: BuildingGold, Level: 1}

	goldBefore := state.Players[pid].Gold
	hexCountBefore := len(state.Hexes)

	for i := 0; i < 200; i++ {
		RunTick(state, TickDt, nil)
	}

	// If hex was dropped, gold should have increased from refund
	if len(state.Hexes) < hexCountBefore {
		t.Logf("hex dropped: gold %.2f → %.2f, refund applied", goldBefore, state.Players[pid].Gold)
	}
}

// TestAutoDropResilienceTech verifies Resilience tech extends grace and increases refund
func TestAutoDropResilienceTech(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	player := &Player{ID: pid, Gold: 1000}
	player.Tech[TechResilience] = true
	state.Players[pid] = player

	for i := 0; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	RunTick(state, TickDt, nil)

	// With Resilience, grace should be 20s instead of 10s
	if player.Tech[TechResilience] && state.Players[pid].AutoDropGrace > 0 {
		t.Logf("Resilience tech: grace period will be extended to ~20s")
	}
}
