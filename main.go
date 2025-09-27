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
	//RESET = "XXXXXXX"
	RED = "\033[31m"
	//RED = "RRRRRRRR"
	GREEN = "\033[32m"
	//GREEN = "GGGGGGGG"
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
	// based on pStr
	pIdx int

	pIdxLowerLimit int

	pUnderlines string

	pSlice []string

	// Current word being attempted, initially 0
	wordIdx int 

	inStr string

	// Length of user input string
	inLen int
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
			if m.inLen > 0 {
				m.inStr = m.inStr[:len(m.inStr) - CHAR_LEN ]//CHAR_LEN]
				m.inLen--

				// Decrement underline pointer if greater than lower limit
				if m.pIdx > m.pIdxLowerLimit {
					m.pIdx--
					m.pUnderlines = updateUnderlines(m.pIdx, m.pIdx+1, m.pUnderlines)
				}
			}

		default:
			if m.inLen < len(m.pSlice[m.wordIdx]) + 1{

				// Update underline pointer and add colors to characters for output
				if m.pIdx < len(m.pUnderlines) - 1 {

					if in == string(m.pStr[m.pIdx]) {
						in = GREEN + in + RESET
					} else { 
					in = RED + in + RESET 
					}
					m.inLen ++

					m.pIdx++
					m.pUnderlines = updateUnderlines(m.pIdx, m.pIdx-1, m.pUnderlines)
				}
					m.inStr += in

				// User typed correctly 
				if removeColors(m.inStr) == (m.pSlice[m.wordIdx] + " ") {

					m.wordIdx++
					m.inStr = ""
					m.pIdxLowerLimit = m.pIdx
					m.inLen = 0
				}

				// Exit if finished
				if m.wordIdx > len(m.pSlice) - 1 {
					fmt.Println("Winner!")
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

// Removes ANSI colors from a string
func removeColors(s string) string {
	s = strings.ReplaceAll(s, RESET, "")
	s = strings.ReplaceAll(s, GREEN, "")
	s = strings.ReplaceAll(s, RED, "")
	return s
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
	pStr := "In Golang, string replacement is primarily handled by functions within the strings package. The two main functions for this purpose are"
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

