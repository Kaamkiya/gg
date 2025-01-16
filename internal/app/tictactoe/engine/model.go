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
	xcolor   lipgloss.Style
	ocolor   lipgloss.Style
}

const (
	size = 3
)

func GetModel() tea.Model {
	board := NewBoard(size)
	engine := NewEngine(100)

	return Game{
		board:    board,
		engine:   engine,
		turn:     P1,
		gameover: false,
		xcolor:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")),
		ocolor:   lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff")),
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
		return "o"
	} else if cell == P2 {
		return "x"
	}

	return ""
}

func (g Game) View() string {
	s := fmt.Sprintf("%s | %s | %s\n", printCell(g.board, 0), printCell(g.board, 1), printCell(g.board, 2))
	s += "---------\n"
	s += fmt.Sprintf("%s | %s | %s\n", printCell(g.board, 3), printCell(g.board, 4), printCell(g.board, 5))
	s += "---------\n"
	s += fmt.Sprintf("%s | %s | %s\n", printCell(g.board, 6), printCell(g.board, 7), printCell(g.board, 8))

	if g.gameover {
		if g.winner != 0 {
			s += fmt.Sprintf("\n\nWinner: %s", printPlayer(g.winner))
		} else {
			s += "\n\nDraw!\n"
		}
		s += "\n[Q]uit -- [N]ext match"
	} else {
		s += fmt.Sprintf("\n\n%s's turn", printPlayer(g.turn))
	}

	return s
}
