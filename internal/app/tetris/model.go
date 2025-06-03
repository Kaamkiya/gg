package tetris

import (
	"strings"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() GameState {
	return GameState{
		nil,
		nil,
		NewGameboard(color.Colors),
	}
}

func (gs *GameState) Init() tea.Cmd {
	return nil
}

func (gs *GameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q":
			return gs, tea.Quit
		case "h", "H", "left":
			gs.HandleLeft()
		case "l", "L", "right":
			gs.HandleRight()
		case "j", "J", "down":
			gs.HandleDown()
		case "z", "Z":
			gs.HandleLeftRotate()
		case "x", "X":
			gs.HandleRightRotate()
		}
	case TickMsg:
		gs.HandleTick()
	}

	return gs, nil
}

func (gs *GameState) View() string {
	boardBuilder := strings.Builder{}
	boardBuilder.Grow(Height * Width * 4)

	for i := range Height {
		lineBuilder := strings.Builder{}
		lineBuilder.Grow(Width * 2)

		for j := range Width {
			nextChar := gs.gameBoard.Colors[gs.gameBoard.Grid[i][j]].Render("  ")
			lineBuilder.WriteString(nextChar + nextChar)
		}

		lineBuilder.WriteString("\n")

		line := lineBuilder.String()
		boardBuilder.WriteString(line)
		boardBuilder.WriteString(line)

	}

	return boardBuilder.String()
}
