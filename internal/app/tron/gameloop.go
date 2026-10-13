package tron

import (
	"strings"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	viewBuilder.Grow((gameboard.Height + 8) * (gameboard.Width + 3) * 2)

	message := gs.buildMessage()

	for _, line := range message {
		viewBuilder.WriteString(centerText(line, gameboard.Width*2+2))
		viewBuilder.WriteString("\n")
	}

	for i := 0; i < 4-len(message); i++ {
		viewBuilder.WriteString("\n")
	}

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

func (gs *gameState) buildMessage() []string {
	if gs.gameEnded {
		if gs.winner >= 0 {
			winningColor := gs.players[gs.winner].Color
			return []string{
				lipgloss.NewStyle().Foreground(lipgloss.Color(winningColor.Hex())).Render(winningColor.String()) + " wins!",
				"Press 'R' to restart",
			}
		} else {
			return []string{
				"Draw!",
				"Press 'R' to restart",
			}
		}
	}

	return []string{}
}

func centerText(text string, size int) string {
	fillSize := (size - lipgloss.Width(text)) / 2
	filling := strings.Repeat(gameboard.Colors[gameboard.None].Render(" "), fillSize)
	return filling + text + filling
}
