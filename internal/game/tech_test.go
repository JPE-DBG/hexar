package game

import (
	"testing"
)

// helpers shared across tech tests
func techState() *GameState {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500, TP: 1000}
	state.Players[p2] = &Player{ID: p2, Gold: 500, TP: 1000}

	capHex := Hex{Q: 0, R: 0}
	adjHex := Hex{Q: 1, R: 0}
	enemyCap := Hex{Q: 5, R: 0}
	enemyAdj := Hex{Q: 4, R: 0}

	state.Hexes[capHex] = &HexState{Owner: p1, Capital: true}
	state.Hexes[adjHex] = &HexState{Owner: p1}
	state.Hexes[enemyCap] = &HexState{Owner: p2, Capital: true}
	state.Hexes[enemyAdj] = &HexState{Owner: p2}
	return state
}

func TestBlitz(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechBlitz] = true

	capHex := Hex{Q: 0, R: 0}
	unclaimedHex := Hex{Q: 1, R: 0}
	state.Hexes[capHex] = &HexState{Owner: p1, Capital: true}
	state.Hexes[unclaimedHex] = &HexState{}

	// With Blitz and 0 gold, claim should succeed (free)
	a := Action{Type: ActionClaim, Player: p1, Target: unclaimedHex}
	if err := ValidateClaim(state, a); err != nil {
		t.Fatalf("Blitz claim with 0g should succeed, got: %v", err)
	}
	ApplyClaim(state, a)
	if state.Players[p1].Gold != 0 {
		t.Errorf("Blitz claim should cost 0g, gold = %v", state.Players[p1].Gold)
	}
	if state.Hexes[unclaimedHex].Owner != p1 {
		t.Error("Blitz claim should transfer ownership")
	}
}

func TestBlitzWithoutTech(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0} // no Blitz

	capHex := Hex{Q: 0, R: 0}
	unclaimedHex := Hex{Q: 1, R: 0}
	state.Hexes[capHex] = &HexState{Owner: p1, Capital: true}
	state.Hexes[unclaimedHex] = &HexState{}

	a := Action{Type: ActionClaim, Player: p1, Target: unclaimedHex}
	if err := ValidateClaim(state, a); err != ErrInsufficientGold {
		t.Errorf("expected ErrInsufficientGold without Blitz at 0g, got: %v", err)
	}
}

func TestReclamation(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 200}
	state.Players[p1].Tech[TechReclamation] = true
	state.Players[p2] = &Player{ID: p2, Gold: 200}

	adjHex := Hex{Q: -1, R: 0}
	targetHex := Hex{Q: 0, R: 0}
	// targetHex was previously owned by p1
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[targetHex] = &HexState{Owner: p2, PreviousOwner: p1}

	a := Action{Type: ActionAttack, Player: p1, Target: targetHex}
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("Reclamation attack should validate: %v", err)
	}
	goldBefore := state.Players[p1].Gold
	ApplyAttack(state, a)
	spent := goldBefore - state.Players[p1].Gold
	if spent != ReclamationCost {
		t.Errorf("Reclamation should cost %v, spent %v", ReclamationCost, spent)
	}
}

func TestReclamationFallback(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 200}
	state.Players[p1].Tech[TechReclamation] = true
	state.Players[p2] = &Player{ID: p2, Gold: 200}

	adjHex := Hex{Q: -1, R: 0}
	targetHex := Hex{Q: 0, R: 0}
	// targetHex was NOT previously owned by p1
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[targetHex] = &HexState{Owner: p2, PreviousOwner: p2}

	a := Action{Type: ActionAttack, Player: p1, Target: targetHex}
	goldBefore := state.Players[p1].Gold
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("attack should validate: %v", err)
	}
	ApplyAttack(state, a)
	spent := goldBefore - state.Players[p1].Gold
	if spent != AttackCost {
		t.Errorf("Without matching PreviousOwner, should cost %v, spent %v", AttackCost, spent)
	}
}

