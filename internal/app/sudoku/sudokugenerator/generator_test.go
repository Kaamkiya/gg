package sudokugenerator

import "testing"

func TestGen(t *testing.T) {
	m := Model{}
	m.Init()

	m.grid = make([][]int, 9)
	for i := range m.grid {
		m.grid[i] = make([]int, 9)
	}
	m.generate()

	for r, row := range m.grid {
		for c, cell := range row {
			if !m.unusedInBox(r-r%3, c-c%3, cell) || !m.unusedInCol(c, cell) || !m.unusedInRow(r, cell) {
				t.Fatalf("Invalid Sudoku generated: %d overlaps", cell)
			}
		}
	}

	m.emptyCells(20)
	c := 0
	for _, r := range m.grid {
		for _, n := range r {
			if n == 0 {
				c++
			}
		}
	}

	if c != 20 {
		t.Fatalf("Not enough empty cells: wanted=20 got=%d", c)
	}
}
