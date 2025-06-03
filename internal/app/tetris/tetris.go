package tetris

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg struct{}

func Run() {
	initialModel := initialModel()
	p := tea.NewProgram(&initialModel)

	gameTicker := time.NewTicker(300 * time.Millisecond)

	go func() {
		for {
			<-gameTicker.C
			p.Send(TickMsg{})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("An error: %v", err)
		os.Exit(1)
	}

	gameTicker.Stop()
	fmt.Println("")
}
