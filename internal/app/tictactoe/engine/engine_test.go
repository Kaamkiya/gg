package engine

import (
	"testing"
)

var testCases = []struct {
	input    []int
	expected int
}{
	// #0: first row
	{
		input:    []int{1, 1, 0, -1, 0, -1, 0, 0, 0},
		expected: 2,
	},
	// #1: first col
	{
		input:    []int{1, 0, 0, 1, -1, 0, 0, -1, 0},
		expected: 6,
	},
	// #2: second col
	{
		input:    []int{0, 1, 0, 0, 1, -1, 0, 0, -1},
		expected: 7,
	},
	// #3: diagonal left (\)
	{
		input:    []int{1, -1, 0, 0, 1, -1, 0, 0, 0},
		expected: 8,
	},
	// #4: diagonal right (/)
	{
		input:    []int{0, -1, 1, 0, 1, -1, 0, 0, 0},
		expected: 6,
	},
	// #5: middle row
	{
		input:    []int{0, 0, 0, 1, 0, 1, -1, -1, 0},
		expected: 4,
	},
	// #6: last row
	{
		input:    []int{0, 0, 0, -1, -1, 0, 1, 1, 0},
		expected: 8,
	},
	// #7: last col
	{
		input:    []int{0, 0, 1, -1, 0, 1, 0, 0, 0},
		expected: 8,
	},
	// #8: No move
	{
		input:    []int{1, -1, 1, -1, -1, 1, 1, 1, -1},
		expected: -1, // Indicates no move left to win
	},
}

func TestEngine_Solve(t *testing.T) {
	BOARD_SIZE := 3
	engine := NewEngine(DEPTH)

	for _, tc := range testCases {
		t.Run("Testing solve", func(t *testing.T) {
			board := NewBoard(BOARD_SIZE)
			board.Load(tc.input)

			move := engine.ai.Solve(board)

			if move != tc.expected {
				t.Errorf("expected move %d, got %d", tc.expected, move)
			}
		})
	}
}

func TestEngine_CheckWin(t *testing.T) {
	BOARD_SIZE := 3
	board := NewBoard(BOARD_SIZE)
	engine := NewEngine(DEPTH)

	t.Run("Empty board", func(t *testing.T) {
		if engine.CheckWin(board, 0) {
			t.Error("expected no win")
		}
	})

	t.Run("Horizontal win", func(t *testing.T) {
		board.SetCell(0, P1)
		board.SetCell(1, P1)
		board.SetCell(2, P1)
		if !engine.CheckWin(board, 2) {
			t.Error("expected win")
		}
	})

	t.Run("Vertical win", func(t *testing.T) {
		board = NewBoard(BOARD_SIZE)
		board.SetCell(0, P1)
		board.SetCell(3, P1)
		board.SetCell(6, P1)
		if !engine.CheckWin(board, 6) {
			t.Error("expected win")
		}
	})

	t.Run("Left diagonal win", func(t *testing.T) {
		board = NewBoard(BOARD_SIZE)
		board.SetCell(0, P1)
		board.SetCell(4, P1)
		board.SetCell(8, P1)
		if !engine.CheckWin(board, 8) {
			t.Error("expected win")
		}
	})

	t.Run("Right diagonal win", func(t *testing.T) {
		board = NewBoard(BOARD_SIZE)
		board.SetCell(2, P1)
		board.SetCell(4, P1)
		board.SetCell(6, P1)
		if !engine.CheckWin(board, 6) {
			t.Error("expected win")
		}
	})
}

func TestEngine_GetLegalMoves(t *testing.T) {
	BOARD_SIZE := 4
	board := NewBoard(BOARD_SIZE)
	engine := NewEngine(DEPTH)
	moves := []int{}

	t.Run("Empty board", func(t *testing.T) {
		moves = engine.GetLegalMoves(board)
		if len(moves) != BOARD_SIZE*BOARD_SIZE {
			t.Errorf("expected %d moves, got %d", BOARD_SIZE*BOARD_SIZE, len(moves))
		}
	})

	t.Run("Full board", func(t *testing.T) {
		for _, move := range moves {
			board.SetCell(move, P1)
		}
		moves := engine.GetLegalMoves(board)
		if len(moves) != 0 {
			t.Errorf("expected 0 moves, got %d", len(moves))
		}
	})

	t.Run("One empty cell", func(t *testing.T) {
		board.SetCell(0, EMPTY)
		moves = engine.GetLegalMoves(board)
		if len(moves) != 1 {
			t.Errorf("expected 1 move, got %d", len(moves))
		}
	})
}
