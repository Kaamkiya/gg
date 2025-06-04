package tetris

import (
	"strings"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	tea "github.com/charmbracelet/bubbletea"
)

const gameProgressTickDelay time.Duration = 300 * time.Millisecond

type GameProgressTick struct{}

func initialModel() GameState {
	return GameState{
		nil,
		nil,
		NewGameboard(color.Colors),
		false,
	}
}

func (gs *GameState) Init() tea.Cmd {
	return func() tea.Msg {
		return GameProgressTick{}
	}
}

func (gs *GameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "Q" {
			return gs, tea.Quit
		} else if !gs.isAnimating {
			switch msg.String() {
			case "h", "H", "left":
				gs.HandleLeft()
			case "l", "L", "right":
				gs.HandleRight()
			case "j", "J", "down":
				gs.HandleDown()
			case "z", "Z":
				gs.HandleLeftRotate()
			case "x", "X":
				gs.HandleRightRotate()
			}
		}
	case GameProgressTick:
		return gs, gs.HandleGameProgressTick()
	case LineAnimationTick:
		return gs, gs.handleLineAnimationTick(msg)
	}

	return gs, nil
}

func (gs *GameState) View() string {
	boardBuilder := strings.Builder{}
	boardBuilder.Grow(Height*Width*5 + 22*13)

	sideBarLines := buildSidebar(gs)

	for i := range Height {
		lineBuilder := strings.Builder{}
		lineBuilder.Grow(Width * 2)

		for j := range Width {
			nextChar := gs.gameBoard.Colors[gs.gameBoard.Grid[i][j]].Render("  ")
			lineBuilder.WriteString(nextChar + nextChar)
		}

		line := lineBuilder.String()
		boardBuilder.WriteString(line)

		if 2*i < len(sideBarLines) {
			boardBuilder.WriteString(sideBarLines[2*i])
		}
		boardBuilder.WriteString("\n")

		boardBuilder.WriteString(line)
		if 2*i+1 < len(sideBarLines) {
			boardBuilder.WriteString(sideBarLines[2*i+1])
		}
		boardBuilder.WriteString("\n")
	}

	return boardBuilder.String()
}

func buildSidebar(gs *GameState) []string {
	sidebarLines := make([]string, 13)
	sidebarLines[0] = "      Next Shape      "
	sidebarLines[1] = "                      "

	if gs.nextShape != nil {
		grid := gs.nextShape.GetGrid()

		for i := range 4 {
			if i >= len(grid) {
				sidebarLines[i+2] = "                      "
			} else {
				lineBuilder := strings.Builder{}
				spaceLength := (22 - len(grid[i])) / 2
				lineBuilder.WriteString(strings.Repeat(" ", spaceLength))

				for j := range grid[i] {
					if grid[i][j] {
						lineBuilder.WriteString(gs.gameBoard.Colors[gs.nextShape.GetColor()].Render(" "))
					} else {
						lineBuilder.WriteString(" ")
					}
				}
				lineBuilder.WriteString(strings.Repeat(" ", spaceLength))

				sidebarLines[i+2] = lineBuilder.String()
			}
		}
	}

	sidebarLines[6] = "                      "
	sidebarLines[7] = "   Your score is      "
	sidebarLines[8] = "   1111               "
	sidebarLines[9] = "                      "
	sidebarLines[10] = "  hjl or arrows to    "
	sidebarLines[11] = "  move, z,x to rotate "
	sidebarLines[12] = "  q or ctl+c to quit  "

	return sidebarLines
}
