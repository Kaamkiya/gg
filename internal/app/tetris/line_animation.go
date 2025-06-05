package tetris

import (
	"maps"
	"slices"
	"time"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
	tea "github.com/charmbracelet/bubbletea"
)

// lineAnimationInterval is the animation refresh interval
const lineAnimationInterval time.Duration = 100 * time.Millisecond

// lineAnimationTick is a tea.Msg that contains the lines to change in a map
// where the key is the line index and the value is the colors to apply.
// Additionally, it holds how many animations (color changes) are left for the
// animation to complete.
type lineAnimationTick struct {
	linesToUpdate      map[int][width]color.Color
	animationCountDown int
}

func (gs *gameState) constructLineAnimationMsg(completedLines []int) lineAnimationTick {
	completedLineMap := make(map[int][width]color.Color, len(completedLines))
	animationCountdown := 2

	if len(completedLines) == 3 {
		animationCountdown = 4
	}

	if len(completedLines) == 4 {
		animationCountdown = 6
	}

	highlightedLine := [width]color.Color{}
	for i := range width {
		highlightedLine[i] = color.Beige
	}

	for _, v := range completedLines {
		completedLineMap[v] = highlightedLine

	}

	return lineAnimationTick{
		completedLineMap,
		animationCountdown,
	}
}

// handleLineAnimationTick performs the grid updates for the flashing animation when
// lines are completed. If the animation is complete (animationCountDown set to 0) it
// resumes the game. Otherwse it swaps the lines color and continues with the animation.
func (gs *gameState) handleLineAnimationTick(animationTick lineAnimationTick) tea.Cmd {
	if animationTick.animationCountDown == 0 {
		gs.removeCompletedLines(slices.Collect(maps.Keys(animationTick.linesToUpdate)))
		return func() tea.Msg {
			return gameProgressTick{}
		}
	}

	animationTick.animationCountDown--
	newLinesToUpdateMap := make(map[int][width]color.Color, len(animationTick.linesToUpdate))
	for k, v := range animationTick.linesToUpdate {
		newLinesToUpdateMap[k] = gs.gameBoard.Grid[k]
		gs.gameBoard.Grid[k] = v
	}

	return tea.Tick(lineAnimationInterval, func(time.Time) tea.Msg {
		return lineAnimationTick{
			newLinesToUpdateMap,
			animationTick.animationCountDown,
		}
	})
}