func TestVanguard(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500, VanguardTimer: VanguardWindow}
	state.Players[p1].Tech[TechVanguard] = true
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	adjHex := Hex{Q: -1, R: 0}
	targetHex := Hex{Q: 0, R: 0}
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[targetHex] = &HexState{Owner: p2}

	a := Action{Type: ActionAttack, Player: p1, Target: targetHex}
	goldBefore := state.Players[p1].Gold
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("Vanguard attack should validate: %v", err)
	}
	ApplyAttack(state, a)
	spent := goldBefore - state.Players[p1].Gold
	if spent != VanguardAttackCost {
		t.Errorf("Vanguard active should cost %v, spent %v", VanguardAttackCost, spent)
	}
}

func TestVanguardExpired(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500, VanguardTimer: 0} // expired
	state.Players[p1].Tech[TechVanguard] = true
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	adjHex := Hex{Q: -1, R: 0}
	targetHex := Hex{Q: 0, R: 0}
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[targetHex] = &HexState{Owner: p2}

	a := Action{Type: ActionAttack, Player: p1, Target: targetHex}
	goldBefore := state.Players[p1].Gold
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("attack should validate: %v", err)
	}
	ApplyAttack(state, a)
	spent := goldBefore - state.Players[p1].Gold
	if spent != AttackCost {
		t.Errorf("Expired Vanguard should cost %v, spent %v", AttackCost, spent)
	}
}

func TestReclamationVanguardStacking(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 1000, VanguardTimer: 5.0} // Active Vanguard
	state.Players[p1].Tech[TechReclamation] = true
	state.Players[p1].Tech[TechVanguard] = true
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	adjHex := Hex{Q: -1, R: 0}
	targetHex := Hex{Q: 0, R: 0}
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[targetHex] = &HexState{Owner: p2, PreviousOwner: p1} // Was owned by p1

	// Attack cost should be 0 (both discounts apply)
	cost := effectiveAttackCost(state, p1, targetHex)
	if cost != 0 {
		t.Errorf("Expected cost 0 with both Reclamation + Vanguard, got %.1f", cost)
	}

	// Apply attack, verify gold deduction
	a := Action{Type: ActionAttack, Player: p1, Target: targetHex}
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("attack should validate: %v", err)
	}
	ApplyAttack(state, a)
	if state.Players[p1].Gold != 1000 {
		t.Errorf("Expected 1000 gold (0g attack), got %.1f", state.Players[p1].Gold)
	}
}

func TestWarChest(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechWarChest] = true
	state.Players[p2] = &Player{ID: p2, Gold: 0}

	targetHex := Hex{Q: 1, R: 0}
	state.Hexes[targetHex] = &HexState{Owner: p2}

	goldBefore := state.Players[p1].Gold
	transferHex(state, targetHex, p1)
	goldAfter := state.Players[p1].Gold
	if goldAfter-goldBefore != WarChestRefund {
		t.Errorf("War Chest should refund %v g on capture, got %v", WarChestRefund, goldAfter-goldBefore)
	}
}

func TestWarChestNoRefundFromUnclaimed(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechWarChest] = true

	targetHex := Hex{Q: 1, R: 0}
	state.Hexes[targetHex] = &HexState{Owner: NoPlayer} // unclaimed

	goldBefore := state.Players[p1].Gold
	transferHex(state, targetHex, p1)
	if state.Players[p1].Gold != goldBefore {
		t.Error("War Chest should not refund when capturing unclaimed hex")
	}
}

func TestGarrisonBoost(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p2] = &Player{ID: p2, Gold: 500}
	state.Players[p2].Tech[TechGarrison] = true

	// Defender hex at (0,0), attacker hex at (-1,0)
	// Two adjacent owned hexes for defender: (1,0) and (0,1) → garrison boost = 2
	defHex := Hex{Q: 0, R: 0}
	atkHex := Hex{Q: -1, R: 0}
	adj1 := Hex{Q: 1, R: 0}
	adj2 := Hex{Q: 0, R: 1}
	state.Hexes[defHex] = &HexState{Owner: p2}
	state.Hexes[atkHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[adj1] = &HexState{Owner: p2}
	state.Hexes[adj2] = &HexState{Owner: p2}

	boost := garrisonBoostForDefender(state, p2, defHex)
	if boost != GarrisonMaxBoost {
		t.Errorf("Garrison boost should be %d with 2+ adjacent hexes, got %d", GarrisonMaxBoost, boost)
	}
}

