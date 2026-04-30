// Package game manages the state and rules of a Hex game session.
package game

import (
	"errors"

	"github.com/JPE-DBG/hexar/internal/board"
)

// State represents the current phase of a game.
type State int

const (
	InProgress State = iota
	WonByX
	WonByO
)

// String returns a human-readable description of the game state.
func (s State) String() string {
	switch s {
	case WonByX:
		return "Player X wins!"
	case WonByO:
		return "Player O wins!"
	default:
		return "In progress"
	}
}

// Game holds all information for a single Hex game.
type Game struct {
	board       *board.Board
	currentTurn board.Cell
	state       State
}

// New creates a new Game on a board of the given size.
// The first move always belongs to PlayerX.
func New(size int) (*Game, error) {
	b, err := board.New(size)
	if err != nil {
		return nil, err
	}
	return &Game{
		board:       b,
		currentTurn: board.PlayerX,
		state:       InProgress,
	}, nil
}

// Board returns the underlying board (read-only use intended).
func (g *Game) Board() *board.Board { return g.board }

// CurrentTurn returns whose turn it is.
func (g *Game) CurrentTurn() board.Cell { return g.currentTurn }

// State returns the current game state.
func (g *Game) State() State { return g.state }

// IsOver returns true when someone has won.
func (g *Game) IsOver() bool { return g.state != InProgress }

// Play places a stone for the current player at (row, col).
// Returns an error if the game is already over or the move is invalid.
func (g *Game) Play(row, col int) error {
	if g.IsOver() {
		return errors.New("game is already over")
	}

	if err := g.board.Set(row, col, g.currentTurn); err != nil {
		return err
	}

	if g.board.HasWon(g.currentTurn) {
		if g.currentTurn == board.PlayerX {
			g.state = WonByX
		} else {
			g.state = WonByO
		}
		return nil
	}

	// Switch turns.
	if g.currentTurn == board.PlayerX {
		g.currentTurn = board.PlayerO
	} else {
		g.currentTurn = board.PlayerX
	}
	return nil
}
