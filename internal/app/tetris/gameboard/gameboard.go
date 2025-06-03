package gameboard

import (
	"strings"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	"github.com/charmbracelet/lipgloss"
)

const Height = 20
const Width = 10

type Gameboard struct {
	Colors map[color.Color]lipgloss.Style
	Grid   [Height][Width]color.Color
}

func NewGameboard(colors map[color.Color]lipgloss.Style) *Gameboard {
	grid := [Height][Width]color.Color{}

	return &Gameboard{colors, grid}
}

func (board *Gameboard) Render() string {
	boardBuilder := strings.Builder{}
	boardBuilder.Grow(Height * Width * 4)

	for i := range Height {
		lineBuilder := strings.Builder{}
		lineBuilder.Grow(Width * 2)

		for j := range Width {
			nextChar := board.Colors[board.Grid[i][j]].Render("  ")
			lineBuilder.WriteString(nextChar + nextChar)
		}

		lineBuilder.WriteString("\n")

		line := lineBuilder.String()
		boardBuilder.WriteString(line)
		boardBuilder.WriteString(line)

	}

	return boardBuilder.String()
}
