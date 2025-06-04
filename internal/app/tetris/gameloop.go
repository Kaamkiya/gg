package tetris

import (
	"strconv"
	"strings"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	"github.com/Kaamkiya/gg/internal/app/tetris/shape"
	tea "github.com/charmbracelet/bubbletea"
)

const gameProgressTickDelay time.Duration = 300 * time.Millisecond

type GameProgressTick struct{}

func initialModel() GameState {
	return GameState{
		nil,
		nil,
		NewGameboard(color.Colors),
		shape.NewRandomizer(),
		0,
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
		} else if !gs.isPaused {
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
			case "p", "P":
				gs.isPaused = true
				return gs, nil
			}
		} else {
			if msg.String() == "p" || msg.String() == "P" {
				gs.isPaused = false
				return gs, tea.Tick(gameProgressTickDelay, func(time.Time) tea.Msg { return GameProgressTick{} })
			}
		}
	case GameProgressTick:
		if gs.isPaused {
			return gs, nil
		}
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
	sidebarLines := make([]string, 14)
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

	scoreStr := strconv.FormatUint(uint64(gs.score), 10)
	sidebarLines[6] = "                      "
	sidebarLines[7] = "   Your score is      "
	sidebarLines[8] = strings.Repeat(" ", 22-len(scoreStr)) + scoreStr
	sidebarLines[9] = "                      "
	sidebarLines[10] = "  hjl/←↓→ to move    "
	sidebarLines[11] = "  z,x to rotate      "
	sidebarLines[12] = "  q/ctl+c to quit    "
	sidebarLines[13] = "  p to pause"

	return sidebarLines
}
