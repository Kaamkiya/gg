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

// height is the game area height counted in Tetris squares
const height = 20

// width is the game area height counted in Tetris squares
const width = 10

const lineAnimationInterval time.Duration = 100 * time.Millisecond

// gameboard represents the Tetris game area. The Grid is a fixed-size array
// where each box contains a Color. The color of each box is used both for displaying
// and for calculating if a box is empty where if the box is of color Black, it is
// considered empty.
type gameboard struct {
	Colors map[color.Color]lipgloss.Style
	Grid   [height][width]color.Color
}

// lineAnimationTick is a tea.Msg that contains which lines the lines to change
// and their colors and how many animations (color changes) are left for the
// animation to complete.
type lineAnimationTick struct {
	linesToUpdate      map[int][width]color.Color
	animationCountDown int
}

func newGameboard(colors map[color.Color]lipgloss.Style) *gameboard {
	grid := [height][width]color.Color{}

	return &gameboard{colors, grid}
}

// gameState contains the application state.
//   - nextShape is the shape that will be dropped after the current one.
//   - currentShape is the shape that is being dropped currently.
//   - gameboard is the playing area
//   - shapeRandomizer is used to find which shape is going to be dropped next.
//   - isPaused is a flag which is true when the game is paused.
type gameState struct {
	nextShape       *shape.Shape
	currentShape    *shape.Shape
	gameBoard       *gameboard
	shapeRandomizer *shape.Randomizer
	score           uint
	isPaused        bool
}

// handleGameProgressTick updates the game state to simulate the current shape
// dropping a line. The basic flow is:
//  1. Create new shapes if needed
//  2. Drop the current shape one line
//  3. Check if any lines are completed
//  4. Start the line clearing animation if needed, otherwise schedule the
//     the next tick.
func (gs *gameState) handleGameProgressTick() tea.Cmd {
	middleX := (width / 2) - 1
	if gs.nextShape == nil {
		newShape := shape.CreateNew(middleX, 0, gs.shapeRandomizer)
		gs.nextShape = &newShape
	}

	nextCmd := tea.Tick(gameProgressTickDelay, func(time.Time) tea.Msg {
		return gameProgressTick{}
	})

	if gs.currentShape == nil {
		newShape := shape.CreateNew(middleX, 0, gs.shapeRandomizer)
		gs.currentShape = gs.nextShape
		gs.nextShape = &newShape
		gs.addShape(gs.currentShape)
		return nextCmd
	}

	gs.score += 1

	if !gs.applyTransformation(gs.currentShape.MoveDown) {
		_, posY := gs.currentShape.GetPosition()
		completedLines := gs.checkForCompleteLines(posY, posY+gs.currentShape.GetHeight()-1)

		gs.currentShape = nil

		if len(completedLines) != 0 {
			lineAnimationMsg := gs.constructLineAnimationMsg(completedLines)
			return gs.handleLineAnimationTick(lineAnimationMsg)
		} else if posY == 0 {
			return tea.Quit
		}
	}

	return nextCmd
}

func (gs *gameState) handleLeft() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.MoveLeft)
}

func (gs *gameState) handleRight() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.MoveRight)
}

func (gs *gameState) handleDown() {
	if gs.currentShape == nil {
		return
	}

	succesfull := gs.applyTransformation(gs.currentShape.MoveDown)

	if succesfull {
		gs.score += 2
	}
}

func (gs *gameState) handleLeftRotate() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.RotateLeft)
}

func (gs *gameState) handleRightRotate() {
	if gs.currentShape == nil {
		return
	}

	gs.applyTransformation(gs.currentShape.RotateRight)
}

