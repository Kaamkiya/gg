package main

import (
	"fmt"
	"os"
	"strings"

	//"strings"

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
	CHAR_LEN = len(RED) + 1 + len(RESET)
)

// Define your model
type model struct {
	pStr string
	
	// Current character the user is trying to solve for
	pIdx int

	pIdxLowerLimit int

	pUnderlines string

	pSlice []string

	// Current word being attempted, initially 0
	curIdx int 

	inStr string
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
		case "enter", "ctrl+w":

		case "backspace":
			if len(m.inStr) > 0 {
				m.inStr = m.inStr[:len(m.inStr)- 1]

				if m.pIdx > m.pIdxLowerLimit {
					m.pIdx--
					m.pUnderlines = updateUnderlines(m.pIdx, m.pIdx+1, m.pUnderlines)
				}
			}

		default:
			if len(m.inStr) < len(m.pSlice[m.curIdx]) + 1{

				m.inStr += in
				if m.pIdx < len(m.pUnderlines) - 1 {
					m.pIdx++
					m.pUnderlines = updateUnderlines(m.pIdx, m.pIdx-1, m.pUnderlines)
				}

				// User typed correctly 
				if m.inStr == (m.pSlice[m.curIdx] + " ") {
					m.curIdx++
					m.inStr = ""
					m.pIdxLowerLimit = m.pIdx
				}

				if m.curIdx > len(m.pSlice) - 1 {
					fmt.Println("Winner!")
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

// Takes the new idx to mark where the current pointer should point to in the string
func updateUnderlines(newI, old int, s string) string {
	s = s[:old] + " " + s[old+1:]
	s = s[:newI] + UNDERLINE_CHAR + s[newI+1:]
	return s

}

// View renders the UI
func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s", m.pStr, m.pUnderlines, m.inStr,
	)
}

func main() {
	//Test("hello ")
	pStr := "hello, world!"
	pSlice := strings.Split(pStr, " ")
	pUnderlines := UNDERLINE_CHAR
	for range(pStr) {
		pUnderlines += " "
	}

	p := tea.NewProgram(model{
		pStr: pStr,
		pSlice: pSlice,
		pUnderlines: pUnderlines,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

