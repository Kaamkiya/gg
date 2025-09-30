package sudoku

import (
	"fmt"
	"strconv"

	"github.com/Kaamkiya/gg/internal/app/sudoku/sudokugenerator"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	origGrid [][]int
	grid     [][]int

	cursorx int
	cursory int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursory > 0 {
				m.cursory--
			}
		case "down", "j":
			if m.cursory < 8 {
				m.cursory++
			}
		case "left", "h":
			if m.cursorx > 0 {
				m.cursorx--
			}
		case "right", "l":
			if m.cursorx < 8 {
				m.cursorx++
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			m.setSquare(msg.String())
		}
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for i, r := range m.grid {
		for j, c := range r {
			if j%3 == 0 && j != 0 {
				s += " | "
			}

			if j == m.cursorx && i == m.cursory {
				col := lipgloss.NewStyle().Background(lipgloss.Color("#0000ff")).Render
				if c == 0 {
					s += col(" . ")
				} else {
					s += col(fmt.Sprintf(" %d ", c))
				}
			} else {
				if c == 0 {
					s += " . "
				} else {
					s += fmt.Sprintf(" %d ", c)
				}
			}
		}

		s += "\n"

		if i == 2 || i == 5 {
			s += "--------------------------------\n"
		}
	}

	s += fmt.Sprintf("\n\norig: %v\n\ncurr: %v", m.origGrid, m.grid)

	return s
}

func (m *model) setSquare(button string) {
	if m.origGrid[m.cursory][m.cursorx] == 0 {
		m.grid[m.cursory][m.cursorx], _ = strconv.Atoi(button)
	}
}

func initialModel() tea.Model {
	g := sudokugenerator.Model{}
	g.Init()

	grid := make([][]int, 9)
	orig := make([][]int, 9)

	for i := range 9 {
		grid[i] = make([]int, 9)
		orig[i] = make([]int, 9)

		for j := range 9 {
			grid[i][j] = g.Grid[j][i]
			orig[i][j] = g.Grid[j][i]
		}
	}

	return model{
		grid:     grid,
		origGrid: orig,
	}
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
