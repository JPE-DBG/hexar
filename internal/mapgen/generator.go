package mapgen

import "hexar/internal/game"

const mapRadius = 4

func GenerateTestMap() map[game.Hex]*game.HexState {
	hexes := make(map[game.Hex]*game.HexState)

	for q := -mapRadius; q <= mapRadius; q++ {
		r1 := max(-mapRadius, -q-mapRadius)
		r2 := min(mapRadius, -q+mapRadius)
		for r := r1; r <= r2; r++ {
			hexes[game.Hex{Q: q, R: r}] = &game.HexState{}
		}
	}

	return hexes
}

func SpawnPositions() [2]game.Hex {
	return [2]game.Hex{
		{Q: -mapRadius, R: 0},
		{Q: mapRadius, R: 0},
	}
}
