package tron

import (
	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
)

type gameState struct {
	gameboard gameboard.Gameboard
}

func newGameState() gameState {
	return gameState{
		gameboard: gameboard.NewGameboard(gameboard.Red, gameboard.Blue),
	}
}
