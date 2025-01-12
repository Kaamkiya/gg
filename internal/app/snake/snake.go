package snake

import (
	"fmt"
	"math/rand/v2"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type moveMsg struct{}

type vector struct {
	x int
	y int
}

func (v vector) add(other vector) vector {
	return vector{
		x: v.x + other.x,
		y: v.y + other.y,
	}
}
func (v vector) equals(other vector) bool {
	return v.x == other.x && v.y == other.y
}

var (
	dirUp    = vector{0, -1}
	dirDown  = vector{0, 1}
	dirLeft  = vector{-1, 0}
	dirRight = vector{1, 0}
)

type player struct {
	body []vector
	dir  vector
}

func (p *player) move(m model, foodPos vector) {
	head := p.body[0].add(p.dir)
	p.body = append([]vector{head}, p.body...)
	if !head.equals(foodPos) {
		p.body = p.body[:len(p.body)-1]
	}
}

func (p player) headChar() rune {
	if p.dir.equals(vector{0, -1}) {
		return '^'
	}
	if p.dir.equals(vector{0, 1}) {
		return 'v'
	}
	if p.dir.equals(vector{1, 0}) {
		return '>'
	}
	return '<'
}

type model struct {
	foodPos vector
	player  player
}

func (m *model) setRandomFoodPos() {
	m.foodPos = vector{
		x: rand.IntN(20),
		y: rand.IntN(20),
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
		case "k", "up":
			if m.player.dir != dirDown {
				m.player.dir = dirUp
			}
		case "j", "down":
			if m.player.dir != dirUp {
				m.player.dir = dirDown
			}
		case "h", "left":
			if m.player.dir != dirRight {
				m.player.dir = dirLeft
			}
		case "l", "right":
			if m.player.dir != dirLeft {
				m.player.dir = dirRight
			}
		}
	case moveMsg:
		m.player.move(m, m.foodPos)

		head := m.player.body[0]

		if head.x >= 20 || head.x < 0 || head.y < 0 || head.y >= 20 {
			return m, tea.Quit
		}

		for i, b := range m.player.body {
			if i < 2 {
				continue
			}
			if b.equals(head) {
				return m, tea.Quit
			}
		}

		if head.x == m.foodPos.x && head.y == m.foodPos.y {
			m.setRandomFoodPos()
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "----------------------\n"

	for y := 0; y < 20; y++ {
		s += "|"
		for x := 0; x < 20; x++ {
			drew := false
			for i, b := range m.player.body {
				if b.x == x && b.y == y {
					if i == 0 {
						s += string(m.player.headChar())
					} else {
						s += "*"
					}
					drew = true
				}
			}
			if !drew {
				if x == m.foodPos.x && y == m.foodPos.y {
					s += "0"
					drew = true
				}
			}
			if !drew {
				s += " "
			}
		}
		s += "|\n"
	}

	s += "----------------------\n"
	s += fmt.Sprintf("Score: %d\n", len(m.player.body))
	return s
}

func initialModel() tea.Model {
	return model{
		foodPos: vector{
			x: rand.IntN(20),
			y: rand.IntN(20),
		},
		player: player{
			body: []vector{{6, 6}},
			dir:  dirRight,
		},
	}
}

func Run() {
	p := tea.NewProgram(initialModel())

	go func() {
		for {
			p.Send(moveMsg{})
			time.Sleep(200 * time.Millisecond)
		}
	}()

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
