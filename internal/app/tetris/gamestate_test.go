package tetris

import (
	"testing"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
)

func TestASingleLineIsRemoved(t *testing.T) {
	gamestate := GameState{
		nil,
		nil,
		NewGameboard(color.Colors),
	}

	for i := range Width {
		gamestate.gameBoard.Grid[Height-1][i] = color.Blue
	}

	gamestate.handleCompletedLines(19, 19)

	if !gamestate.isLineEmpty(19) {
		t.Fatal("Completed single line not removed")
	}

}

func TestMultipleLinesAreRemoved(t *testing.T) {
	gamestate := GameState{
		nil,
		nil,
		NewGameboard(color.Colors),
	}

	for i := range Width {
		gamestate.gameBoard.Grid[Height-1][i] = color.Blue
		gamestate.gameBoard.Grid[Height-3][i] = color.Blue
		gamestate.gameBoard.Grid[Height-4][i] = color.Blue
	}

	gamestate.gameBoard.Grid[Height-2][0] = color.Blue
	gamestate.gameBoard.Grid[Height-5][0] = color.Blue

	gamestate.handleCompletedLines(16, 19)

	if gamestate.gameBoard.Grid[Height-1][0] != color.Blue && gamestate.gameBoard.Grid[Height-1][1] != color.Black {
		t.Fatal("Second to last line didn't drop when last line was completed")
	}

	if gamestate.gameBoard.Grid[Height-2][0] != color.Blue && gamestate.gameBoard.Grid[Height-2][1] != color.Black {
		t.Fatal("Fifth to last line didn't drop when third to last line was completed")
	}

	if !gamestate.isLineEmpty(Height - 3) {
		t.Fatal("Lines didn't move correctly when lines where completed")
	}

}
