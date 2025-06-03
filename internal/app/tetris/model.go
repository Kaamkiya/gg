package tetris

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Color int

const (
	Black Color = iota
	Blue
	Green
	Orange
	Pink
	Teal
	Purple
)

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9f6f2"))

var colors = map[Color]lipgloss.Style{
	Black:  defaultStyle.Background(lipgloss.Color("#000000")),
	Blue:   defaultStyle.Background(lipgloss.Color("#eee4da")),
	Green:  defaultStyle.Background(lipgloss.Color("#4CA74F")),
	Orange: defaultStyle.Background(lipgloss.Color("#CF6209")),
	Pink:   defaultStyle.Background(lipgloss.Color("#D85B85")),
	Teal:   defaultStyle.Background(lipgloss.Color("#2692E8")),
	Purple: defaultStyle.Background(lipgloss.Color("#9047A3")),
}

type model struct {
	score int
}

func InitialModel() model {
	return model{
		0,
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
		default:
			m.score++
		}
	}

	return m, nil
}

func (m model) View() string {
	gameboard := NewGameboard(colors)
	return gameboard.Render()
}
