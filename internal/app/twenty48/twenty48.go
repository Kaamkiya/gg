package twenty48

import (
	"math/rand/v2"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	// TODO: add a score counter.
	colors map[int]lipgloss.Style
	grid   [4][4]int
}

func initialModel() tea.Model {
	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f9f6f2"))
	c := func(s string) lipgloss.Color {
		return lipgloss.Color(s)
	}

	m := model{
		colors: map[int]lipgloss.Style{
			0:    defaultStyle.Background(c("#3c3a32")),
			2:    defaultStyle.Background(c("#eee4da")).Foreground(c("#000000")),
			4:    defaultStyle.Background(c("#ede0c8")).Foreground(c("#000000")),
			8:    defaultStyle.Background(c("#f2b179")),
			16:   defaultStyle.Background(c("#f59563")),
			32:   defaultStyle.Background(c("#f67c5f")),
			64:   defaultStyle.Background(c("#f65e3b")),
			128:  defaultStyle.Background(c("#edcf72")),
			256:  defaultStyle.Background(c("#edcc61")),
			512:  defaultStyle.Background(c("#edc850")),
			1024: defaultStyle.Background(c("#edc53f")),
			2048: defaultStyle.Background(c("#edc22e")),
		},
		grid: [4][4]int{},
	}

	// The board needs to start with two starting tiles.
	m.AddTile()
	m.AddTile()
	return m
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
		case "left", "h":
			beforeMerge := m.grid
			m.MergeTilesLeft()
			m.ValidateTile(beforeMerge)
		case "down", "j":
			/* Instead of creating a separate method to merge down,
			 * we rotate the grid. This is because the
			 * m.MergeTilesLeft() method is *much* more complex
			 * than m.Rotate90(), so it's simpler to rotate, merge,
			 * then rotate back than to create a separate function.
			 */
			beforeMerge := m.grid
			m.Rotate90(false)
			m.MergeTilesLeft()
			m.Rotate90(true)

			m.ValidateTile(beforeMerge)
		case "up", "k":
			beforeMerge := m.grid
			m.Rotate90(true)
			m.MergeTilesLeft()
			m.Rotate90(false)
			m.ValidateTile(beforeMerge)
		case "right", "l":
			beforeMerge := m.grid
			m.Rotate90(false)
			m.Rotate90(false)
			m.MergeTilesLeft()
			m.Rotate90(true)
			m.Rotate90(true)
			m.ValidateTile(beforeMerge)
		}
	}

	if m.CheckForWin() {
		return m, tea.Quit
	}

	// The game is over when there are no possible merges.
	if !m.CanMove() {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			/* The tiles don't look like this: |  256 |, they look
			 * like this: --------
			 *            |      |
			 *            |  256 |
			 *            |      |
			 *            --------
			 * For that reason, we add empty spaces. It provides a
			 * row of padding, so the game looks better.
			 */
			s += m.colors[m.grid[y][x]].Render("      ")
		}
		s += "\n"
		for x := 0; x < 4; x++ {
			stringifiedNum := strconv.Itoa(m.grid[y][x])
			if stringifiedNum == "0" {
				stringifiedNum = "."
			}

			/* Add spaces before the number so that the width of
			 * the tiles is even.
			 */
			for i := 0; i < 5-len(stringifiedNum); i++ {
				s += m.colors[m.grid[y][x]].Render(" ")
			}
			s += m.colors[m.grid[y][x]].Render(stringifiedNum + " ")
		}
		s += "\n"
		for x := 0; x < 4; x++ {
			// This is for the bottom line of padding.
			s += m.colors[m.grid[y][x]].Render("      ")
		}
		s += "\n"
	}

	s += "\nhjkl or arrows to move"

	return s
}

func (m *model) MergeTilesLeft() {
	for i := range m.grid {
		stopMerge := 0
		for j := 1; j < len(m.grid[i]); j++ {
			if m.grid[i][j] != 0 {
				for k := j; k > stopMerge; k-- {
					switch {
					case m.grid[i][k-1] == 0:
						m.grid[i][k-1] = m.grid[i][k]
						m.grid[i][k] = 0
					case m.grid[i][k-1] == m.grid[i][k]:
						m.grid[i][k-1] += m.grid[i][k]
						m.grid[i][k] = 0
						stopMerge = k
					default:
						break
					}
				}
			}
		}
	}
}

func (m *model) AddTile() bool {
	empty := []int{}
	for y, row := range m.grid {
		for x, cell := range row {
			if cell == 0 {
				empty = append(empty, y*4+x)
			}
		}
	}

	if len(empty) == 0 {
		return false
	}

	rndSrc := rand.NewPCG(
		uint64(time.Now().UnixNano()),
		uint64(time.Now().UnixMilli()),
	)
	rnd := rand.New(rndSrc)

	cell := empty[rnd.IntN(len(empty))]

	if rnd.IntN(10) < 9 {
		m.grid[cell/len(m.grid)][cell%len(m.grid)] = 2
	} else {
		m.grid[cell/len(m.grid)][cell%len(m.grid)] = 4
	}

	return true
}

func (m *model) Rotate90(counterClockWise bool) {
	rotatedGrid := [4][4]int{}
	for i, row := range m.grid {
		rotatedGrid[i] = [4]int{}
		for j := range row {
			if counterClockWise {
				rotatedGrid[i][j] = m.grid[j][len(m.grid)-i-1]
			} else {
				rotatedGrid[i][j] = m.grid[len(m.grid)-j-1][i]
			}
		}
	}
	m.grid = rotatedGrid
}

func (m model) CheckForWin() bool {
	for _, row := range m.grid {
		for x := range row {
			if row[x] == 2048 {
				return true
			}
		}
	}

	return false
}

func (m *model) ValidateTile(beforeMerge [4][4]int) (tea.Model, tea.Cmd) {
	// Check if the grid has changed after handling the merge logic.
	if m.grid != beforeMerge {
		// If unable to add a tile, quit.
		if !m.AddTile() {
			return m, tea.Quit
		}
	}

	return m, nil
}

// Validates that movement is possible. Returns true only if there is at least one empty tile or any adjacent equal pairs present.
func (m model) CanMove() bool {
	// Checks for empty tiles.
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if m.grid[y][x] == 0 {
				return true
			}
		}
	}

	// Checks if there are any horizontal merges within the grid.
	for y := 0; y < 4; y++ {
		for x := 0; x < 3; x++ {
			if m.grid[y][x] == m.grid[y][x+1] {
				return true
			}
		}
	}

	// Checks if there are any vertical merges within the grid.
	for x := 0; x < 4; x++ {
		for y := 0; y < 3; y++ {
			if m.grid[y][x] == m.grid[y+1][x] {
				return true
			}
		}
	}

	return false
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
