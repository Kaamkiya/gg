package main

import (
	"fmt"

	"github.com/Kaamkiya/gg/internal/app/connect4"
	"github.com/Kaamkiya/gg/internal/app/dodger"
	"github.com/Kaamkiya/gg/internal/app/hangman"
	"github.com/Kaamkiya/gg/internal/app/maze"
	"github.com/Kaamkiya/gg/internal/app/pong"
	"github.com/Kaamkiya/gg/internal/app/tictactoe"
	"github.com/Kaamkiya/gg/internal/app/twenty48"

	"github.com/charmbracelet/huh"
)

func main() {
	var game string

	fmt.Println("gg - a tui for small offline games")

	err := huh.NewSelect[string]().
		Title("choose a game:").
		Options(
			huh.NewOption("2048", "twenty48"),
			huh.NewOption("dodger", "dodger"),
			huh.NewOption("maze", "maze"),
			huh.NewOption("hangman", "hangman"),
			huh.NewOption("connect 4 (2 player)", "connect4"),
			huh.NewOption("pong (2 player)", "pong"),
			huh.NewOption("tictactoe (2 player)", "tictactoe"),
		).
		Value(&game).
		Run()

	if err != nil {
		fmt.Println("Error: failed to run selection menu.")
		panic(err)
	}

	switch game {
	case "maze":
		maze.Run()
	case "pong":
		pong.Run()
	case "tictactoe":
		tictactoe.Run()
	case "dodger":
		dodger.Run()
	case "hangman":
		hangman.Run()
	case "twenty48":
		twenty48.Run()
	case "connect4":
		connect4.Run()
	default:
		panic("This game either doesn't exist or hasn't been implemented.")
	}
}
