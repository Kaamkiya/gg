package sudokugenerator

import "testing"

func TestGen(t *testing.T) {
	m := Model{}
	m.Init()

	m.Grid = make([][]int, 9)
	for i := range m.Grid {
		m.Grid[i] = make([]int, 9)
	}
	m.generate()

	for i := range 9 {
		digit := i + 1
		for r := range 9 {
			if m.unusedInCol(r, digit) {
				t.Fatalf("Digit %d is not present in column %d", digit, r)
			}
			if m.unusedInRow(r, digit) {
				t.Fatalf("Digit %d is not present in row %d", digit, r)
			}
			if m.unusedInBox(r-r%3, r-r%3, digit) {
				t.Fatalf("Digit %d is not present in box starting at (%d, %d)", digit, r-r%3, r-r%3)
			}
		}
	}

	m.emptyCells(20)
	c := 0
	for _, r := range m.Grid {
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
