package color

import "github.com/charmbracelet/lipgloss"

type Color int

const (
	None Color = iota
	Blue
	Green
	Orange
	Pink
	Teal
	Purple
	Magenta
	Beige
)

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9f6f2"))

var Colors = map[Color]lipgloss.Style{
	None:    defaultStyle.Background(lipgloss.NoColor{}),
	Blue:    defaultStyle.Background(lipgloss.Color("#063970")),
	Green:   defaultStyle.Background(lipgloss.Color("#4CA74F")),
	Orange:  defaultStyle.Background(lipgloss.Color("#CF6209")),
	Pink:    defaultStyle.Background(lipgloss.Color("#D85B85")),
	Teal:    defaultStyle.Background(lipgloss.Color("#2692E8")),
	Purple:  defaultStyle.Background(lipgloss.Color("#9047A3")),
	Magenta: defaultStyle.Background(lipgloss.Color("#CA1F7B")),
	Beige:   defaultStyle.Background(lipgloss.Color("#FFFDD0")),
}
