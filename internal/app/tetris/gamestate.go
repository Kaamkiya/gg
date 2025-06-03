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

func (gs *GameState) HandleTick() {
	middleX := (gameboard.Width / 2) - 1
	if gs.nextShape == nil {
		newShape := shape.CreateNew(middleX, 0)
		gs.nextShape = &newShape
	}

	if gs.currentShape == nil {
		newShape := shape.CreateNew(middleX, 0)
		gs.currentShape = gs.nextShape
		gs.nextShape = &newShape
		gs.addShape(gs.currentShape)
		return
	}

	if !gs.applyTransformation(gs.currentShape.MoveDown) {
		_, posY := gs.currentShape.GetPosition()
		gs.handleCompletedLines(posY, posY+gs.currentShape.GetHeight()-1)
		gs.currentShape = nil
	}
}

func (gs *GameState) HandleLeft() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.MoveLeft)
}

func (gs *GameState) HandleRight() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.MoveRight)
}

func (gs *GameState) HandleDown() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.MoveDown)
}

func (gs *GameState) HandleLeftRotate() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.RotateLeft)
}

func (gs *GameState) HandleRightRotate() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.RotateRight)
}

func (gs *GameState) applyTransformation(tranformation func() shape.Shape) bool {
	newShape := tranformation()

	gs.deleteShape(gs.currentShape)

	if gs.isShapeValid(newShape) {
		gs.currentShape = &newShape
		gs.addShape(gs.currentShape)

		return true
	} else {
		gs.addShape(gs.currentShape)
	}

	return false
}

func (gs *GameState) isShapeValid(shape shape.Shape) bool {
	shapeGrid := shape.GetGrid()
	posX, posY := shape.GetPosition()

	if posX < 0 {
		return false
	}

	if posX+len(shapeGrid[0]) > gameboard.Width || posY+len(shapeGrid) > gameboard.Height {
		return false
	}

	for i := range shapeGrid {
		for j := range shapeGrid[i] {
			if shapeGrid[i][j] {
				if gs.gameBoard.Grid[posY+i][posX+j] != color.Black {
					return false
				}
			}
		}
	}

	return true
}

func (gs *GameState) addShape(shape *shape.Shape) {
	gs.modidfyColorGridFromShape(shape, shape.GetColor())
}

func (gs *GameState) deleteShape(shape *shape.Shape) {
	gs.modidfyColorGridFromShape(shape, color.Black)
}

func (gs *GameState) modidfyColorGridFromShape(shape *shape.Shape, color color.Color) {
	shapeGrid := shape.GetGrid()
	posX, posY := shape.GetPosition()

	for i := range shapeGrid {
		for j := range shapeGrid[i] {
			if shapeGrid[i][j] {
				gs.gameBoard.Grid[posY+i][posX+j] = color
			}
		}
	}
}

func (gs *GameState) handleCompletedLines(from, to int) {
	completedLines := gs.checkForCompleteLines(from, to)

	if len(completedLines) == 0 {
		return
	}

	distanceToCopyFrom := 1
	nextCompletedLine := 1

	for i := completedLines[0]; i >= 0; i-- {
		if i-distanceToCopyFrom < 0 {
			return
		}

		if gs.isLineEmpty(i) {
			return
		}

		for nextCompletedLine < len(completedLines) && completedLines[nextCompletedLine] == i-distanceToCopyFrom {
			nextCompletedLine++
			distanceToCopyFrom++
		}

		for j := range gameboard.Width {
			gs.gameBoard.Grid[i][j] = gs.gameBoard.Grid[i-distanceToCopyFrom][j]
		}
	}

}

func (gs *GameState) checkForCompleteLines(from, to int) []int {
	completedLines := make([]int, 0, 4)
	for i := to; i >= from; i-- {
		if gs.isLineCompleted(i) {
			completedLines = append(completedLines, i)
		}
	}

	return completedLines
}

func (gs *GameState) isLineCompleted(line int) bool {
	for i := range gameboard.Width {
		if gs.gameBoard.Grid[line][i] == color.Black {
			return false
		}
	}

	return true
}

func (gs *GameState) isLineEmpty(line int) bool {
	for i := range gameboard.Width {
		if gs.gameBoard.Grid[line][i] != color.Black {
			return false
		}
	}

	return true
}