func TestGarrisonReducesCounterSpendCap(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p2] = &Player{ID: p2, Gold: 500}
	state.Players[p2].Tech[TechGarrison] = true

	defHex := Hex{Q: 0, R: 0}
	atkHex := Hex{Q: -1, R: 0}
	adj1 := Hex{Q: 1, R: 0}
	adj2 := Hex{Q: 0, R: 1}
	state.Hexes[defHex] = &HexState{Owner: p2}
	state.Hexes[atkHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}
	state.Hexes[adj1] = &HexState{Owner: p2}
	state.Hexes[adj2] = &HexState{Owner: p2}

	// Garrison = +2, so CounterSpend cap = 3 - 2 = 1
	state.Battles = []Battle{{
		AttackerHex: atkHex, DefenderHex: defHex,
		Attacker: p1, Defender: p2, TimeLeft: 10.0, CounterBoost: 0,
	}}

	a := Action{Type: ActionCounterSpend, Player: p2, Target: defHex}
	// First spend should be allowed (CounterBoost 0 < effectiveCap 1)
	if err := ValidateCounterSpend(state, a); err != nil {
		t.Fatalf("first counter-spend should succeed: %v", err)
	}
	ApplyCounterSpend(state, a)

	// Second spend should be rejected (CounterBoost 1 >= effectiveCap 1)
	if err := ValidateCounterSpend(state, a); err == nil {
		t.Error("second counter-spend should be rejected when garrison fills cap")
	}
}

func TestIronGrip(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechIronGrip] = true

	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: p1} // Power 0 normally

	power := hexEffectivePower(state, p1, hex)
	if power != 1 {
		t.Errorf("Iron Grip should give +1 effective power, got %d", power)
	}
}

func TestIronGripCapital(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechIronGrip] = true

	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: p1, Capital: true} // Power 1 from capital

	power := hexEffectivePower(state, p1, hex)
	if power != 2 {
		t.Errorf("Iron Grip on capital should give power 2 (1+1), got %d", power)
	}
}

func TestSiegeMasteryDuration(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p1].Tech[TechSiegeMastery] = true
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	// Power diff = 2 → standard battle
	adjHex := Hex{Q: -1, R: 0}
	targetHex := Hex{Q: 0, R: 0}
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 3}    // Power 3
	state.Hexes[targetHex] = &HexState{Owner: p2, Building: BuildingPower, Level: 1} // Power 1

	a := Action{Type: ActionAttack, Player: p1, Target: targetHex}
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("attack should validate: %v", err)
	}
	ApplyAttack(state, a)

	if len(state.Battles) == 0 {
		t.Fatal("expected a battle to be created")
	}
	b := state.Battles[0]
	// Normal duration = 5 + (3+1)/2 = 7s; with SiegeMastery: 7 * 0.6 = 4.2s
	expectedDuration := 7.0 * SiegeMasteryMult
	if b.TimeLeft != expectedDuration {
		t.Errorf("SiegeMastery battle duration: expected %v, got %v", expectedDuration, b.TimeLeft)
	}
}

func TestSiegeMasteryMinDuration(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p1].Tech[TechSiegeMastery] = true
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	// Power diff = 2, both powers very low → duration 5 + (2+1)/2 = 6.5; ×0.6 = 3.9 → clamped to 3.0
	adjHex := Hex{Q: -1, R: 0}
	targetHex := Hex{Q: 0, R: 0}
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 2} // Power 2
	state.Hexes[targetHex] = &HexState{Owner: p2}                                 // Power 0

	a := Action{Type: ActionAttack, Player: p1, Target: targetHex}
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("attack should validate: %v", err)
	}
	ApplyAttack(state, a)
	if len(state.Battles) == 0 {
		t.Fatal("expected a battle to be created")
	}
	b := state.Battles[0]
	// diff=2, standard battle; duration = 5+(2+0)/2=6; ×0.6 → ≥ SiegeMasteryMinDur (3.0) so no clamp
	// Compute expected using runtime multiplication to avoid compile-time constant precision difference
	rawDuration := 5.0 + float64(2+0)/2.0
	expectedDuration := max(SiegeMasteryMinDur, rawDuration*SiegeMasteryMult)
	if b.TimeLeft != expectedDuration {
		t.Errorf("SiegeMastery duration: expected %v, got %v", expectedDuration, b.TimeLeft)
	}
}

