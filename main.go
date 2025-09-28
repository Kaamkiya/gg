package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	YELLOW         = "\033[33m"
	BLUE           = "\033[34m"
	RESET_LEN      = len(RESET)
	COLOR_LEN      = len(RED)
	UNDERLINE_CHAR = "^"
	CHAR_LEN       = len(RED) + 1 + len(RESET)
)

type TickMsg time.Time

type State struct {
	score int

	// As decimal
	accuracy float32

	// Time elapsed
	time int

	// Number of incorrect presses
	errors int

	// Number of successful presses
	hits int

	// Stores characters that have been attempted
	// Ensures no double counting for hits or errors
	seenSet map[int]int
}

// Define your model
type Model struct {
	cfg *Config 

	// Current string ID in play
	pStrsID int

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

	state State
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init runs when the program starts
func (m Model) Init() tea.Cmd {
	// No initial command
	return doTick()
}

// Deletes a character from user input and updates the underline string
func backspace(m *Model) {
	if m.inLen > 0 {
		m.inStr = m.inStr[:len(m.inStr)-CHAR_LEN] //CHAR_LEN]
		m.inLen--

		// Decrement underline pointer if greater than lower limit
		if m.pIdx > m.pIdxLowerLimit {
			m.pIdx--
			m.pUnderlines = updateUnderlines(m.pIdx, m.pIdx+1, m.pUnderlines)
		}
	}
}

// Types a character from the user input, and updates the underline string
func typeChar(m *Model, in string) {
	lastWordIncr := 1
	space := " "
	if m.wordIdx == len(m.pSlice)-1 {
		lastWordIncr = 2
		space = ""
	}

	if m.inLen < len(m.pSlice[m.wordIdx])+lastWordIncr {
		// Update underline pointer and add colors to characters for output
		if m.pIdx < len(m.pUnderlines)-1 {

			if in == string(m.pStr[m.pIdx]) {
				in = GREEN + in + RESET

				if _, ok := m.state.seenSet[m.pIdx]; !ok {
					m.state.hits++
					m.state.seenSet[m.pIdx] = 1
				}

			} else {
				in = RED + in + RESET
				if _, ok := m.state.seenSet[m.pIdx]; !ok {
					m.state.errors++
					m.state.seenSet[m.pIdx] = 1
				}
			}
			m.inLen++

			m.pIdx++
			m.pUnderlines = updateUnderlines(m.pIdx, m.pIdx-1, m.pUnderlines)

			m.inStr += in
		}
	}

	// User typed correctly
	if removeColors(m.inStr) == (m.pSlice[m.wordIdx] + space) {
		m.wordIdx++
		m.inStr = ""
		m.pIdxLowerLimit = m.pIdx
		m.inLen = 0
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		in := msg.String()

		switch in {
		case "ctrl+c":
			// Quit on q or ctrl+c
			return m, tea.Quit
		case "enter", "ctrl+w":

		case "backspace":
			backspace(&m)

		default:
			typeChar(&m, in)

			// Exit if finished
			if m.wordIdx > len(m.pSlice)-1 {
				// Select next prompt
				m.pStrsID = setNewPrompt(m.cfg, m.pStrsID)

				// Reinitialize variables
				m.pStr = m.cfg.Prompts[m.pStrsID-1].Text
				m.pSlice = strings.Split(m.pStr, " ")
				m.pUnderlines = UNDERLINE_CHAR
				for range m.pStr {
					m.pUnderlines += " "
				}
				m.state.score++
				m.inStr = ""
				m.pIdx = 0
				m.pIdxLowerLimit = 0
				m.wordIdx = 0
				m.inLen = 0

				// Wipe the seen set
				m.state.seenSet = make(map[int]int)
				//return m, tea.Quit
			}
		}
	case TickMsg:
		m.state.time++
		return m, doTick()
	}

	return m, nil
}

func setNewPrompt(cfg *Config, curID int) int {
	// For now just increment by one
	return curID + 1
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
func (m Model) View() string {

	pStr := m.pStr

	// Updates the '|' cursor line to the prompt string
	if m.pIdx+1 < len(m.pStr)+1 {
		pStr = m.pStr[:m.pIdx] + "|" + string(m.pStr[m.pIdx]) + m.pStr[m.pIdx+1:]
	}

	m.state.accuracy = 100

	if m.state.hits > 0 {
		m.state.accuracy = (1.0 - (float32(m.state.errors) / float32(m.state.hits))) * 100
	}

	return fmt.Sprintf(
		"%s\n%s\n%s\nScore: %d\nTime elapsed (s): %vs\nAccuracy: %.0f%%\n\n", pStr, m.pUnderlines, m.inStr, m.state.score, m.state.time, m.state.accuracy,
	)
}

func main() {
	// Parse 'library.yaml' for a list of prompts
	cfg, err := parseYAML("")
	if err != nil {
		log.Fatal(err)
	}

	//pStr := "In Golang, string replacement is primarily handled by functions within the strings package. The two main functions for this purpose are"
	pStrsID := 1
	pStr := cfg.Prompts[pStrsID-1].Text
	pSlice := strings.Split(pStr, " ")
	pUnderlines := UNDERLINE_CHAR
	for range pStr {
		pUnderlines += " "
	}

	p := tea.NewProgram(Model{
		cfg: cfg,
		pStrsID: pStrsID,
		pStr:        pStr,
		pSlice:      pSlice,
		pUnderlines: pUnderlines,
		state:       State{
			seenSet: make(map[int]int),
			hits: 1,
		},
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
