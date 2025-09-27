package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	RESET = "\033[0m"
	RED = "\033[31m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"
	BLUE = "\033[34m"
	RESET_LEN = len(RESET)
	COLOR_LEN = len(RED)
	UNDERLINE_CHAR = "^"
)

// Define your model
type model struct {
	prompt string
	promptUnderline string
	input string
	i int
}

// Init runs when the program starts
func (m model) Init() tea.Cmd {
	// No initial command
	return nil
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		in := msg.String()
		switch in {
		case "q", "ctrl+c":
			// Quit on q or ctrl+c
			return m, tea.Quit
		case "enter":

		case "ctrl+w":
			// Delete block from input
			if len(m.input) > 1 {
				m.input = strings.TrimRight(m.input, " ")
				

				for i := len(m.input) - 1; i >= 0; i-- {
					char := m.input[i]

					if char == ' ' {
						m.input = m.input[:i+1]
						break
					} else if i == 1 {
						m.input = ""
						break
					}
				}


				// Rollback currently attempted word
				promptCopy := m.prompt[:m.i]
				strings.TrimRight(promptCopy, " ")

				res := 0
				for i := len(promptCopy) - 1; i >= 0; i-- {
					char := promptCopy[i]

					if char == ' ' {
						res = i+1
						break
					} else if i == 1 {
						break
					}
				}

				if m.i < len(m.promptUnderline) {
					m.promptUnderline = m.promptUnderline[:m.i] + " " + m.promptUnderline[m.i+1:]
				}

				m.i = res

				if m.i < len(m.promptUnderline) {
					m.promptUnderline = m.promptUnderline[:m.i] + UNDERLINE_CHAR + m.promptUnderline[m.i+1:]
				}
			}


		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]

				if m.i < len(m.promptUnderline) {
					m.promptUnderline = m.promptUnderline[:m.i] + " " + m.promptUnderline[m.i+1:]
				}

				if m.i > 0 {
					m.i--
				}

				if m.i < len(m.promptUnderline) {
					m.promptUnderline = m.promptUnderline[:m.i] + UNDERLINE_CHAR + m.promptUnderline[m.i+1:]
				}

			}

		default:
			if len(m.input) < len(m.prompt) {
				m.input += in

				if m.i < len(m.promptUnderline) {
					m.promptUnderline = m.promptUnderline[:m.i] + " " + m.promptUnderline[m.i+1:]
				}

				if m.i < len(m.promptUnderline) {
					m.i++
				}

				if m.i < len(m.promptUnderline) {
					m.promptUnderline = m.promptUnderline[:m.i] + UNDERLINE_CHAR + m.promptUnderline[m.i+1:]
				}
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s", m.prompt, m.promptUnderline, m.input,
	)
}

func main() {
	//Test("hello ")
	newModel := model{prompt: "Type this text out in full"}
	newUnderline := UNDERLINE_CHAR
	for range(newModel.prompt[1:]) {
		newUnderline += " "
	}

	newModel.promptUnderline = newUnderline

	p := tea.NewProgram(newModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