func TestSiegeMasteryTieBreak(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p1].Tech[TechSiegeMastery] = true
	state.Players[p2] = &Player{ID: p2, Gold: 0}

	atkHex := Hex{Q: -1, R: 0}
	defHex := Hex{Q: 0, R: 0}
	state.Hexes[atkHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 1} // Power 1
	state.Hexes[defHex] = &HexState{Owner: p2, Building: BuildingPower, Level: 1} // Power 1

	state.Battles = []Battle{{
		AttackerHex: atkHex, DefenderHex: defHex,
		Attacker: p1, Defender: p2, TimeLeft: 0,
	}}

	resolveBattle(state, &state.Battles[0])
	// Tie, attacker has SiegeMastery → attacker wins
	if state.Hexes[defHex].Owner != p1 {
		t.Error("SiegeMastery tie should give win to attacker")
	}
}

func TestFortify(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p1].Tech[TechFortify] = true
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	capHex := Hex{Q: 0, R: 0}
	adjHex := Hex{Q: -1, R: 0}
	state.Hexes[capHex] = &HexState{Owner: p1, Capital: true}
	state.Hexes[adjHex] = &HexState{Owner: p1}

	// Fortify own hex
	a := Action{Type: ActionFortify, Player: p1, Target: capHex}
	if err := ValidateFortify(state, a); err != nil {
		t.Fatalf("Fortify should validate: %v", err)
	}
	ApplyFortify(state, a)
	if state.Players[p1].Gold != 500-FortifyCost {
		t.Errorf("Fortify should cost %vg", FortifyCost)
	}
	if state.Hexes[capHex].FortifyTimer != FortifyDuration {
		t.Errorf("Fortify should set FortifyTimer to %v", FortifyDuration)
	}
}

func TestFortifyBlocksInstantTakeover(t *testing.T) {
	state := NewGameState()
	p1, p2 := PlayerID(1), PlayerID(2)
	state.Players[p1] = &Player{ID: p1, Gold: 500}
	state.Players[p2] = &Player{ID: p2, Gold: 500}

	adjHex := Hex{Q: -1, R: 0}
	defHex := Hex{Q: 0, R: 0}
	// diff = 5 > InstantTakeoverMinDiff(3), but FortifyTimer > 0
	state.Hexes[adjHex] = &HexState{Owner: p1, Building: BuildingPower, Level: 5} // Power 5
	state.Hexes[defHex] = &HexState{Owner: p2, FortifyTimer: FortifyDuration}     // Power 0 + Fortified

	a := Action{Type: ActionAttack, Player: p1, Target: defHex}
	if err := ValidateAttack(state, a); err != nil {
		t.Fatalf("attack should validate: %v", err)
	}
	ApplyAttack(state, a)
	// Should create a battle instead of instant takeover
	if len(state.Battles) == 0 {
		t.Error("Fortify should block instant takeover and force standard battle")
	}
	if state.Hexes[defHex].Owner != p2 {
		t.Error("Fortified hex should not be instantly captured")
	}
}

func TestProsperity(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechProsperity] = true

	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: p1, Building: BuildingGold, Level: 1}

	income := hexIncomeForPlayer(state.Hexes[hex], state.Players[p1])
	// L1 base: (2 + 0.6) × 1.5 = 3.9; + Prosperity 1.0 = 4.9
	expected := (BaseIncomePerSec+GoldPerLevel*1)*GoldBonusMultiplier + ProsperityBonus
	if income != expected {
		t.Errorf("Prosperity L1 income: expected %v, got %v", expected, income)
	}
}

