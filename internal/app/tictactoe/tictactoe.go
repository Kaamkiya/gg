package tictactoe

import (
	"fmt"
	"strconv"

	"github.com/Kaamkiya/gg/internal/app/tictactoe/engine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var winner = ' '

type model struct {
	turn   rune
	board  [9]rune
	xcolor lipgloss.Style
	ocolor lipgloss.Style
}

func initialModel() tea.Model {
	return model{
		turn: 'x',
		board: [9]rune{
			'1', '2', '3',
			'4', '5', '6',
			'7', '8', '9',
		},
		xcolor: lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")),
		ocolor: lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff")),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// There shouldn't be an error, because this is only called for integers
			position, _ := strconv.Atoi(msg.String())

			if m.board[position-1] != 'x' && m.board[position-1] != 'o' {
				m.board[position-1] = m.turn

				if m.turn == 'x' {
					m.turn = 'o'
				} else {
					m.turn = 'x'
				}
			}

			if m.CheckForWin() != ' ' {
				winner = m.CheckForWin()
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("%c | %c | %c\n", m.board[0], m.board[1], m.board[2])
	s += "---------\n"
	s += fmt.Sprintf("%c | %c | %c\n", m.board[3], m.board[4], m.board[5])
	s += "---------\n"
	s += fmt.Sprintf("%c | %c | %c\n", m.board[6], m.board[7], m.board[8])

	s += fmt.Sprintf("\n\n%c's turn", m.turn)

	return s
}

func (m model) CheckForWin() rune {
	// Check over each row to see if someone won.
	for i := 0; i < 9; i += 3 {
		o := 0
		x := 0
		for _, c := range m.board[i : i+3] {
			if c == 'x' {
				x++
			} else if c == 'o' {
				o++
			}
		}

		if o == 3 {
			return 'o'
		} else if x == 3 {
			return 'x'
		}
	}

	// Check the columns
	for i := 0; i < 3; i++ {
		if m.board[i] == 'x' && m.board[i+3] == 'x' && m.board[i+6] == 'x' {
			return 'x'
		}
		if m.board[i] == 'o' && m.board[i+3] == 'o' && m.board[i+6] == 'o' {
			return 'o'
		}
	}

	// And finally, check the diagonals.
	if m.board[0] == 'x' && m.board[4] == 'x' && m.board[8] == 'x' ||
		m.board[2] == 'x' && m.board[4] == 'x' && m.board[6] == 'x' {
		return 'x'
	}

	if m.board[0] == 'o' && m.board[4] == 'o' && m.board[8] == 'o' ||
		m.board[2] == 'o' && m.board[4] == 'o' && m.board[6] == 'o' {
		return 'o'
	}

	return ' '
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}

	fmt.Printf("%c wins\n", winner)
}

func RunVsAi() {
	p := tea.NewProgram(engine.GetModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
