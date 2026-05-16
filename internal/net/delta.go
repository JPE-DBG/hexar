package net

func buildDelta(prev, curr *SnapshotMsg) *DeltaMsg {
	delta := &DeltaMsg{
		Type:          MsgDelta,
		Elapsed:       curr.Elapsed,
		Over:          curr.Over,
		Winner:        curr.Winner,
		WinReason:     curr.WinReason,
		Waiting:       curr.Waiting,
		Paused:        curr.Paused,
		PauseTimeLeft: curr.PauseTimeLeft,
	}

	for key, p := range curr.Players {
		prevP, ok := prev.Players[key]
		if ok && techEqual(prevP.Tech, p.Tech) {
			stripped := *p
			stripped.Tech = nil
			delta.Players = append(delta.Players, &stripped)
		} else {
			delta.Players = append(delta.Players, p)
		}
	}

	delta.Battles = curr.Battles

	for key, hex := range curr.Hexes {
		prevHex, ok := prev.Hexes[key]
		if !ok || hexChanged(prevHex, hex) {
			delta.HexChanges = append(delta.HexChanges, hex)
		}
	}

	return delta
}

func techEqual(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// hexChanged checks all mutable HexDTO fields. Coordinates (Q, R) are immutable.
// If new fields are added to HexDTO, update this function or deltas will desync.
func hexChanged(prev, curr *HexDTO) bool {
	return prev.Owner != curr.Owner ||
		prev.Building != curr.Building ||
		prev.Level != curr.Level ||
		prev.Capital != curr.Capital ||
		fortifyChanged(prev.FortifyTimer, curr.FortifyTimer) ||
		prev.PreviousOwner != curr.PreviousOwner
}

func fortifyChanged(prev, curr float64) bool {
	if (prev > 0) != (curr > 0) {
		return true
	}
	return int(prev) != int(curr)
}
