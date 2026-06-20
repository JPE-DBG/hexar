package game

func RunUnits(state *GameState, dt float64) {
	// Advance move timers and move units.
	for _, u := range state.Units {
		u.MoveTimer += dt
		if u.MoveTimer >= SoldierMovePeriod {
			u.MoveTimer = 0
			advanceUnit(state, u)
		}
	}

	// Resolve combat between units on the same hex.
	resolveUnitCombat(state, dt)

	// Check if any unit has reached the enemy capital.
	checkCapitalContact(state)

	// Remove dead units.
	removeDeadUnits(state)
}

func spawnUnit(state *GameState, owner PlayerID, ut UnitType) {
	capitalHex := findCapitalHex(state, owner)
	if capitalHex == nil {
		return
	}
	state.Units = append(state.Units, &Unit{
		ID:    state.NextUnitID,
		Owner: owner,
		Type:  ut,
		Pos:   *capitalHex,
		HP:    SoldierHP,
	})
	state.NextUnitID++
}

func findCapitalHex(state *GameState, owner PlayerID) *Hex {
	for h, hs := range state.Hexes {
		if hs.Capital && hs.Owner == owner {
			hCopy := h
			return &hCopy
		}
	}
	return nil
}

func advanceUnit(state *GameState, u *Unit) {
	enemyCapital := findCapitalHex(state, opponentOf(u.Owner))
	if enemyCapital == nil {
		return
	}
	// Don't move if in combat (enemy unit on same hex).
	if hasEnemyOnSameHex(state, u) {
		return
	}
	next := nextHexToward(state, u.Pos, *enemyCapital)
	if next != u.Pos {
		u.Pos = next
	}
}

func nextHexToward(state *GameState, from, target Hex) Hex {
	best := from
	bestDist := from.Distance(target)
	for _, n := range from.Neighbors() {
		if _, ok := state.Hexes[n]; !ok {
			continue
		}
		d := n.Distance(target)
		if d < bestDist {
			bestDist = d
			best = n
		}
	}
	return best
}

func hasEnemyOnSameHex(state *GameState, u *Unit) bool {
	for _, other := range state.Units {
		if other.ID != u.ID && other.Owner != u.Owner && other.Pos == u.Pos {
			return true
		}
	}
	return false
}

func resolveUnitCombat(state *GameState, dt float64) {
	// Advance attack timers for units in combat.
	for _, u := range state.Units {
		if u.HP <= 0 {
			continue
		}
		if !hasEnemyOnSameHex(state, u) {
			continue
		}
		u.AttackTimer += dt
		if u.AttackTimer >= SoldierAttackPeriod {
			u.AttackTimer = 0
			// Deal damage to one enemy on the same hex.
			for _, enemy := range state.Units {
				if enemy.Owner != u.Owner && enemy.Pos == u.Pos && enemy.HP > 0 {
					enemy.HP -= SoldierAttackPower
					break
				}
			}
		}
	}
}

func checkCapitalContact(state *GameState) {
	for _, u := range state.Units {
		if u.HP <= 0 {
			continue
		}
		enemyCapital := findCapitalHex(state, opponentOf(u.Owner))
		if enemyCapital == nil {
			continue
		}
		if u.Pos == *enemyCapital {
			// Check if a defending unit is present; attack it instead of the capital.
			hasDefender := false
			for _, other := range state.Units {
				if other.Owner != u.Owner && other.Pos == u.Pos && other.HP > 0 {
					hasDefender = true
					break
				}
			}
			if !hasDefender {
				state.Players[opponentOf(u.Owner)].CapitalHP -= CapitalDamagePerUnit
				u.HP = 0
			}
		}
	}
}

func removeDeadUnits(state *GameState) {
	alive := state.Units[:0]
	for _, u := range state.Units {
		if u.HP > 0 {
			alive = append(alive, u)
		}
	}
	state.Units = alive
}
