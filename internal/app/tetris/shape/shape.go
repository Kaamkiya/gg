// Package shape is responsible for creating and transforming the game shapes.
// The constructed shape cannot be modified outside the package and each
// transformation is producing a new shape.
package shape

import (
	"github.com/Kaamkiya/gg/internal/app/tetris/color"
)

const (
	I int = iota
	L
	J
	T
	Z
	S
	O
)

type Shape struct {
	posX  int
	posY  int
	grid  [][]bool
	color color.Color
}

func createI(posX int, posY int) Shape {
	return Shape{
		posX,
		posY,
		[][]bool{
			{true},
			{true},
			{true},
			{true},
		},
		color.Teal,
	}
}

func createJ(posX int, posY int) Shape {
	return Shape{
		posX,
		posY,
		[][]bool{
			{false, true},
			{false, true},
			{true, true},
		},
		color.Green,
	}
}

func createL(posX int, posY int) Shape {
	return Shape{
		posX,
		posY,
		[][]bool{
			{true, false},
			{true, false},
			{true, true},
		},
		color.Orange,
	}
}

func createZ(posX int, posY int) Shape {
	return Shape{
		posX,
		posY,
		[][]bool{
			{false, true},
			{true, true},
			{true, false},
		},
		color.Purple,
	}
}

func createS(posX int, posY int) Shape {
	return Shape{
		posX,
		posY,
		[][]bool{
			{true, false},
			{true, true},
			{false, true},
		},
		color.Pink,
	}
}

func createO(posX int, posY int) Shape {
	return Shape{
		posX,
		posY,
		[][]bool{
			{true, true},
			{true, true},
		},
		color.Blue,
	}
}

func createT(posX int, posY int) Shape {
	return Shape{
		posX,
		posY,
		[][]bool{
			{true, true, true},
			{false, true, false},
		},
		color.Magenta,
	}
}

func CreateNew(posX, posY int, randomizer *Randomizer) Shape {
	switch randomizer.nextInt(7) {
	case L:
		return createL(posX, posY)
	case I:
		return createI(posX, posY)
	case J:
		return createJ(posX, posY)
	case O:
		return createO(posX, posY)
	case S:
		return createS(posX, posY)
	case Z:
		return createZ(posX, posY)
	default:
		return createT(posX, posY)
	}
}

func (s Shape) MoveDown() Shape {
	return Shape{
		s.posX,
		s.posY + 1,
		copyGrid(s.grid),
		s.color,
	}
}

func (s Shape) MoveRight() Shape {
	return Shape{
		s.posX + 1,
		s.posY,
		copyGrid(s.grid),
		s.color,
	}
}

func (s Shape) MoveLeft() Shape {
	return Shape{
		s.posX - 1,
		s.posY,
		copyGrid(s.grid),
		s.color,
	}
}

func (s Shape) RotateRight() Shape {
	newGrid := make([][]bool, len(s.grid[0]))

	for i := range s.grid[0] {
		newLine := make([]bool, len(s.grid))
		for j := range s.grid {
			newLine[len(s.grid)-1-j] = s.grid[j][i]
		}
		newGrid[i] = newLine
	}

	return Shape{
		s.posX,
		s.posY,
		newGrid,
		s.color,
	}
}

func (s Shape) RotateLeft() Shape {
	newGrid := make([][]bool, len(s.grid[0]))

	for i := range s.grid[0] {
		newLine := make([]bool, len(s.grid))

		for j := range s.grid {
			newLine[j] = s.grid[j][i]
		}

		newGrid[len(s.grid[0])-1-i] = newLine
	}

	return Shape{
		s.posX,
		s.posY,
		newGrid,
		s.color,
	}
}

func (s Shape) GetColor() color.Color {
	return s.color
}

func (s Shape) GetPosition() (int, int) {
	return s.posX, s.posY
}

func (s Shape) GetGrid() [][]bool {
	return copyGrid(s.grid)
}

func (s Shape) GetHeight() int {
	return len(s.grid)
}

func copyGrid(grid [][]bool) [][]bool {
	duplicate := make([][]bool, len(grid))
	for i := range grid {
		duplicate[i] = make([]bool, len(grid[i]))
		copy(duplicate[i], grid[i])
	}

	return duplicate
}
