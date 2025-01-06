package mazegenerator

import (
	"fmt"
	"math/rand"
)

const (
	WALL  = '#'
	PATH  = ' '
	START = 'S'
	END   = 'E'
)

type Cell struct {
	x, y int
}

var DIRS = []Cell{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func (c Cell) Diff(other Cell) int {
	return (c.x-other.x)*(c.x-other.x) + (c.y-other.y)*(c.y-other.y)
}

type Maze struct {
	Width, Height int
	Start, End    Cell
	Grid          [][]rune
}

func NewMaze(width, height int) *Maze {
	grid := make([][]rune, height)

	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = WALL
		}
	}

	startX := rand.Intn(width/4) + 1
	startY := rand.Intn(height/4) + 1

	grid[startY][startX] = START

	return &Maze{
		Width:  width,
		Height: height,
		Start:  Cell{startX, startY},
		Grid:   grid,
	}
}

func (m *Maze) Set(x, y int, val rune) {
	m.Grid[y][x] = val
}

func (m Maze) Get(x, y int) rune {
	return m.Grid[y][x]
}

func (m Maze) GetStartPos() (x, y int) {
	return m.Start.x, m.Start.y
}

func (m Maze) GetEndPos() (x, y int) {
	return m.End.x, m.End.y
}

func (m *Maze) SetEnd(x, y int) {
	m.Grid[y][x] = END
	m.End = Cell{x, y}
}

func (m Maze) IsInner(x, y int) bool {
	return x > 0 && x < m.Width-1 && y > 0 && y < m.Height-1
}

func (m Maze) IsWall(x, y int) bool {
	return m.Grid[y][x] == WALL
}

func (m Maze) GetNeighbors(x, y int, findWall bool) []Cell {
	var neighbors []Cell
	for _, dir := range DIRS {
		dx, dy := x+dir.x, y+dir.y
		if !m.IsInner(dx, dy) {
			continue
		}
		if findWall == m.IsWall(dx, dy) {
			neighbors = append(neighbors, Cell{dx, dy})
		}
	}

	return neighbors
}

func (m *Maze) MakePath(cell Cell) {
	m.Set(cell.x, cell.y, PATH)
}

func (m Maze) Print() {
	for _, row := range m.Grid {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}
