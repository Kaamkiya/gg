package tetris

import (
	"maps"
	"slices"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	"github.com/Kaamkiya/gg/internal/app/tetris/shape"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const Height = 20
const Width = 10
const LineAnimationInterval time.Duration = 100 * time.Millisecond

type Gameboard struct {
	Colors map[color.Color]lipgloss.Style
	Grid   [Height][Width]color.Color
}

type LineAnimationTick struct {
	linesToUpdate      map[int][Width]color.Color
	animationCountDown int
}

func NewGameboard(colors map[color.Color]lipgloss.Style) *Gameboard {
	grid := [Height][Width]color.Color{}

	return &Gameboard{colors, grid}
}

type GameState struct {
	nextShape    *shape.Shape
	currentShape *shape.Shape
	gameBoard    *Gameboard
	isAnimating  bool
}

func (gs *GameState) HandleGameProgressTick() tea.Cmd {
	middleX := (Width / 2) - 1
	if gs.nextShape == nil {
		newShape := shape.CreateNew(middleX, 0)
		gs.nextShape = &newShape
	}

	nextCmd := tea.Tick(gameProgressTickDelay, func(t time.Time) tea.Msg {
		return GameProgressTick{}
	})

	if gs.currentShape == nil {
		newShape := shape.CreateNew(middleX, 0)
		gs.currentShape = gs.nextShape
		gs.nextShape = &newShape
		gs.addShape(gs.currentShape)
		return nextCmd
	}

	if !gs.applyTransformation(gs.currentShape.MoveDown) {
		_, posY := gs.currentShape.GetPosition()
		completedLines := gs.checkForCompleteLines(posY, posY+gs.currentShape.GetHeight()-1)

		gs.currentShape = nil

		if len(completedLines) != 0 {
			gs.isAnimating = true
			lineAnimationMsg := gs.constructLineAnimationMsg(completedLines)
			return gs.handleLineAnimationTick(lineAnimationMsg)
		}
	}

	return nextCmd
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

	if posX+len(shapeGrid[0]) > Width || posY+len(shapeGrid) > Height {
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

func (gs *GameState) constructLineAnimationMsg(completedLines []int) LineAnimationTick {
	completedLineMap := make(map[int][Width]color.Color, len(completedLines))

	highlightColor := color.Beige
	animationCountdown := 2

	if len(completedLines) == 3 {
		animationCountdown = 4
	}

	if len(completedLines) == 4 {
		animationCountdown = 6
	}

	highlightedLine := [Width]color.Color{}
	for i := range Width {
		highlightedLine[i] = highlightColor
	}

	for _, v := range completedLines {
		completedLineMap[v] = highlightedLine

	}

	return LineAnimationTick{
		completedLineMap,
		animationCountdown,
	}
}

func (gs *GameState) handleLineAnimationTick(animationTick LineAnimationTick) tea.Cmd {
	if animationTick.animationCountDown == 0 {
		gs.isAnimating = false
		gs.removeCompletedLines(slices.Collect(maps.Keys(animationTick.linesToUpdate)))
		return func() tea.Msg {
			return GameProgressTick{}
		}
	}

	animationTick.animationCountDown--
	newLinesToUpdateMap := make(map[int][Width]color.Color, len(animationTick.linesToUpdate))
	for k, v := range animationTick.linesToUpdate {
		newLinesToUpdateMap[k] = gs.gameBoard.Grid[k]
		gs.gameBoard.Grid[k] = v
	}

	return tea.Tick(LineAnimationInterval, func(time.Time) tea.Msg {
		return LineAnimationTick{
			newLinesToUpdateMap,
			animationTick.animationCountDown,
		}
	})
}

func (gs *GameState) removeCompletedLines(completedLines []int) {
	slices.Sort(completedLines)
	slices.Reverse(completedLines)
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

		for j := range Width {
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
	for i := range Width {
		if gs.gameBoard.Grid[line][i] == color.Black {
			return false
		}
	}

	return true
}

func (gs *GameState) isLineEmpty(line int) bool {
	for i := range Width {
		if gs.gameBoard.Grid[line][i] != color.Black {
			return false
		}
	}

	return true
}
