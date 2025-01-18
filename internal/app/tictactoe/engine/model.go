package engine

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Game struct {
	board    *Board
	engine   *Engine
	turn     Player
	winner   Player
	gameover bool
	colors   map[string]lipgloss.Style
}

const (
	size = 3
)

func GetModel() tea.Model {
	board := NewBoard(size)
	engine := NewEngine(100)

	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f9f6f2"))
	c := func(s string) lipgloss.Color {
		return lipgloss.Color(s)
	}

	const (
		yellow = "#FF9E3B"
		dark   = "#3C3A32"
		gray   = "#717C7C"
		light  = "#DCD7BA"
		red    = "#E63D3D"
		green  = "#98BB6C"
		blue   = "#7E9CD8"
	)

	return Game{
		board:    board,
		engine:   engine,
		turn:     P1,
		winner:   0,
		gameover: false,
		colors: map[string]lipgloss.Style{
			"board":  defaultStyle.Background(c(dark)),
			"text":   defaultStyle.Background(c(dark)).Foreground(c(light)),
			"line":   defaultStyle.Background(c(dark)).Foreground(c(gray)),
			"p1":     defaultStyle.Background(c(dark)).Foreground(c(yellow)),
			"p2":     defaultStyle.Background(c(dark)).Foreground(c(red)),
			"hi":     defaultStyle.Foreground(c(green)),
			"status": defaultStyle.Foreground(c(blue)),
		},
	}
}

func (g Game) Init() tea.Cmd {
	return nil
}

func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gameOverMsg:
		g.winner = msg.winner
		g.turn = g.engine.GetOpponent(g.turn)
		g.gameover = true
		return g, nil

	case nextTurnMsg:
		g.turn = g.engine.GetOpponent(g.turn)
		return g, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return g, tea.Quit
		case "n", "N":
			g.nextMatch()

			if g.turn == P2 {
				return g, aiTurnCmd(&g)
			}

			return g, nil
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// There shouldn't be an error, because this is only called for integers
			index, _ := strconv.Atoi(msg.String())
			index -= 1
			cell, err := g.board.GetCell(index)
			if err != nil {
				log.Fatal(err)
			}

			if cell == EMPTY {
				g.engine.PlayMove(g.board, P1, index)
				g.turn = g.engine.GetOpponent(g.turn)

				isover, win := g.engine.CheckGameOver(g.board, index)

				if isover {
					if win > 0 {
						g.winner = g.turn
					} else {
						g.winner = 0
					}
					g.gameover = true
					return g, nil
				}

				if g.turn == P2 {
					return g, aiTurnCmd(&g)
				}
			}
		}
	}

	return g, nil
}

type gameOverMsg struct {
	winner Player
}

type nextTurnMsg struct {
	move int
}

// Handle AI turn
func aiTurnCmd(g *Game) tea.Cmd {
	return func() tea.Msg {
		rollout := g.board.Copy()
		move := g.engine.ai.Solve(rollout)

		g.engine.PlayMove(g.board, P2, move)

		isover, win := g.engine.CheckGameOver(g.board, move)
		if isover {
			if win > 0 {
				return gameOverMsg{winner: P2}
			}

			return gameOverMsg{winner: 0}
		}

		return nextTurnMsg{}
	}
}

func (g *Game) nextMatch() {
	g.board = NewBoard(size)
	g.gameover = false
	g.winner = 0
	randLvl := rand.IntN(50) + 50
	g.engine = NewEngine(randLvl)
}

func printCell(board *Board, index int) string {
	cell, err := board.GetCell(index)
	if err != nil {
		panic(err)
	}

	sign := printPlayer(cell)

	if sign == "" {
		return fmt.Sprintf("%d", index+1)
	}

	return sign
}

func printPlayer(cell int) string {
	if cell == P1 {
		return "O"
	} else if cell == P2 {
		return "X"
	}

	return ""
}

func (g Game) View() string {
	renderCell := func(index int) string {
		cell, _ := g.board.GetCell(index)
		var style lipgloss.Style
		content := ""

		switch cell {
		case P1:
			style = g.colors["p1"]
			content = "O"
		case P2:
			style = g.colors["p2"]
			content = "X"
		default: // Empty cell, show index
			style = g.colors["text"]
			content = strconv.Itoa(index + 1)
		}

		return style.Render(content)
	}

	board := "\n"
	for i := 0; i < 3; i++ {
		board += g.colors["board"].Render(" ")
		board += renderCell(i * 3)
		board += g.colors["line"].Render(" | ")
		board += renderCell(i*3 + 1)
		board += g.colors["line"].Render(" | ")
		board += renderCell(i*3 + 2)
		board += g.colors["board"].Render(" ")

		if i < 2 {
			board += "\n" + g.colors["line"].Render("---+---+---") + "\n"
		}
	}

	status := ""
	if g.gameover {
		if g.winner != 0 {
			status += g.colors["hi"].Render("\n  Winner: ")
			status += g.colors["hi"].Render(printPlayer(g.winner))
		} else {
			status += g.colors["hi"].Render("\n   Draw!")
		}
		status += g.colors["status"].Render("\n\n[Q]uit -- [N]ext match")
	} else {
		status = g.colors["status"].Render(fmt.Sprintf("\n  %s's turn", printPlayer(g.turn)))
	}

	return board + status
}
