package tetris

import (
	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	"github.com/Kaamkiya/gg/internal/app/tetris/gameboard"
	"github.com/Kaamkiya/gg/internal/app/tetris/shape"
)

type GameState struct {
	nextShape    *shape.Shape
	currentShape *shape.Shape
	gameBoard    *gameboard.Gameboard
}

func (gameState *GameState) HandleTick() {
	middleX := gameboard.Width / 2
	if gameState.nextShape == nil {
		newShape := shape.CreateNew(middleX, 0)
		gameState.nextShape = &newShape
	}

	if gameState.currentShape == nil {
		newShape := shape.CreateNew(middleX, 0)
		gameState.currentShape = gameState.nextShape
		gameState.nextShape = &newShape
		gameState.addShape(gameState.currentShape)
		return
	}
}

func (gamestate *GameState) addShape(shape *shape.Shape) {
	gamestate.modidfyColorGridFromShape(shape, shape.GetColor())
}

func (gamestate *GameState) deleteShape(shape *shape.Shape) {
	gamestate.modidfyColorGridFromShape(shape, color.Black)
}

func (gamestate *GameState) modidfyColorGridFromShape(shape *shape.Shape, color color.Color) {
	shapeGrid := shape.GetGrid()
	posX, posY := shape.GetPosition()

	for i := range shapeGrid {
		for j := range shapeGrid[i] {
			if shapeGrid[i][j] {
				gamestate.gameBoard.Grid[posY+i][posX+j] = color
			}
		}
	}
}
