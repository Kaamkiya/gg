package tron

import (
	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
)

type gameState struct {
	gameboard gameboard.Gameboard
	players   []player
}

func newGameState() gameState {
	return gameState{
		gameboard: gameboard.NewGameboard(gameboard.Red, gameboard.Blue),
		players: []player{
			{gameboard.Right, gameboard.Red},
			{gameboard.Left, gameboard.Blue},
		},
	}
}

func (gs *gameState) movePlayers() {
	for _, player := range gs.players {
		gs.gameboard.MovePlayer(player.Color, player.Direction, 1)
	}
}
