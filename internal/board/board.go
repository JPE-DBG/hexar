// Package board provides the hexagonal board for the Hex game.
//
// The board is an n×n grid of hexagonal cells. Each cell can be empty or owned
// by one of the two players. Player 1 (X) wins by connecting the top edge to
// the bottom edge; Player 2 (O) wins by connecting the left edge to the right
// edge.
package board

import "fmt"

// Cell represents the state of a single cell on the board.
type Cell int

const (
	Empty Cell = iota
	PlayerX
	PlayerO
)

// String returns the display character for a cell.
func (c Cell) String() string {
	switch c {
	case PlayerX:
		return "X"
	case PlayerO:
		return "O"
	default:
		return "."
	}
}

// Board is an n×n hex grid.
type Board struct {
	size  int
	cells [][]Cell
}

// New creates a new empty Board of the given size.
func New(size int) (*Board, error) {
	if size < 2 {
		return nil, fmt.Errorf("board size must be at least 2, got %d", size)
	}
	cells := make([][]Cell, size)
	for i := range cells {
		cells[i] = make([]Cell, size)
	}
	return &Board{size: size, cells: cells}, nil
}

// Size returns the board dimension.
func (b *Board) Size() int { return b.size }

// Get returns the cell value at (row, col). Returns Empty for out-of-bounds coordinates.
func (b *Board) Get(row, col int) Cell {
	if row < 0 || row >= b.size || col < 0 || col >= b.size {
		return Empty
	}
	return b.cells[row][col]
}

// Set places a cell value at (row, col). Returns an error if the position is
// out of bounds or already occupied.
func (b *Board) Set(row, col int, c Cell) error {
	if row < 0 || row >= b.size || col < 0 || col >= b.size {
		return fmt.Errorf("position (%d,%d) is out of bounds", row, col)
	}
	if b.cells[row][col] != Empty {
		return fmt.Errorf("position (%d,%d) is already occupied", row, col)
	}
	b.cells[row][col] = c
	return nil
}

// neighbors returns the up-to-six hex neighbors of (row, col).
func (b *Board) neighbors(row, col int) [][2]int {
	candidates := [][2]int{
		{row - 1, col}, {row - 1, col + 1},
		{row, col - 1}, {row, col + 1},
		{row + 1, col - 1}, {row + 1, col},
	}
	valid := candidates[:0]
	for _, c := range candidates {
		r, cc := c[0], c[1]
		if r >= 0 && r < b.size && cc >= 0 && cc < b.size {
			valid = append(valid, c)
		}
	}
	return valid
}

// HasWon returns true if the given player has a winning connection.
//
//   - PlayerX must connect row 0 to row n-1.
//   - PlayerO must connect col 0 to col n-1.
func (b *Board) HasWon(player Cell) bool {
	visited := make([][]bool, b.size)
	for i := range visited {
		visited[i] = make([]bool, b.size)
	}

	var queue [][2]int

	// Seed the BFS from the player's starting edge.
	for i := 0; i < b.size; i++ {
		var r, c int
		if player == PlayerX {
			r, c = 0, i
		} else {
			r, c = i, 0
		}
		if b.cells[r][c] == player && !visited[r][c] {
			visited[r][c] = true
			queue = append(queue, [2]int{r, c})
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		row, col := cur[0], cur[1]

		// Check if we reached the target edge.
		if player == PlayerX && row == b.size-1 {
			return true
		}
		if player == PlayerO && col == b.size-1 {
			return true
		}

		for _, nb := range b.neighbors(row, col) {
			nr, nc := nb[0], nb[1]
			if !visited[nr][nc] && b.cells[nr][nc] == player {
				visited[nr][nc] = true
				queue = append(queue, [2]int{nr, nc})
			}
		}
	}
	return false
}

// Display returns a human-readable string representation of the board.
func (b *Board) Display() string {
	header := "   "
	for col := 0; col < b.size; col++ {
		header += fmt.Sprintf("%2d", col)
	}
	result := header + "\n"

	for row := 0; row < b.size; row++ {
		indent := ""
		for i := 0; i < row; i++ {
			indent += " "
		}
		line := fmt.Sprintf("%2d %s", row, indent)
		for col := 0; col < b.size; col++ {
			if col > 0 {
				line += " "
			}
			line += b.cells[row][col].String()
		}
		result += line + "\n"
	}
	return result
}
