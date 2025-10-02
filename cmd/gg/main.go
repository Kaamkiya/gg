package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Kaamkiya/gg/internal/app/blackjack"
	"github.com/Kaamkiya/gg/internal/app/connect4"
	"github.com/Kaamkiya/gg/internal/app/dodger"
	"github.com/Kaamkiya/gg/internal/app/hangman"
	"github.com/Kaamkiya/gg/internal/app/maze"
	"github.com/Kaamkiya/gg/internal/app/pong"
	"github.com/Kaamkiya/gg/internal/app/snake"
	"github.com/Kaamkiya/gg/internal/app/solitaire"
	"github.com/Kaamkiya/gg/internal/app/sudoku"
	"github.com/Kaamkiya/gg/internal/app/tetris"
	"github.com/Kaamkiya/gg/internal/app/tictactoe"
	"github.com/Kaamkiya/gg/internal/app/twenty48"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Game represents a game option in the menu
type Game struct {
	name        string
	description string
	icon        string
	value       string
	category    string
}

// Model represents the main menu state
type Model struct {
	games         []Game
	cursor        int
	selected      string
	quitting      bool
	width         int
	height        int
	categories    map[string][]Game
	categoryOrder []string
	visualGames   []Game // Games in visual order (category by category)
}

// InitialModel creates the initial state
func InitialModel() Model {
	games := []Game{
		{"Blackjack", "Classic card game against the dealer", "🃏", "blackjack", "Card Games"},
		{"Solitaire", "Klondike solitaire with 1-card draw", "🂱", "solitaire", "Card Games"},
		{"2048", "Slide and merge numbers to reach 2048", "🔢", "twenty48", "Puzzle Games"},
		{"Sudoku", "Fill the 9x9 grid with numbers 1-9", "🧩", "sudoku", "Puzzle Games"},
		{"Dodger", "Avoid obstacles in this arcade game", "🎯", "dodger", "Arcade Games"},
		{"Maze", "Navigate through randomly generated mazes", "🌐", "maze", "Puzzle Games"},
		{"Hangman", "Guess the word before the hangman is drawn", "🎭", "hangman", "Word Games"},
		{"Snake", "Classic snake game - eat food and grow", "🐍", "snake", "Arcade Games"},
		{"Tetris", "Stack falling blocks to clear lines", "🧱", "tetris", "Arcade Games"},
		{"Connect 4", "Two-player strategy game", "🔴", "connect4", "Strategy Games"},
		{"Pong", "Classic two-player paddle game", "🏓", "pong", "Arcade Games"},
		{"Tic Tac Toe", "Two-player classic game", "⭕", "tictactoe", "Strategy Games"},
		{"Tic Tac Toe vs AI", "Play against an AI opponent", "🤖", "tictactoe-ai", "Strategy Games"},
	}

	// Group games by category
	categories := make(map[string][]Game)
	for _, game := range games {
		categories[game.category] = append(categories[game.category], game)
	}

	categoryOrder := []string{"Card Games", "Puzzle Games", "Arcade Games", "Strategy Games", "Word Games"}

	// Create visual games list (games in display order)
	var visualGames []Game
	for _, category := range categoryOrder {
		if games, exists := categories[category]; exists {
			visualGames = append(visualGames, games...)
		}
	}

	return Model{
		games:         games,
		cursor:        0,
		selected:      "",
		quitting:      false,
		categories:    categories,
		categoryOrder: categoryOrder,
		visualGames:   visualGames,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// getCurrentGame returns the currently selected game
func (m Model) getCurrentGame() Game {
	return m.visualGames[m.cursor]
}

// getTotalGames returns the total number of games
func (m Model) getTotalGames() int {
	return len(m.visualGames)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < m.getTotalGames()-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = m.getCurrentGame().value
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Background(lipgloss.Color("#2D3748")).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 2).
		Margin(1, 0)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A0AEC0")).
		Align(lipgloss.Center).
		Margin(0, 0, 2, 0)

	gameItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E2E8F0")).
		Padding(0, 1).
		Margin(0, 1)

	selectedGameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2D3748")).
		Background(lipgloss.Color("#68D391")).
		Bold(true).
		Padding(0, 1).
		Margin(0, 1)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A0AEC0")).
		Italic(true).
		Margin(0, 0, 0, 2)

	categoryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4A5568")).
		Bold(true).
		Margin(1, 0, 0, 0)

	instructionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#718096")).
		Align(lipgloss.Center).
		Margin(2, 0)

	// Build the UI
	var content strings.Builder

	// Title
	content.WriteString(titleStyle.Render("🎮 GG - Game Collection"))
	content.WriteString("\n")
	content.WriteString(subtitleStyle.Render("Choose your game and start playing!"))
	content.WriteString("\n\n")

	// Render games by category
	visualIndex := 0
	for _, category := range m.categoryOrder {
		if games, exists := m.categories[category]; exists {
			content.WriteString(categoryStyle.Render("  ┌─ " + category))
			content.WriteString("\n")

			for _, game := range games {
				style := gameItemStyle
				if visualIndex == m.cursor {
					style = selectedGameStyle
				}

				gameLine := fmt.Sprintf("│%s %s", game.icon, game.name)
				content.WriteString(style.Render(gameLine))
				content.WriteString("\n")

				// Add description for selected game
				if visualIndex == m.cursor {
					content.WriteString(descriptionStyle.Render("│  " + game.description))
					content.WriteString("\n")
				}

				visualIndex++
			}

			content.WriteString("  └")
			content.WriteString(strings.Repeat("─", 50))
			content.WriteString("\n")
		}
	}

	// Instructions
	content.WriteString("\n")
	content.WriteString(instructionsStyle.Render("↑/↓: Navigate  •  Enter: Select  •  q: Quit"))

	return content.String()
}

func main() {
	// Initialize the program
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())

	// Run the program
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Get the final model and selected game
	if model, ok := m.(Model); ok && model.selected != "" {
		// Run the selected game
		switch model.selected {
		case "blackjack":
			blackjack.Run()
		case "solitaire":
			solitaire.Run()
		case "maze":
			maze.Run()
		case "pong":
			pong.Run()
		case "tictactoe":
			tictactoe.Run()
		case "tictactoe-ai":
			tictactoe.RunVsAi()
		case "dodger":
			dodger.Run()
		case "hangman":
			hangman.Run()
		case "twenty48":
			twenty48.Run()
		case "connect4":
			connect4.Run()
		case "snake":
			snake.Run()
		case "sudoku":
			sudoku.Run()
		case "tetris":
			tetris.Run()
		default:
			fmt.Println("This game either doesn't exist or hasn't been implemented.")
		}
	}
}
