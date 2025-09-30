package connect4

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	board [6][7]rune // [y][x]
	turn  rune

	xStyle lipgloss.Style
	oStyle lipgloss.Style
}

func initialModel() tea.Model {
	board := [6][7]rune{}
	for y := range board {
		for x := range board[y] {
			board[y][x] = ' '
		}
	}

	return model{
		board:  board,
		turn:   'x',
		xStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		oStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
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
		case "1", "2", "3", "4", "5", "6", "7":
			/* Don't check for errors because there can't be one.
			 * This only gets called if an integer was inputted.
			 */
			col, _ := strconv.Atoi(msg.String())
			col-- // Go is 0 indexed, inputs are not.

			// A piece can only go in that column if it's not full.
			if m.board[0][col] == ' ' {
				for y := len(m.board) - 1; y >= 0; y-- {
					if m.board[y][col] == ' ' {
						m.board[y][col] = m.turn

						if m.turn == 'x' {
							m.turn = 'o'
						} else {
							m.turn = 'x'
						}

						break
					}
				}
			}
		}
	}

	if m.CheckForWin() != ' ' {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	s := "| 1 | 2 | 3 | 4 | 5 | 6 | 7 |\n"
	s += "+---------------------------+\n"

	for _, row := range m.board {
		s += "| "
		for _, cell := range row {
			style := m.oStyle

			if cell == 'x' {
				style = m.xStyle
			}

			s += style.Render(string(cell)) + " | "
		}
		s += "\n"
	}

	s += "+---------------------------+\n"

	switch m.CheckForWin() {
	case ' ':
		s += fmt.Sprintf("\n%c's turn\n", m.turn)
	case 't':
		s += "\ntie!\n"
	default:
		s += fmt.Sprintf("\n%c wins!\n", m.CheckForWin())
	}

	return s
}

func (m model) CheckForWin() rune {
	// Check for a win horizontally.
	for y := 0; y < len(m.board[0]); y++ {
		for x := 0; x < len(m.board)-3; x++ {
			tile1 := m.board[x][y]
			tile2 := m.board[x+1][y]
			tile3 := m.board[x+2][y]
			tile4 := m.board[x+3][y]
			if tile1 == 'x' && tile2 == 'x' && tile3 == 'x' && tile4 == 'x' {
				return 'x'
			}
			if tile1 == 'o' && tile2 == 'o' && tile3 == 'o' && tile4 == 'o' {
				return 'o'
			}
		}
	}

	// Check for a win vertically.
	for y := 0; y < len(m.board[0])-3; y++ {
		for x := 0; x < len(m.board); x++ {
			tile1 := m.board[x][y]
			tile2 := m.board[x][y+1]
			tile3 := m.board[x][y+2]
			tile4 := m.board[x][y+3]
			if tile1 == 'x' && tile2 == 'x' && tile3 == 'x' && tile4 == 'x' {
				return 'x'
			}
			if tile1 == 'o' && tile2 == 'o' && tile3 == 'o' && tile4 == 'o' {
				return 'o'
			}
		}
	}

	for y := 0; y < len(m.board[0])-3; y++ {
		for x := 0; x < len(m.board)-3; x++ {
			// Check left down diagonal.
			tile1 := m.board[x][y]
			tile2 := m.board[x+1][y+1]
			tile3 := m.board[x+2][y+2]
			tile4 := m.board[x+3][y+3]
			if tile1 == 'x' && tile2 == 'x' && tile3 == 'x' && tile4 == 'x' {
				return 'x'
			}
			if tile1 == 'o' && tile2 == 'o' && tile3 == 'o' && tile4 == 'o' {
				return 'o'
			}

			// Check right down diagonal.
			tile1 = m.board[x+3][y]
			tile2 = m.board[x+2][y+1]
			tile3 = m.board[x+1][y+2]
			tile4 = m.board[x][y+3]
			if tile1 == 'x' && tile2 == 'x' && tile3 == 'x' && tile4 == 'x' {
				return 'x'
			}
			if tile1 == 'o' && tile2 == 'o' && tile3 == 'o' && tile4 == 'o' {
				return 'o'
			}
		}
	}

	emptyCount := len(m.board) * len(m.board[0])

	for _, row := range m.board {
		for _, cell := range row {
			if cell != ' ' {
				emptyCount--
			}
		}
	}

	if emptyCount <= 0 {
		return 't'
	}

	return ' '
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
