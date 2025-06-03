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
	return gs.gameBoard.Render()
}
