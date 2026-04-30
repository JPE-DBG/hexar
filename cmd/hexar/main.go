// Command hexar is a command-line implementation of the Hex board game.
//
// Two players (X and O) take turns placing stones on an n×n hexagonal grid.
//   - Player X tries to connect the top edge to the bottom edge.
//   - Player O tries to connect the left edge to the right edge.
//
// Usage:
//
//	hexar [size]
//
// The optional size argument sets the board dimension (default: 7).
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/JPE-DBG/hexar/internal/game"
)

const defaultSize = 7

func main() {
	size := defaultSize
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil || n < 2 {
			fmt.Fprintf(os.Stderr, "Invalid board size %q — must be an integer >= 2\n", os.Args[1])
			os.Exit(1)
		}
		size = n
	}

	g, err := game.New(size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create game: %v\n", err)
		os.Exit(1)
	}

	printRules(size)

	scanner := bufio.NewScanner(os.Stdin)
	for !g.IsOver() {
		fmt.Print(g.Board().Display())
		fmt.Printf("Turn: Player %s — enter row col: ", g.CurrentTurn())

		if !scanner.Scan() {
			break
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) != 2 {
			fmt.Println("Please enter two numbers: row col")
			continue
		}
		row, errR := strconv.Atoi(parts[0])
		col, errC := strconv.Atoi(parts[1])
		if errR != nil || errC != nil {
			fmt.Println("Invalid input — please enter two integers.")
			continue
		}

		if err := g.Play(row, col); err != nil {
			fmt.Printf("Invalid move: %v\n", err)
			continue
		}
	}

	fmt.Print(g.Board().Display())
	fmt.Println(g.State())
}

func printRules(size int) {
	fmt.Printf("=== HEXAR — %d×%d board ===\n", size, size)
	fmt.Println("Player X: connect TOP ↕ to BOTTOM")
	fmt.Println("Player O: connect LEFT ↔ to RIGHT")
	fmt.Println("Enter moves as: <row> <col>  (zero-indexed)")
	fmt.Println()
}
