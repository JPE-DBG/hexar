package game

import "testing"

func TestAutoDropGracePeriodTriggered(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 1000}

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

func TestAutoDropGracePeriodClears(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 1000}

	for i := 0; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	RunTick(state, TickDt, nil)
	if !state.Players[pid].AutoDropActive {
		t.Fatal("auto-drop should be active with negative income")
	}

	// Remove hexes until income is positive (keep 10: net +10/sec)
	for i := 10; i < 31; i++ {
		delete(state.Hexes, Hex{Q: i, R: 0})
	}

	RunTick(state, TickDt, nil)

	if state.Players[pid].AutoDropActive {
		t.Error("expected AutoDropActive=false when income is positive")
	}
	if state.Players[pid].AutoDropGrace != 0 {
		t.Errorf("expected AutoDropGrace=0, got %.2f", state.Players[pid].AutoDropGrace)
	}
}

func TestAutoDropProtectsCapital(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 10000}

	capital := Hex{Q: 0, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}
	for i := 1; i < 31; i++ {
		state.Hexes[Hex{Q: i, R: 0}] = &HexState{Owner: pid}
	}

	for i := 0; i < 200; i++ {
		RunTick(state, TickDt, nil)
	}

	if state.Hexes[capital].Owner != pid {
		t.Error("expected capital to never be auto-dropped")
	}
}

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

	battleHex := Hex{Q: 30, R: 0}
	attackerHex := Hex{Q: 31, R: 0}
	state.Hexes[attackerHex] = &HexState{Owner: pid2}
	state.Battles = append(state.Battles, Battle{
		AttackerHex: attackerHex,
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

func TestAutoDropSelectionAlgorithm(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 0}

	capital := Hex{Q: 0, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}

	// Two power hexes: same income (2/sec), different investment
	powerL1 := Hex{Q: 1, R: 0}
	state.Hexes[powerL1] = &HexState{Owner: pid, Building: BuildingPower, Level: 1} // invested 60

	powerL2 := Hex{Q: 2, R: 0}
	state.Hexes[powerL2] = &HexState{Owner: pid, Building: BuildingPower, Level: 2} // invested 180

	// Call autoDropLowestHex directly: capital is protected, powerL1 (60) < powerL2 (180)
	autoDropLowestHex(state, pid)

	if state.Hexes[powerL1].Owner == pid {
		t.Error("expected powerL1 (lower investment) to be dropped, not powerL2")
	}
	if state.Hexes[powerL2].Owner != pid {
		t.Error("expected powerL2 (higher investment) to be retained")
	}
}

func TestAutoDropRefund(t *testing.T) {
	state := NewGameState()
	pid := PlayerID(1)
	state.Players[pid] = &Player{ID: pid, Gold: 100}

	capital := Hex{Q: 0, R: 0}
	state.Hexes[capital] = &HexState{Owner: pid, Capital: true}

	dropHex := Hex{Q: 1, R: 0}
	state.Hexes[dropHex] = &HexState{Owner: pid, Building: BuildingGold, Level: 1}

	goldBefore := state.Players[pid].Gold
	autoDropLowestHex(state, pid)

	expectedRefund := TotalInvested(BuildingGold, 1) * AutoDropRefund // 60 * 0.5 = 30
	if state.Players[pid].Gold != goldBefore+expectedRefund {
		t.Errorf("gold = %.2f, want %.2f (refund=%.2f)", state.Players[pid].Gold, goldBefore+expectedRefund, expectedRefund)
	}
	if state.Hexes[dropHex].Owner != NoPlayer {
		t.Error("expected dropped hex to be unclaimed")
	}
}

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

	// Grace period should be ResilienceGracePeriod (20s), not AutoDropGracePeriod (10s)
	if state.Players[pid].AutoDropGrace != ResilienceGracePeriod {
		t.Errorf("expected AutoDropGrace=%.1f with Resilience, got %.2f", ResilienceGracePeriod, state.Players[pid].AutoDropGrace)
	}

	// Verify 70% refund: use a fresh state with only capital + building hex
	state2 := NewGameState()
	state2.Players[pid] = &Player{ID: pid, Gold: 100}
	state2.Players[pid].Tech[TechResilience] = true
	state2.Hexes[Hex{Q: 0, R: 0}] = &HexState{Owner: pid, Capital: true}
	dropHex := Hex{Q: 1, R: 0}
	state2.Hexes[dropHex] = &HexState{Owner: pid, Building: BuildingGold, Level: 1}

	goldBefore := state2.Players[pid].Gold
	autoDropLowestHex(state2, pid)

	expectedRefund := TotalInvested(BuildingGold, 1) * ResilienceDropRefund // 60 * 0.7 = 42
	if state2.Players[pid].Gold != goldBefore+expectedRefund {
		t.Errorf("Resilience refund: gold = %.2f, want %.2f (refund=%.2f)", state2.Players[pid].Gold, goldBefore+expectedRefund, expectedRefund)
	}
}
