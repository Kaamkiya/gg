package tron

import (
	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
)

type gameState struct {
	gameboard      gameboard.Gameboard
	players        []*player
	gameEnded      bool
	winner         int
	startCountDown int
}

func newGameState() gameState {
	return gameState{
		gameboard: gameboard.NewGameboard(gameboard.Red, gameboard.Blue),
		players: []*player{
			{Direction: gameboard.Right, Color: gameboard.Red},
			{Direction: gameboard.Left, Color: gameboard.Blue},
		},
		startCountDown: 3,
	}
}

func (gs *gameState) updatePlayers() {
	for _, player := range gs.players {
		if !gs.gameboard.UpdatePlayer(player.Color, player.Direction, 1, gs.gameEnded) {
			player.Crashed = true
		}
	}
}

func (gs *gameState) checkForWinner() (gameEnded bool, winner int) {
	var aliveIndex, aliveCount int
	for i, player := range gs.players {
		if !player.Crashed {
			aliveCount = aliveCount + 1
			aliveIndex = i
		}
		if aliveCount > 1 {
			return false, -1
		}
	}

	if aliveCount == 0 {
		return true, -1
	}

	if aliveCount == 1 {
		return true, aliveIndex
	}

	return false, -1
}
