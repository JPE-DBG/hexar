package mapgen

import "hexar/internal/game"

func GenerateTestMap() map[game.Hex]*game.HexState {
	hexes := make(map[game.Hex]*game.HexState)

	radius := 4
	for q := -radius; q <= radius; q++ {
		r1 := max(-radius, -q-radius)
		r2 := min(radius, -q+radius)
		for r := r1; r <= r2; r++ {
			hexes[game.Hex{Q: q, R: r}] = &game.HexState{}
		}
	}

	return hexes
}

func SpawnPositions() [2]game.Hex {
	return [2]game.Hex{
		{Q: -4, R: 0},
		{Q: 4, R: 0},
	}
}
