package game

import (
	"math"
	"testing"
)

func TestCounterSpendCap(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	attHex := Hex{Q: 0, R: 0}
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 2}
	state.Hexes[defHex] = &HexState{Owner: p2}

	state.Battles = []Battle{{
		AttackerHex:  attHex,
		DefenderHex:  defHex,
		Attacker:     p1,
		Defender:     p2,
		TimeLeft:     7.0,
		CounterBoost: 2,
	}}

	// Third spend should succeed (boost 2 < cap 3)
	a := Action{Type: ActionCounterSpend, Player: p2, Target: defHex}
	if err := ValidateCounterSpend(state, a); err != nil {
		t.Fatalf("expected valid counter-spend, got: %v", err)
	}
	ApplyCounterSpend(state, a)
	if state.Battles[0].CounterBoost != 3 {
		t.Errorf("boost = %d, want 3", state.Battles[0].CounterBoost)
	}

	// Fourth spend should fail (boost 3 >= cap 3)
	if err := ValidateCounterSpend(state, a); err != ErrCounterSpendCap {
		t.Errorf("expected ErrCounterSpendCap, got: %v", err)
	}
}

func TestCounterSpendTimeCapPreventsSpend(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	attHex := Hex{Q: 0, R: 0}
	defHex := Hex{Q: 1, R: 0}
	state.Hexes[attHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 2}
	state.Hexes[defHex] = &HexState{Owner: p2}

	// TimeLeft=0.9 → int(0.9)=0 → time cap = 0, no spending allowed
	state.Battles = []Battle{{
		AttackerHex: attHex,
		DefenderHex: defHex,
		Attacker:    p1,
		Defender:    p2,
		TimeLeft:    0.9,
	}}

	a := Action{Type: ActionCounterSpend, Player: p2, Target: defHex}
	if err := ValidateCounterSpend(state, a); err != ErrCounterSpendCap {
		t.Errorf("expected ErrCounterSpendCap for 0.9s remaining, got: %v", err)
	}
}

func TestCounterSpendFlipsBattle(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1}
	state.Players[p2] = &Player{ID: p2}

	attHex := Hex{Q: 0, R: 0}
	defHex := Hex{Q: 1, R: 0}
	// Attacker power 3, defender base power 1
	state.Hexes[attHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[defHex] = &HexState{Owner: p2, Capital: true} // Power 1

	// CounterBoost=2: attacker 3 > 1+2=3 is false → defender keeps hex
	state.Battles = []Battle{{
		AttackerHex:  attHex,
		DefenderHex:  defHex,
		Attacker:     p1,
		Defender:     p2,
		TimeLeft:     0.1,
		CounterBoost: 2,
	}}

	RunBattles(state, 0.1)

	if state.Hexes[defHex].Owner != p2 {
		t.Error("expected defender to keep hex with counter-spend boost")
	}
}

func TestUnlockTechDeductsTP(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, TP: 20.0}

	a := Action{Type: ActionUnlockTech, Player: pid, TechID: TechBlitz}
	if err := ValidateUnlockTech(state, a); err != nil {
		t.Fatalf("expected valid unlock, got: %v", err)
	}
	ApplyUnlockTech(state, a)

	if state.Players[pid].TP != 0 {
		t.Errorf("TP = %.1f, want 0", state.Players[pid].TP)
	}
	if !state.Players[pid].Tech[TechBlitz] {
		t.Error("expected TechBlitz to be owned")
	}

	// Double-unlock fails
	if err := ValidateUnlockTech(state, a); err != ErrTechAlreadyOwned {
		t.Errorf("expected ErrTechAlreadyOwned, got: %v", err)
	}
}

func TestUnlockTechInvalidID(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, TP: 1000.0}

	a := Action{Type: ActionUnlockTech, Player: pid, TechID: TechID(99)}
	if err := ValidateUnlockTech(state, a); err != ErrInvalidTechID {
		t.Errorf("expected ErrInvalidTechID, got: %v", err)
	}
}

func TestTPAccumulation(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid}

	h := Hex{Q: 0, R: 0}
	state.Hexes[h] = &HexState{Owner: pid, Building: BuildingResearch, Level: 1}

	// Run 100 ticks = 10 seconds. Research L1: 0.2 TP/sec × 10s = 2.0 TP
	// Note: auto-drop will not fire because maintenance for 1 hex = 1.0/sec, income = 2.0/sec (net positive)
	for range 100 {
		RunEconomy(state, TickDt)
	}

	expected := ResearchPerLevel * 1 * 100 * TickDt // 0.2 * 1 * 100 * 0.1 = 2.0
	if math.Abs(state.Players[pid].TP-expected) > 0.01 {
		t.Errorf("TP = %.3f, want %.3f", state.Players[pid].TP, expected)
	}
}

func TestAutoDropGraceTrigger(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid}

	// 31 hexes: income 62/sec, maintenance 63/sec → net -1/sec
	for i := 0; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	RunAutoDropPhase(state, TickDt)

	p := state.Players[pid]
	if !p.AutoDropActive {
		t.Error("expected AutoDropActive=true after one tick of negative income")
	}
	if p.AutoDropGrace != AutoDropGracePeriod {
		t.Errorf("grace = %.1f, want %.1f", p.AutoDropGrace, AutoDropGracePeriod)
	}
}

func TestAutoDropForceDrop(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid}

	// Gold L1 hex (income 3.9/sec) — should NOT be dropped (highest income)
	goldHex := Hex{Q: 0, R: 0}
	state.Hexes[goldHex] = &HexState{Owner: pid, Building: BuildingGold, Level: 1}

	// 31 base hexes (income 2.0/sec each) — one of these should be dropped
	for i := 1; i <= 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}
	// Total: 32 hexes, income = 3.9 + 31*2 = 65.9/sec, maintenance = 10+20+36 = 66/sec → net -0.1/sec

	// Run 101 ticks: grace triggers at tick 1, fires at tick 101
	for range 101 {
		RunTick(state, TickDt, nil)
	}

	ownedCount := 0
	for _, hs := range state.Hexes {
		if hs.Owner == pid {
			ownedCount++
		}
	}
	if ownedCount != 31 {
		t.Errorf("owned hexes = %d, want 31 (one dropped)", ownedCount)
	}
	// Gold L1 hex must still be owned (it has highest income)
	if state.Hexes[goldHex].Owner != pid {
		t.Error("expected Gold L1 hex to be kept (highest income)")
	}
}

func TestDropHexManualClearsFlag(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 0, AutoDropActive: true}

	// Hex with Gold L1 building: TotalInvested(Gold, 1) * 0.5 = 60*(2-1)*0.5 = 30 refund
	target := Hex{Q: 0, R: 0}
	state.Hexes[target] = &HexState{Owner: pid, Building: BuildingGold, Level: 1}

	a := Action{Type: ActionDropHex, Player: pid, Target: target}
	if err := ValidateDropHex(state, a); err != nil {
		t.Fatalf("expected valid drop, got: %v", err)
	}
	ApplyDropHex(state, a)

	if state.Hexes[target].Owner != NoPlayer {
		t.Error("expected hex to be unclaimed after drop")
	}
	if state.Players[pid].AutoDropActive {
		t.Error("expected AutoDropActive=false after drop")
	}
	expectedRefund := 30.0
	if math.Abs(state.Players[pid].Gold-expectedRefund) > 0.01 {
		t.Errorf("gold = %.2f, want %.2f (refund)", state.Players[pid].Gold, expectedRefund)
	}
}
