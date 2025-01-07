package mazegenerator

import tea "github.com/charmbracelet/bubbletea"

type MazeModel struct {
	maze *Maze
}

const (
	width  = 30
	height = 20
	algo   = "prim"
)

func GetModel() tea.Model {
	maze := GenerateMaze(width, height, algo)

	return MazeModel{
		maze,
	}
}

func (m MazeModel) Init() tea.Cmd {
	return nil
}

func (m MazeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "g":
			m.generate()
		}
	}

	return m, nil
}

func (m MazeModel) View() string {
	s := "\n"

	startX, startY := m.maze.GetStartPos()

	for i, row := range m.maze.Grid {
		for j := range m.maze.Grid[i] {
			if i == startX && j == startY {
				s += "@"
			} else if row[j] == 'E' {
				s += "X"
			} else if row[j] == '#' {
				s += string(rune(9608))
			} else {
				s += " "
			}
		}
		s += "\n"
	}

	s += "\n[G]enerate new maze \n"

	return s
}

func (m *MazeModel) generate() {
	m.maze = GenerateMaze(width, height, algo)
}
