package dodger

import (
	"fmt"
	"math/rand/v2"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type spawnBlockMsg struct{}
type moveBlockMsg struct{}

type vector struct {
	x int
	y int
}

type model struct {
	size   vector   // The size of the screen.
	player vector   // The position of the player.
	blocks []vector // The positions of each block on the screen.
	score  int      // The amount of blocks that have gone off-screen.

	blockStyle  lipgloss.Style
	playerStyle lipgloss.Style
}

func initialModel() tea.Model {
	size := vector{30, 20}
	return model{
		size:        size,
		player:      vector{int(size.x / 2), size.y - 1},
		blocks:      []vector{},
		score:       0,
		blockStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc")),
		playerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaff")),
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
		case "left", "a":
			m.player.x--
			if m.player.x < 0 {
				m.player.x = m.size.x - 1
			}
		case "right", "d":
			m.player.x++
			if m.player.x >= m.size.x {
				m.player.x = 0
			}
		}
	case spawnBlockMsg:
		m.blocks = append(m.blocks, vector{rand.IntN(m.size.x), 0})
	case moveBlockMsg:
		m.moveBlocks()
	}

	for _, b := range m.blocks {
		if b.x == m.player.x && b.y == m.player.y {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("\nScore: %d\n", m.score)

	for y := 0; y < m.size.y; y++ {
		for x := 0; x < m.size.x; x++ {
			drew := false
			for _, b := range m.blocks {
				if b.x == x && b.y == y {
					s += m.blockStyle.Render(string(rune(0x2022))) // 0x2022 is a unicode bullet point.
					drew = true
				}
			}
			if !drew {
				if x == m.player.x && y == m.player.y {
					s += m.playerStyle.Render(string(rune(0x2205))) // 0x2205 is a unicode rectangle.
				} else {
					s += " "
				}
			}
		}
		s += "\n"
	}

	s += "wasd or arrows to move"

	return s
}

func (m *model) moveBlocks() {
	for i := range m.blocks {
		m.blocks[i].y++
	}

	for i := range m.blocks {
		if m.blocks[i].y > m.size.y {
			m.blocks = append(m.blocks[:i], m.blocks[i+1:]...)
			m.score++
			// There can only be one block to be removed every time.
			break
		}
	}
}

func Run() {
	prog := tea.NewProgram(initialModel())

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			prog.Send(spawnBlockMsg{})
			prog.Send(moveBlockMsg{})
		}
	}()

	if _, err := prog.Run(); err != nil {
		panic(err)
	}
}
