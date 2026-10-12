package tron

import (
	"fmt"
	"os"

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
}
