package tron

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

func Run() {
	width, height, err := term.GetSize(os.Stdout.Fd())

	if err != nil {
		fmt.Println("Error while getting window dimensions: ", err)
		os.Exit(1)
	}

	if height < 50 || width < 200 {
		fmt.Println("Tron requires a terminal with a size of at least 200x50, please resize accordingly.")
		os.Exit(0)
	}

	initialModel := newGameState()
	p := tea.NewProgram(&initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("An error: %v", err)
		os.Exit(1)
	}

	fmt.Println("")

}
