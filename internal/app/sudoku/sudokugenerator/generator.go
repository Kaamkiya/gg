package sudokugenerator

import (
	"math/rand/v2"
	"slices"
)

type Model struct {
	Grid [][]int
}

func (m *Model) unusedInBox(row, col, n int) bool {
	for i := range 3 {
		for j := range 3 {
			if m.Grid[row+i][col+j] == n {
				return false
			}
		}
	}

	return true
}

func (m *Model) fillBox(row, col int) {
	var n int
	for i := range 3 {
		for j := range 3 {
			for !m.unusedInBox(row, col, n) {
				n = rand.IntN(9) + 1
			}
			m.Grid[row+i][col+j] = n
		}
	}
}

func (m *Model) unusedInCol(col, n int) bool {
	for j := range 9 {
		if m.Grid[j][col] == n {
			return false
		}
	}

	return true
}

func (m *Model) unusedInRow(row, n int) bool {
	return !slices.Contains(m.Grid[row], n)
}

func (m *Model) isSafe(row, col, n int) bool {
	return m.unusedInBox(row-row%3, col-col%3, n) && m.unusedInCol(col, n) && m.unusedInRow(row, n)
}

func (m *Model) fillRemaining(row, col int) bool {
	if row == 9 {
		return true
	}

	if col == 9 {
		return m.fillRemaining(row+1, 0)
	}

	if m.Grid[row][col] != 0 {
		return m.fillRemaining(row, col+1)
	}

	for n := 1; n < 10; n++ {
		if m.isSafe(row, col, n) {
			m.Grid[row][col] = n
			if m.fillRemaining(row, col+1) {
				return true
			}
			m.Grid[row][col] = 0
		}
	}

	return false
}

func (m *Model) emptyCells(amount int) {
	for amount > 0 {
		id := rand.IntN(81)
		i := id / 9
		j := id % 9

		if m.Grid[i][j] != 0 {
			m.Grid[i][j] = 0
			amount--
		}
	}
}

func (m *Model) generate() {
	for i := 0; i < 9; i += 3 {
		m.fillBox(i, i)
	}

	m.fillRemaining(0, 0)
}

func (m *Model) Init() {
	m.Grid = make([][]int, 9)
	for i := range m.Grid {
		m.Grid[i] = make([]int, 9)
	}

	m.generate()
	m.emptyCells(54)
}