func TestCompoundGrowth(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechCompoundGrowth] = true

	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: p1, Building: BuildingGold, Level: 1}

	income := hexIncomeForPlayer(state.Hexes[hex], state.Players[p1])
	// L1 base: (2 + 0.6) × 1.5 = 3.9; × 1.25 = 4.875
	expected := (BaseIncomePerSec + GoldPerLevel*1) * GoldBonusMultiplier * CompoundGrowthMultiplier
	if income != expected {
		t.Errorf("Compound Growth L1 income: expected %v, got %v", expected, income)
	}
}

func TestCompoundGrowthAndProsperity(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechCompoundGrowth] = true
	state.Players[p1].Tech[TechProsperity] = true

	hex := Hex{Q: 0, R: 0}
	state.Hexes[hex] = &HexState{Owner: p1, Building: BuildingGold, Level: 1}

	income := hexIncomeForPlayer(state.Hexes[hex], state.Players[p1])
	// (2 + 0.6) × 1.5 × 1.25 + 1.0 = 4.875 + 1.0 = 5.875
	expected := (BaseIncomePerSec+GoldPerLevel*1)*GoldBonusMultiplier*CompoundGrowthMultiplier + ProsperityBonus
	if income != expected {
		t.Errorf("Compound Growth + Prosperity L1: expected %v, got %v", expected, income)
	}
}

func TestSupplyLines(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechSupplyLines] = true

	// 30 hexes: tier1=10×0.9=9, tier2=10×1.8=18, tier3=10×2.7=27 → 54 total
	maintenance := calcMaintenanceForPlayer(30, state.Players[p1])
	expected := 10*SupplyLinesTier1 + 10*SupplyLinesTier2 + 10*SupplyLinesTier3
	if maintenance != expected {
		t.Errorf("Supply Lines 30 hexes: expected %v, got %v", expected, maintenance)
	}

	// Without tech: 10×1 + 10×2 + 10×3 = 60
	maintenance2 := calcMaintenanceForPlayer(30, &Player{})
	if maintenance2 != 60.0 {
		t.Errorf("Without Supply Lines 30 hexes: expected 60, got %v", maintenance2)
	}
}

func TestResilience(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 0}
	state.Players[p1].Tech[TechResilience] = true

	// Auto-drop with Resilience: grace period should be ResilienceGracePeriod
	// AutoDropActive transitions: start false, net<0 → active with grace = ResilienceGracePeriod
	// Add a hex so player has income
	capHex := Hex{Q: 0, R: 0}
	extraHex := Hex{Q: 1, R: 0}
	state.Hexes[capHex] = &HexState{Owner: p1, Capital: true}
	state.Hexes[extraHex] = &HexState{Owner: p1}
	// 2 hexes, maintenance = 2×1 = 2/s, income = 4/s → net = +2/s (positive, no drop)

	RunAutoDropPhase(state, TickDt)
	if state.Players[p1].AutoDropActive {
		t.Error("AutoDrop should not be active when net income is positive")
	}
}

func TestResilienceDropRefund(t *testing.T) {
	state := NewGameState()
	p1 := PlayerID(1)
	state.Players[p1] = &Player{ID: p1, Gold: 100}
	state.Players[p1].Tech[TechResilience] = true

	capHex := Hex{Q: 0, R: 0}
	dropHex := Hex{Q: 1, R: 0}
	state.Hexes[capHex] = &HexState{Owner: p1, Capital: true}
	state.Hexes[dropHex] = &HexState{Owner: p1, Building: BuildingGold, Level: 1}

	a := Action{Type: ActionDropHex, Player: p1, Target: dropHex}
	goldBefore := state.Players[p1].Gold
	ApplyDropHex(state, a)
	refundReceived := state.Players[p1].Gold - goldBefore
	expectedRefund := TotalInvested(BuildingGold, 1) * ResilienceDropRefund
	if refundReceived != expectedRefund {
		t.Errorf("Resilience drop refund: expected %v, got %v", expectedRefund, refundReceived)
	}
}
