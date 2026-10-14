package tron

import (
	"strconv"
	"strings"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const tickDelay = 50 * time.Millisecond

type gameProgressTick struct{}

type gameStartTick struct{}

// Init implements [tea.Model].
func (gs *gameState) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*1500, func(time.Time) tea.Msg {
		return gameStartTick{}
	})
}

// Update implements [tea.Model].
func (gs *gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "Q" {
			return gs, tea.Quit
		} else if gs.gameEnded && (msg.String() == "r" || msg.String() == "R") {
			gameState := newGameState()
			return &gameState, tea.Tick(time.Millisecond*1500, func(time.Time) tea.Msg {
				return gameStartTick{}
			})
		} else {
			firstPlayer := gs.players[0]
			lastPlayer := gs.players[len(gs.players)-1]
			switch msg.String() {
			case "w":
				firstPlayer.changeDirection(gameboard.Up)
			case "a":
				firstPlayer.changeDirection(gameboard.Left)
			case "s":
				firstPlayer.changeDirection(gameboard.Down)
			case "d":
				firstPlayer.changeDirection(gameboard.Right)
			case " ":
				firstPlayer.nitroBoost()
			case "up":
				lastPlayer.changeDirection(gameboard.Up)
			case "left":
				lastPlayer.changeDirection(gameboard.Left)
			case "down":
				lastPlayer.changeDirection(gameboard.Down)
			case "right":
				lastPlayer.changeDirection(gameboard.Right)
			case "0":
				lastPlayer.nitroBoost()
			}
		}
	case gameProgressTick:
		if gs.startCountDown > 0 {
			return gs, nil
		}

		gs.updatePlayers()
		gs.gameEnded, gs.winner = gs.checkForWinner()
		return gs, tea.Tick(tickDelay, func(time.Time) tea.Msg {
			return gameProgressTick{}
		})
	case gameStartTick:
		gs.startCountDown = gs.startCountDown - 1
		if gs.startCountDown > 0 {
			return gs, tea.Tick(time.Second, func(time.Time) tea.Msg {
				return gameStartTick{}
			})

		} else {
			return gs, func() tea.Msg {
				return gameProgressTick{}
			}
		}

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

	for i := 0; i < 3-len(message); i++ {
		viewBuilder.WriteString("\n")
	}

	viewBuilder.WriteString(gs.buildNitrosLine())
	viewBuilder.WriteString("\n")

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

	var color, text string
	if gs.startCountDown > 1 {
		color = gameboard.Red.Hex()
		text = strconv.Itoa(gs.startCountDown)

	} else if gs.startCountDown == 1 {
		color = gameboard.Orange.Hex()
		text = strconv.Itoa(gs.startCountDown)
	} else {
		color = gameboard.Green.Hex()
		text = "GO!"
	}

	return []string{
		lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text),
	}

}

func (gs *gameState) buildNitrosLine() string {
	nitrosStrings := make([]string, len(gs.players))
	var stringTotalSize int

	for i, player := range gs.players {
		label := player.Color.String() + ": "
		nitros := strings.Repeat("🚀", player.Nitros)
		nitrosStrings[i] = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(player.Color.Hex())).
			Render(label + nitros + strings.Repeat(" ", initialNitros-player.Nitros))
		stringTotalSize = lipgloss.Width(nitrosStrings[i]) + stringTotalSize
	}

	remainingGap := (gameboard.Width*2 + 2) - stringTotalSize
	paddingSize := remainingGap / (len(nitrosStrings) - 1)
	paddingRemainder := remainingGap % (len(nitrosStrings) - 1)

	var nitroLine strings.Builder
	for i := 0; i < len(nitrosStrings)-1; i++ {
		nitroLine.WriteString(nitrosStrings[i])

		if paddingRemainder > 0 {
			paddingRemainder--
			nitroLine.WriteString(strings.Repeat(" ", paddingSize+1))
		} else {
			nitroLine.WriteString(strings.Repeat(" ", paddingSize))
		}

	}

	nitroLine.WriteString(nitrosStrings[len(nitrosStrings)-1])

	return nitroLine.String()
}

func centerText(text string, size int) string {
	fillSize := (size - lipgloss.Width(text)) / 2
	filling := strings.Repeat(gameboard.Colors[gameboard.None].Render(" "), fillSize)
	return filling + text + filling
}
