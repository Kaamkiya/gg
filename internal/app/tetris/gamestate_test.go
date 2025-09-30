package tetris

import (
	"testing"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	"github.com/Kaamkiya/gg/internal/app/tetris/shape"
)

func TestASingleLineIsRemoved(t *testing.T) {
	gamestate := gameState{
		nil,
		nil,
		newGameboard(color.Colors),
		shape.NewRandomizer(),
		0,
		&difficulty{
			20,
			1.0,
			300,
		},
		false,
		pieceDrop{
			dropFinished,
			false,
		},
	}

	for i := range width {
		gamestate.gameBoard.Grid[height-1][i] = color.Blue
	}

	lines := gamestate.checkForCompleteLines(19, 19)
	gamestate.removeCompletedLines(lines)

	if !gamestate.isLineEmpty(19) {
		t.Fatal("Completed single line not removed")
	}

}

func TestMultipleLinesAreRemoved(t *testing.T) {
	gamestate := gameState{
		nil,
		nil,
		newGameboard(color.Colors),
		shape.NewRandomizer(),
		0,
		&difficulty{
			20,
			1.0,
			300,
		},
		false,
		pieceDrop{
			dropFinished,
			false,
		},
	}

	for i := range width {
		gamestate.gameBoard.Grid[height-1][i] = color.Blue
		gamestate.gameBoard.Grid[height-3][i] = color.Blue
		gamestate.gameBoard.Grid[height-4][i] = color.Blue
	}

	gamestate.gameBoard.Grid[height-2][0] = color.Blue
	gamestate.gameBoard.Grid[height-5][0] = color.Blue

	lines := gamestate.checkForCompleteLines(16, 19)
	gamestate.removeCompletedLines(lines)

	if gamestate.gameBoard.Grid[height-1][0] != color.Blue && gamestate.gameBoard.Grid[height-1][1] != color.None {
		t.Fatal("Second to last line didn't drop when last line was completed")
	}

	if gamestate.gameBoard.Grid[height-2][0] != color.Blue && gamestate.gameBoard.Grid[height-2][1] != color.None {
		t.Fatal("Fifth to last line didn't drop when third to last line was completed")
	}

	if !gamestate.isLineEmpty(height - 3) {
		t.Fatal("Lines didn't move correctly when lines where completed")
	}

}
