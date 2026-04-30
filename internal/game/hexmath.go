package game

import "math"

type Hex struct {
	Q int `json:"q"`
	R int `json:"r"`
}

var axialDirections = [6]Hex{
	{1, 0}, {1, -1}, {0, -1},
	{-1, 0}, {-1, 1}, {0, 1},
}

func (h Hex) Neighbors() [6]Hex {
	var out [6]Hex
	for i, d := range axialDirections {
		out[i] = Hex{h.Q + d.Q, h.R + d.R}
	}
	return out
}

func (h Hex) IsNeighbor(other Hex) bool {
	dq := h.Q - other.Q
	dr := h.R - other.R
	ds := (-h.Q - h.R) - (-other.Q - other.R)
	return (abs(dq)+abs(dr)+abs(ds))/2 == 1
}

func (h Hex) Distance(other Hex) int {
	dq := h.Q - other.Q
	dr := h.R - other.R
	ds := (-h.Q - h.R) - (-other.Q - other.R)
	return (abs(dq) + abs(dr) + abs(ds)) / 2
}

const hexSize = 30.0

func HexToPixel(h Hex) (float64, float64) {
	x := hexSize * (math.Sqrt(3)*float64(h.Q) + math.Sqrt(3)/2*float64(h.R))
	y := hexSize * (3.0 / 2 * float64(h.R))
	return x, y
}

func PixelToHex(x, y float64) Hex {
	q := (math.Sqrt(3)/3*x - 1.0/3*y) / hexSize
	r := (2.0 / 3 * y) / hexSize
	return hexRound(q, r)
}

func hexRound(q, r float64) Hex {
	s := -q - r
	rq := math.Round(q)
	rr := math.Round(r)
	rs := math.Round(s)

	qDiff := math.Abs(rq - q)
	rDiff := math.Abs(rr - r)
	sDiff := math.Abs(rs - s)

	if qDiff > rDiff && qDiff > sDiff {
		rq = -rr - rs
	} else if rDiff > sDiff {
		rr = -rq - rs
	}

	return Hex{int(rq), int(rr)}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
