package tron

import (
	"fmt"
	"os"

	"github.com/Kaamkiya/gg/internal/app/tron/gameboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
)

type gameConfig struct {
	playerConfigs []playerConfig
}

type playerConfig struct {
	color gameboard.Color
}

func newGameConfig(numOfPlayers int) gameConfig {
	playerConfigs := []playerConfig{{gameboard.Red}, {gameboard.Blue}}

	if numOfPlayers == 3 {
		playerConfigs = append(playerConfigs, playerConfig{gameboard.Purple})
	}

	return gameConfig{
		playerConfigs,
	}
}

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

	var playerCount int

	err = huh.NewSelect[int]().
		Title("How many players?").
		Options(
			huh.NewOption("2 players", 2),
			huh.NewOption("3 players", 3),
		).
		Value(&playerCount).
		Run()

	if err != nil {
		fmt.Println("Error: failed to get number of players.")
		panic(err)
	}

	playerConfigs := []playerConfig{
		{gameboard.Blue},
		{gameboard.Red},
	}

	if playerCount == 3 {
		playerConfigs = append(playerConfigs, playerConfig{gameboard.Purple})
	}
	initialModel := newGameState(gameConfig{playerConfigs: playerConfigs})

	p := tea.NewProgram(initialModel)

	if _, err := p.Run(); err != nil {
		fmt.Printf("An error: %v", err)
		os.Exit(1)
	}

	fmt.Println("")

}
