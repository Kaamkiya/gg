package maze

import (
	"math/rand/v2"

	tea "github.com/charmbracelet/bubbletea"
)

type vector struct {
	x int
	y int
}

type model struct {
	maze   [][]rune
	pos    vector
	endpos vector
}

func initialModel() tea.Model {
	maze := mazes[rand.IntN(len(mazes))]

	startpos := vector{}
	endpos := vector{}

	for y := range maze {
		for x := range maze[y] {
			if maze[y][x] == 'S' {
				startpos.x = y
				startpos.y = x
			}
			if maze[y][x] == 'E' {
				endpos.x = y
				endpos.y = x
			}
		}
	}

	return model{
		maze:   maze,
		pos:    startpos,
		endpos: endpos,
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
		case "up", "w":
			m.MovePlayer("up")
		case "down", "s":
			m.MovePlayer("down")
		case "left", "a":
			m.MovePlayer("left")
		case "right", "d":
			m.MovePlayer("right")

		}
	}

	if m.pos == m.endpos {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for i, row := range m.maze {
		for j := range m.maze[i] {
			if i == m.pos.x && j == m.pos.y {
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

	s += "\n\nwasd or arrows to move\n"

	return s
}

func (m *model) MovePlayer(dir string) {
	switch dir {
	case "left":
		m.pos.y--
		if m.maze[m.pos.x][m.pos.y] == '#' {
			m.pos.y++
		}
	case "right":
		m.pos.y++
		if m.maze[m.pos.x][m.pos.y] == '#' {
			m.pos.y--
		}
	case "up":
		m.pos.x--
		if m.maze[m.pos.x][m.pos.y] == '#' {
			m.pos.x++
		}
	case "down":
		m.pos.x++
		if m.maze[m.pos.x][m.pos.y] == '#' {
			m.pos.x--
		}
	}
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
