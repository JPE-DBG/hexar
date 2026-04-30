package game

import (
	"testing"

	"github.com/JPE-DBG/hexar/internal/board"
)

func TestNewGame(t *testing.T) {
	g, err := New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.CurrentTurn() != board.PlayerX {
		t.Error("expected first turn to be PlayerX")
	}
	if g.State() != InProgress {
		t.Errorf("expected InProgress, got %v", g.State())
	}
	if g.IsOver() {
		t.Error("expected game to not be over at start")
	}
}

func TestNewGameInvalidSize(t *testing.T) {
	_, err := New(1)
	if err == nil {
		t.Error("expected error for invalid size, got nil")
	}
}

func TestPlaySwitchesTurns(t *testing.T) {
	g, _ := New(5)
	if err := g.Play(0, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.CurrentTurn() != board.PlayerO {
		t.Error("expected turn to switch to PlayerO after PlayerX plays")
	}
}

func TestPlayOccupiedCell(t *testing.T) {
	g, _ := New(5)
	_ = g.Play(0, 0)
	if err := g.Play(0, 0); err == nil {
		t.Error("expected error when playing on occupied cell")
	}
}

func TestPlayOutOfBounds(t *testing.T) {
	g, _ := New(3)
	if err := g.Play(10, 10); err == nil {
		t.Error("expected error for out-of-bounds move")
	}
}

func TestPlayAfterGameOver(t *testing.T) {
	g, _ := New(3)
	// X wins by connecting top to bottom via column 0
	_ = g.Play(0, 0) // X
	_ = g.Play(0, 1) // O
	_ = g.Play(1, 0) // X
	_ = g.Play(0, 2) // O
	_ = g.Play(2, 0) // X — X wins
	if !g.IsOver() {
		t.Fatal("expected game to be over")
	}
	if err := g.Play(2, 2); err == nil {
		t.Error("expected error when playing after game over")
	}
}

func TestPlayerXWins(t *testing.T) {
	g, _ := New(3)
	// X connects row 0 → row 2 via col 2
	_ = g.Play(0, 2) // X
	_ = g.Play(1, 0) // O
	_ = g.Play(1, 2) // X
	_ = g.Play(2, 1) // O
	_ = g.Play(2, 2) // X — wins
	if g.State() != WonByX {
		t.Errorf("expected WonByX, got %v", g.State())
	}
	if g.State().String() != "Player X wins!" {
		t.Errorf("unexpected state string: %q", g.State().String())
	}
}

func TestPlayerOWins(t *testing.T) {
	g, _ := New(3)
	// O connects col 0 → col 2 via row 2
	_ = g.Play(0, 0) // X
	_ = g.Play(2, 0) // O
	_ = g.Play(0, 1) // X
	_ = g.Play(2, 1) // O
	_ = g.Play(1, 2) // X
	_ = g.Play(2, 2) // O — wins
	if g.State() != WonByO {
		t.Errorf("expected WonByO, got %v", g.State())
	}
	if g.State().String() != "Player O wins!" {
		t.Errorf("unexpected state string: %q", g.State().String())
	}
}

func TestInProgressString(t *testing.T) {
	g, _ := New(3)
	if g.State().String() != "In progress" {
		t.Errorf("unexpected string: %q", g.State().String())
	}
}

func TestBoard(t *testing.T) {
	g, _ := New(4)
	if g.Board() == nil {
		t.Error("expected non-nil board")
	}
}
