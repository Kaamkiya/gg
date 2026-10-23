package tron

import (
	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
)

const (
	// The number of nitros each player begins with.
	initialNitros = 3

	// Each nitro usage increases the player speed for that many rounds.
	roundsPerNitroUsage = 15

	// The factor which the speed increases with each nitro usage.
	nitroBoost = 2

	startCountDown = 3
)

type gameState struct {
	gameboard      *gameboard.Board
	players        []player
	gameEnded      bool
	winner         int
	startCountDown int
	gameConfig     gameConfig
}

func newGameState(gameConfig gameConfig) *gameState {
	players := make([]player, len(gameConfig.playerConfigs))
	playerColors := make([]gameboard.Color, len(gameConfig.playerConfigs))

	players[0] = newPlayer(gameConfig.playerConfigs[0].color, gameboard.Right)
	playerColors[0] = gameConfig.playerConfigs[0].color
	players[1] = newPlayer(gameConfig.playerConfigs[1].color, gameboard.Left)
	playerColors[1] = gameConfig.playerConfigs[1].color

	if len(gameConfig.playerConfigs) == 3 {
		players[2] = newPlayer(gameConfig.playerConfigs[2].color, gameboard.Up)
		playerColors[2] = gameConfig.playerConfigs[2].color
	}

	return &gameState{
		gameboard:      gameboard.NewBoard(playerColors...),
		players:        players,
		startCountDown: startCountDown,
		gameConfig:     gameConfig,
	}
}

func (gs *gameState) updatePlayers() {
	for i, player := range gs.players {
		gs.gameboard.AnimatePlayer(player.Color)

		if gs.gameEnded || player.Crashed {
			continue
		}

		if gs.gameboard.MovePlayer(player.Color, player.Direction, gs.players[i].speedForNextTick()) {
			gs.players[i].Crashed = true
		}
	}
}

func (gs *gameState) checkForGameEnd() (gameEnded bool, winner int) {
	var aliveIndex, aliveCount int
	// See if at most one player is alive.
	for i, player := range gs.players {
		if !player.Crashed {
			aliveCount++
			aliveIndex = i
		}
		if aliveCount > 1 {
			return false, -1
		}
	}

	// In case two or more players crash on the same round it's a draw.
	if aliveCount == 0 {
		return true, -1
	}

	if aliveCount == 1 {
		return true, aliveIndex
	}

	return false, -1
}
