package tetris

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	initialModel := initialModel()
	p := tea.NewProgram(&initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("An error: %v", err)
		os.Exit(1)
	}

	fmt.Println("")
}
