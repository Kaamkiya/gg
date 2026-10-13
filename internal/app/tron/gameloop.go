package tron

import (
	"fmt"
	"strings"

	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
	tea "github.com/charmbracelet/bubbletea"
)

type gameProgressTick struct{}

// Init implements [tea.Model].
func (gs *gameState) Init() tea.Cmd {
	return nil
}

// Update implements [tea.Model].
func (gs *gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "Q" {
			return gs, tea.Quit
		}
	case gameProgressTick:
		fmt.Println("asdf")

	}
	return gs, nil
}

// View implements [tea.Model].
func (gs *gameState) View() string {
	viewBuilder := strings.Builder{}
	viewBuilder.Grow(gameboard.Height * gameboard.Width * 2)

	for y := range gameboard.Height {
		lineBuilder := strings.Builder{}
		lineBuilder.Grow(gameboard.Width * 2)

		for x := range gameboard.Width {
			nextChar := gameboard.Colors[gs.gameboard.Grid[y][x]].Render("  ")
			lineBuilder.WriteString(nextChar)
		}

		viewBuilder.WriteString(lineBuilder.String())
		viewBuilder.WriteString("\n")
	}

	return viewBuilder.String()
}
