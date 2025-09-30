package pong

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type vector struct {
	x int
	y int
}

type ballBody struct {
	pos vector
	vel vector
}

type moveBallMsg struct{}

type model struct {
	hitCount int

	size vector

	paddle1 vector
	paddle2 vector

	ball ballBody

	colors []lipgloss.Style
}

func initialModel() tea.Model {
	size := vector{30, 15}

	return model{
		hitCount: 0,
		size:     size,
		paddle1:  vector{1, 8},
		paddle2:  vector{size.x - 1, 7},
		ball: ballBody{
			pos: vector{int(15), int(8)},
			vel: vector{1, 1},
		},
		colors: []lipgloss.Style{
			lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaff")),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaaaa")),
		},
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
		case "a":
			m.MovePaddle(1, -1)
		case "d":
			m.MovePaddle(1, 1)
		case "left":
			m.MovePaddle(2, -1)
		case "right":
			m.MovePaddle(2, 1)
		}
	case moveBallMsg:
		if m.ball.pos.y < 0 || m.ball.pos.y >= m.size.y {
			m.ball.vel.y *= -1
		}

		if m.ball.pos.x < 0 || m.ball.pos.x >= m.size.x {
			m.ball.vel.x *= -1
		}

		if m.ball.pos == m.paddle1 || m.ball.pos == m.paddle2 {
			m.ball.vel.x *= -1
			m.hitCount++
		}

		if m.ball.pos.x == 0 || m.ball.pos.x >= m.size.x {
			return m, tea.Quit
		}

		m.ball.pos.x += m.ball.vel.x
		m.ball.pos.y += m.ball.vel.y
	}
	return m, nil
}

func (m model) View() string {
	s := ""

	for i := 0; i < m.size.x; i++ {
		s += m.colors[i%2].Render(string(rune(9608)))

		for j := 0; j < m.size.y; j++ {
			switch (vector{i, j}) {
			case m.ball.pos:
				s += "o"
			case m.paddle1:
				s += "-"
			case m.paddle2:
				s += "-"
			default:
				s += " "
			}
		}

		s += m.colors[i%2].Render(string(rune(9608)))
		s += "\n"
	}

	s += fmt.Sprintf("\nHit count: %d\n", m.hitCount)

	return s
}

func (m *model) MovePaddle(num, amount int) {
	if num == 1 {
		m.paddle1.y += amount
		if m.paddle1.y < 0 || m.paddle1.y >= m.size.y {
			m.paddle1.y -= amount
		}
	} else {
		m.paddle2.y += amount
		if m.paddle2.y < 0 || m.paddle2.y >= m.size.y {
			m.paddle2.y -= amount
		}
	}
}

func Run() {
	p := tea.NewProgram(initialModel())

	go func() {
		for {
			time.Sleep(300 * time.Millisecond)
			p.Send(moveBallMsg{})
		}
	}()

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
