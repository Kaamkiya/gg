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
	dx := c.x - other.x
	dy := c.y - other.y
	return dx*dx + dy*dy
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

func (m Maze) IsBoundary(x, y int) bool {
	vertical := (x == 0 || x == m.Width-1) && y >= 0 && y <= m.Height-1
	horizontal := (y == 0 || y == m.Height-1) && x >= 0 && x <= m.Width-1

	return vertical || horizontal
}

func (m Maze) IsWall(x, y int) bool {
	return m.Grid[y][x] == WALL
}

func (m Maze) GetFrontiers(x, y int, findWall bool) []Cell {
	var frontiers []Cell
	for _, dir := range DIRS {
		dx, dy := x+2*dir.x, y+2*dir.y
		if !m.IsInner(dx, dy) {
			if findWall && m.IsBoundary(dx, dy) {
				frontiers = append(frontiers, Cell{dx, dy})
			}
			continue
		}
		if m.IsWall(dx, dy) == findWall && m.IsWall(x+dir.x, y+dir.y) {
			frontiers = append(frontiers, Cell{dx, dy})
		}
	}

	return frontiers
}

func (m *Maze) MakePath(cell Cell) {
	if !m.IsInner(cell.x, cell.y) || !m.IsWall(cell.x, cell.y) {
		return
	}
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
