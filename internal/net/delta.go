package net

func buildDelta(prev, curr *SnapshotMsg) *DeltaMsg {
	delta := &DeltaMsg{
		Type:    MsgDelta,
		Elapsed: curr.Elapsed,
		Over:    curr.Over,
		Winner:  curr.Winner,
	}

	for _, p := range curr.Players {
		delta.Players = append(delta.Players, p)
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

// hexChanged checks all mutable HexDTO fields. Coordinates (Q, R) are immutable.
// If new fields are added to HexDTO, update this function or deltas will desync.
func hexChanged(prev, curr *HexDTO) bool {
	return prev.Owner != curr.Owner ||
		prev.Building != curr.Building ||
		prev.Level != curr.Level ||
		prev.Capital != curr.Capital ||
		prev.FortifyTimer != curr.FortifyTimer ||
		prev.PreviousOwner != curr.PreviousOwner
}