func (gs *gameState) applyTransformation(tranformation func() shape.Shape) bool {
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

// isShapeValid checks if a shape is valid by checking:
//   - If the shape is inside the gameBoard
//   - If the shape does not overlap with any occupied box.
func (gs *gameState) isShapeValid(shape shape.Shape) bool {
	shapeGrid := shape.GetGrid()
	posX, posY := shape.GetPosition()

	if posX < 0 {
		return false
	}

	if posX+len(shapeGrid[0]) > width || posY+len(shapeGrid) > height {
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

func (gs *gameState) addShape(shape *shape.Shape) {
	gs.modidfyColorGridFromShape(shape, shape.GetColor())
}

func (gs *gameState) deleteShape(shape *shape.Shape) {
	gs.modidfyColorGridFromShape(shape, color.Black)
}

func (gs *gameState) modidfyColorGridFromShape(shape *shape.Shape, color color.Color) {
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

func (gs *gameState) constructLineAnimationMsg(completedLines []int) lineAnimationTick {
	completedLineMap := make(map[int][width]color.Color, len(completedLines))
	animationCountdown := 2

	if len(completedLines) == 3 {
		animationCountdown = 4
	}

	if len(completedLines) == 4 {
		animationCountdown = 6
	}

	highlightedLine := [width]color.Color{}
	for i := range width {
		highlightedLine[i] = color.Beige
	}

	for _, v := range completedLines {
		completedLineMap[v] = highlightedLine

	}

	return lineAnimationTick{
		completedLineMap,
		animationCountdown,
	}
}

// handleLineAnimationTick performs the state updates for the flashing animation when
// lines are completed. If the animation is complete (animationCountDown set to 0) it
// resumes the game. Otherwse it swaps the lines color and continues with the animation.
func (gs *gameState) handleLineAnimationTick(animationTick lineAnimationTick) tea.Cmd {
	if animationTick.animationCountDown == 0 {
		gs.removeCompletedLines(slices.Collect(maps.Keys(animationTick.linesToUpdate)))
		return func() tea.Msg {
			return gameProgressTick{}
		}
	}

	animationTick.animationCountDown--
	newLinesToUpdateMap := make(map[int][width]color.Color, len(animationTick.linesToUpdate))
	for k, v := range animationTick.linesToUpdate {
		newLinesToUpdateMap[k] = gs.gameBoard.Grid[k]
		gs.gameBoard.Grid[k] = v
	}

	return tea.Tick(lineAnimationInterval, func(time.Time) tea.Msg {
		return lineAnimationTick{
			newLinesToUpdateMap,
			animationTick.animationCountDown,
		}
	})
}

func (gs *gameState) removeCompletedLines(completedLines []int) {
	gs.addLineScore(len(completedLines))
	slices.Sort(completedLines)
	slices.Reverse(completedLines)

	// lines are removed with a single pass from bottom to top.The completedLines array
	// is sorted in descending order and the first completed line is replaced by the one
	// above it. If another completed line is encountered during replacing, the distanceToCopyFrom
	// is increased to start copying from two places above etc. The distanceToCopyFrom variable
	// specifies both the lines to skip when replacing and the index of the next completed line in the
	// completedLines array.
	distanceToCopyFrom := 1

	for i := completedLines[0]; i >= 0; i-- {
		if i-distanceToCopyFrom < 0 {
			return
		}

		if gs.isLineEmpty(i) {
			return
		}

		for distanceToCopyFrom < len(completedLines) && completedLines[distanceToCopyFrom] == i-distanceToCopyFrom {
			distanceToCopyFrom++
		}

		for j := range width {
			gs.gameBoard.Grid[i][j] = gs.gameBoard.Grid[i-distanceToCopyFrom][j]
		}
	}
}

func (gs *gameState) addLineScore(completedLinesNum int) {
	switch completedLinesNum {
	case 4:
		gs.score += 800
	case 3:
		gs.score += 500
	case 2:
		gs.score += 300
	default:
		gs.score += 100
	}
}

func (gs *gameState) checkForCompleteLines(from, to int) []int {
	completedLines := make([]int, 0, 4)
	for i := to; i >= from; i-- {
		if gs.isLineCompleted(i) {
			completedLines = append(completedLines, i)
		}
	}

	return completedLines
}

func (gs *gameState) isLineCompleted(line int) bool {
	for i := range width {
		if gs.gameBoard.Grid[line][i] == color.Black {
			return false
		}
	}

	return true
}

func (gs *gameState) isLineEmpty(line int) bool {
	for i := range width {
		if gs.gameBoard.Grid[line][i] != color.Black {
			return false
		}
	}

	return true
}
