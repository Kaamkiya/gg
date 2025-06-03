package tetris

import (
	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	"github.com/Kaamkiya/gg/internal/app/tetris/gameboard"
	tea "github.com/charmbracelet/bubbletea"
)

func initialModel() GameState {
	return GameState{
		nil,
		nil,
		gameboard.NewGameboard(color.Colors),
	}
}

func (gs *GameState) Init() tea.Cmd {
	return nil
}

func (gs *GameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return gs, tea.Quit
		case "a":
			gs.HandleTick()
		}
	}

	return gs, nil
}

func (gs *GameState) View() string {
	return gs.gameBoard.Render()
}
