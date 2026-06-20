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

	// Always send full player state (bar and deck change every tick).
	for _, p := range curr.Players {
		delta.Players = append(delta.Players, p)
	}

	// Always send full unit list (positions change every tick; empty slice clears stale units).
	delta.Units = curr.Units

	for key, hex := range curr.Hexes {
		prevHex, ok := prev.Hexes[key]
		if !ok || hexChanged(prevHex, hex) {
			delta.HexChanges = append(delta.HexChanges, hex)
		}
	}

	return delta
}

// hexChanged checks all mutable HexDTO fields. Coordinates (Q, R) are immutable.
// If new fields are added to HexDTO, update this function or deltas will desync.
func hexChanged(prev, curr *HexDTO) bool {
	return prev.Owner != curr.Owner ||
		prev.Capital != curr.Capital ||
		prev.HasBuilding != curr.HasBuilding
}
