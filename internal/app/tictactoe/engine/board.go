package engine

import "fmt"

const (
	P1    = 1
	P2    = -1
	EMPTY = 0
)

type Player = int

type Board struct {
	Size  int
	Cells []int
}

func NewBoard(size int) *Board {
	cells := make([]int, size*size)
	for i := range cells {
		cells[i] = EMPTY
	}

	return &Board{
		Size:  size,
		Cells: cells,
	}
}

func (b *Board) GetCell(index int) (int, error) {
	if index < 0 || index >= len(b.Cells) {
		return 0, fmt.Errorf("invalid cell index: %d", index)
	}

	return b.Cells[index], nil
}

func (b *Board) SetCell(index int, player int) error {
	if index < 0 || index >= len(b.Cells) {
		return fmt.Errorf("invalid cell index: %d", index)
	}

	b.Cells[index] = player
	return nil
}

func (b *Board) Load(cells []int) error {
	if len(cells) != len(b.Cells) {
		return fmt.Errorf("invalid cells length: %d", len(cells))
	}

	copy(b.Cells, cells)
	return nil
}

func (b *Board) GetRowCol(index int) (int, int, error) {
	if index < 0 || index >= len(b.Cells) {
		return 0, 0, fmt.Errorf("invalid cell index: %d", index)
	}

	return index / b.Size, index % b.Size, nil
}

func (b *Board) ChangePerspective() {
	for i := range b.Cells {
		b.Cells[i] *= -1
	}
}

func (b *Board) Copy() *Board {
	newBoard := NewBoard(b.Size)
	copy(newBoard.Cells, b.Cells)
	return newBoard
}

func (b *Board) Print() {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			cell, _ := b.GetCell(i*b.Size + j)
			if cell == P1 {
				fmt.Print("O")
			} else if cell == P2 {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
