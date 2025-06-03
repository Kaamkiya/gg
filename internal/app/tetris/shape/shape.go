package shape

import "github.com/Kaamkiya/gg/internal/app/tetris/color"

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

func CreateNew(posX, posY int) Shape {
	return createL(posX, posY)
}

func (s Shape) MoveDown() Shape {
	return Shape{
		s.posX,
		s.posY + 1,
		copyGrid(s.grid),
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

func copyGrid(grid [][]bool) [][]bool {
	duplicate := make([][]bool, len(grid))
	for i := range grid {
		duplicate[i] = make([]bool, len(grid[i]))
		copy(duplicate[i], grid[i])
	}

	return duplicate
}
