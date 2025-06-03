package tetris

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const height = 20
const width = 10

type gameboard struct {
	colors map[Color]lipgloss.Style
	grid   [height][width]Color
}

func NewGameboard(colors map[Color]lipgloss.Style) *gameboard {
	grid := [height][width]Color{}

	return &gameboard{colors, grid}
}

func (board *gameboard) Render() string {
	boardBuilder := strings.Builder{}
	boardBuilder.Grow(height * width * 4)

	for i := 0; i < height; i++ {
		lineBuilder := strings.Builder{}
		lineBuilder.Grow(width * 2)

		for j := 0; j < width; j++ {
			nextChar := board.colors[Teal].Render("  ")
			lineBuilder.WriteString(nextChar + nextChar)
		}

		lineBuilder.WriteString("\n")

		line := lineBuilder.String()
		boardBuilder.WriteString(line)
		boardBuilder.WriteString(line)

	}

	return boardBuilder.String()
}
