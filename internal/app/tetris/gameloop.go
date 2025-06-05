package tetris

import (
	"strconv"
	"strings"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	"github.com/Kaamkiya/gg/internal/app/tetris/shape"
	tea "github.com/charmbracelet/bubbletea"
)

type gameProgressTick struct{}

func initialModel() gameState {
	return gameState{
		nil,
		nil,
		newGameboard(color.Colors),
		shape.NewRandomizer(),
		0,
		&difficulty{
			initialDifficulyCountDown,
			initialDifficulyLevel,
			initialGameProgressTickDelay,
		},
		false,
	}
}

func (gs *gameState) Init() tea.Cmd {
	return func() tea.Msg {
		return gameProgressTick{}
	}
}

// Update implements the game loop by handling the tea.Msg structs. There are the following flows:
//   - Base loop: gameProgressTick -> handleGameProgress -> gameProgressTick
//   - Line complete: gameProgressTick -> handleGameProgress -> lineAnimationTick
//   - Line animation ongoing: lineAnimationTick -> handleLineAnimation -> lineAnimationTick
//   - Line animation finished: lineAnimationTick -> handleLineAnimation -> gameProgressTick
func (gs *gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "Q" {
			return gs, tea.Quit
		} else if !gs.isPaused {
			switch msg.String() {
			case "h", "H", "left":
				gs.handleLeft()
			case "l", "L", "right":
				gs.handleRight()
			case "j", "J", "down":
				gs.handleDown()
			case "z", "Z":
				gs.handleLeftRotate()
			case "x", "X":
				gs.handleRightRotate()
			case "p", "P":
				gs.isPaused = true
				return gs, nil
			}
		} else {
			if msg.String() == "p" || msg.String() == "P" {
				gs.isPaused = false
				return gs, tea.Tick(gs.currentDifficulty.gameProgressTickDelay, func(time.Time) tea.Msg { return gameProgressTick{} })
			}
		}
	case gameProgressTick:
		if gs.isPaused {
			return gs, nil
		}
		return gs, gs.handleGameProgressTick()
	case lineAnimationTick:
		return gs, gs.handleLineAnimationTick(msg)
	}

	return gs, nil
}

// View method creates the view by generating the play area and the sidebar. Although the Tetris board size is
// defined by Height and Width, the play area is larger. Each Tetris box is 4 characters wide and 2 characters tall
// so the total play area size is 2 * Height * 4 * Width characters. On each line of the play area, a sidebar
// line is appended.
func (gs *gameState) View() string {
	boardBuilder := strings.Builder{}
	boardBuilder.Grow(height*width*5 + 22*13)

	sideBarLines := buildSidebar(gs)

	for i := range height {
		lineBuilder := strings.Builder{}
		lineBuilder.Grow(width * 2)

		for j := range width {
			nextChar := gs.gameBoard.Colors[gs.gameBoard.Grid[i][j]].Render("    ")
			lineBuilder.WriteString(nextChar)
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

func buildSidebar(gs *gameState) []string {
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
	sidebarLines[13] = "  p to pause         "

	return sidebarLines
}
