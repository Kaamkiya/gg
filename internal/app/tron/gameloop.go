package tron

import (
	"strings"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
	tea "github.com/charmbracelet/bubbletea"
)

const tickDelay = 50 * time.Millisecond

type gameProgressTick struct{}

// Init implements [tea.Model].
func (gs *gameState) Init() tea.Cmd {
	return func() tea.Msg {
		return gameProgressTick{}
	}
}

// Update implements [tea.Model].
func (gs *gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "Q" {
			return gs, tea.Quit
		} else {
			switch msg.String() {
			case "w":
				gs.players[0].changeDirection(gameboard.Up)
			case "a":
				gs.players[0].changeDirection(gameboard.Left)
			case "s":
				gs.players[0].changeDirection(gameboard.Down)
			case "d":
				gs.players[0].changeDirection(gameboard.Right)
			case "i":
				gs.players[1].changeDirection(gameboard.Up)
			case "j":
				gs.players[1].changeDirection(gameboard.Left)
			case "k":
				gs.players[1].changeDirection(gameboard.Down)
			case "l":
				gs.players[1].changeDirection(gameboard.Right)
			}
		}
	case gameProgressTick:
		gs.updatePlayers()
		gs.gameEnded, gs.winner = gs.checkForWinner()
		return gs, tea.Tick(tickDelay, func(time.Time) tea.Msg {
			return gameProgressTick{}
		})

	}
	return gs, nil
}

// View implements [tea.Model].
func (gs *gameState) View() string {
	viewBuilder := strings.Builder{}
	viewBuilder.Grow((gameboard.Height + 5) * (gameboard.Width + 2) * 2)

	viewBuilder.WriteString(strings.Repeat(gameboard.Colors[gameboard.Beige].Render(" "), gameboard.Width*2+2))
	viewBuilder.WriteString("\n")

	for y := range gameboard.Height {
		lineBuilder := strings.Builder{}
		lineBuilder.Grow(gameboard.Width * 2)

		for x := range gameboard.Width {
			if x == 0 {
				lineBuilder.WriteString(gameboard.Colors[gameboard.Beige].Render(" "))
			}

			nextChar := gameboard.Colors[gs.gameboard.Grid[y][x]].Render("  ")
			lineBuilder.WriteString(nextChar)

			if x == gameboard.Width-1 {
				lineBuilder.WriteString(gameboard.Colors[gameboard.Beige].Render(" "))
			}
		}

		viewBuilder.WriteString(lineBuilder.String())
		viewBuilder.WriteString("\n")
	}

	viewBuilder.WriteString(strings.Repeat(gameboard.Colors[gameboard.Beige].Render(" "), gameboard.Width*2+2))
	viewBuilder.WriteString("\n")

	return viewBuilder.String()
}
