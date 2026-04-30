package board

import (
	"strings"
	"testing"
)

func TestNewBoard(t *testing.T) {
	b, err := New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Size() != 5 {
		t.Errorf("expected size 5, got %d", b.Size())
	}
}

func TestNewBoardTooSmall(t *testing.T) {
	_, err := New(1)
	if err == nil {
		t.Error("expected error for size < 2, got nil")
	}
}

func TestSetAndGet(t *testing.T) {
	b, _ := New(3)
	if err := b.Set(1, 1, PlayerX); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := b.Get(1, 1); got != PlayerX {
		t.Errorf("expected PlayerX, got %v", got)
	}
}

func TestSetAlreadyOccupied(t *testing.T) {
	b, _ := New(3)
	_ = b.Set(0, 0, PlayerX)
	if err := b.Set(0, 0, PlayerO); err == nil {
		t.Error("expected error for occupied cell, got nil")
	}
}

func TestSetOutOfBounds(t *testing.T) {
	b, _ := New(3)
	if err := b.Set(5, 5, PlayerX); err == nil {
		t.Error("expected error for out-of-bounds, got nil")
	}
}

func TestGetOutOfBounds(t *testing.T) {
	b, _ := New(3)
	if got := b.Get(-1, 0); got != Empty {
		t.Errorf("expected Empty for out-of-bounds, got %v", got)
	}
}

func TestHasWonPlayerX(t *testing.T) {
	b, _ := New(3)
	// X connects top (row 0) to bottom (row 2) via column 1
	_ = b.Set(0, 1, PlayerX)
	_ = b.Set(1, 1, PlayerX)
	_ = b.Set(2, 1, PlayerX)
	if !b.HasWon(PlayerX) {
		t.Error("expected PlayerX to have won")
	}
	if b.HasWon(PlayerO) {
		t.Error("expected PlayerO to not have won")
	}
}

func TestHasWonPlayerO(t *testing.T) {
	b, _ := New(3)
	// O connects left (col 0) to right (col 2) via row 1
	_ = b.Set(1, 0, PlayerO)
	_ = b.Set(1, 1, PlayerO)
	_ = b.Set(1, 2, PlayerO)
	if !b.HasWon(PlayerO) {
		t.Error("expected PlayerO to have won")
	}
	if b.HasWon(PlayerX) {
		t.Error("expected PlayerX to not have won")
	}
}

func TestHasWonNoWinner(t *testing.T) {
	b, _ := New(3)
	if b.HasWon(PlayerX) || b.HasWon(PlayerO) {
		t.Error("expected no winner on empty board")
	}
}

func TestDisplay(t *testing.T) {
	b, _ := New(3)
	_ = b.Set(0, 0, PlayerX)
	out := b.Display()
	if !strings.Contains(out, "X") {
		t.Error("expected display to contain 'X'")
	}
}

func TestCellString(t *testing.T) {
	cases := []struct {
		cell Cell
		want string
	}{
		{Empty, "."},
		{PlayerX, "X"},
		{PlayerO, "O"},
	}
	for _, tc := range cases {
		if got := tc.cell.String(); got != tc.want {
			t.Errorf("Cell(%d).String() = %q, want %q", tc.cell, got, tc.want)
		}
	}
}
