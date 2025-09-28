package main

import (
	"fmt"
	"log"
	"math/rand"
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
	Score int

	// As decimal
	Accuracy float32

	// Time elapsed
	Time int

	// Number of incorrect presses
	Errors int

	// Number of successful presses
	Hits int

	// Stores indexes of characters that have been attempted
	// Ensures no double counting for Hits or Errors
	SeenIdxSet map[int]int
}

// Define your model
type Model struct {
	Cfg *Config 

	// Current string ID in play
	PStrsID int

	PStr string

	// Current character the user is trying to solve for
	// based on PStr
	PIdx int

	PIdxLowerLimit int

	PUnderlines string

	PSlice []string

	// Current word being attempted, initially 0
	WordIdx int

	InStr string

	// Length of user input string
	InLen int

	State State
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
	if m.InLen > 0 {
		m.InStr = m.InStr[:len(m.InStr)-CHAR_LEN] //CHAR_LEN]
		m.InLen--

		// Decrement underline pointer if greater than lower limit
		if m.PIdx > m.PIdxLowerLimit {
			m.PIdx--
			m.PUnderlines = updateUnderlines(m.PIdx, m.PIdx+1, m.PUnderlines)
		}
	}
}

// Types a character from the user input, and updates the underline string
func typeChar(m *Model, in string) {
	lastWordIncr := 1
	space := " "
	if m.WordIdx == len(m.PSlice)-1 {
		lastWordIncr = 2
		space = ""
	}

	if m.InLen < len(m.PSlice[m.WordIdx])+lastWordIncr {
		// Update underline pointer and add colors to characters for output
		if m.PIdx < len(m.PUnderlines)-1 {

			if in == string(m.PStr[m.PIdx]) {
				in = GREEN + in + RESET

				if _, ok := m.State.SeenIdxSet[m.PIdx]; !ok {
					m.State.Hits++
					m.State.SeenIdxSet[m.PIdx] = 1
				}

			} else {
				in = RED + in + RESET
				if _, ok := m.State.SeenIdxSet[m.PIdx]; !ok {
					m.State.Errors++
					m.State.SeenIdxSet[m.PIdx] = 1
				}
			}
			m.InLen++

			m.PIdx++
			m.PUnderlines = updateUnderlines(m.PIdx, m.PIdx-1, m.PUnderlines)

			m.InStr += in
		}
	}

	// User typed correctly
	if removeColors(m.InStr[m.PIdxLowerLimit:]) == (m.PSlice[m.WordIdx] + space) {
		m.WordIdx++
		m.PIdxLowerLimit = m.PIdx
		m.InStr = ""
		for range(m.PIdxLowerLimit) {
			m.InStr += " "
		}
		m.InLen = 0
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
			if m.WordIdx > len(m.PSlice)-1 {
				// Select next prompt
				m.PStrsID = getNewPromptId(m.Cfg, m.PStrsID)

				if m.PStrsID == -2 {
					fmt.Println("Finished!")
					return m, tea.Quit
				}

				// Reinitialize variables
				m.PStr = m.Cfg.Prompts[m.PStrsID-1].Text
				m.PSlice = strings.Split(m.PStr, " ")
				m.PUnderlines = UNDERLINE_CHAR
				for range m.PStr {
					m.PUnderlines += " "
				}
				m.State.Score++
				m.InStr = ""
				m.PIdx = 0
				m.PIdxLowerLimit = 0
				m.WordIdx = 0
				m.InLen = 0

				// Wipe the seen set
				m.State.SeenIdxSet = make(map[int]int)
			}
		}
	case TickMsg:
		m.State.Time++
		return m, doTick()
	}

	return m, nil
}

func getNewPromptId(cfg *Config, curID int) int {
	cfg.SeenIDs[curID] = 1

	if len(cfg.SeenIDs) - 1 >= len(cfg.Prompts) {
		return -2
	}

	newIdx := rand.Intn(len(cfg.Prompts)) + 1 // 1 indexed
	_, ok := cfg.SeenIDs[newIdx]

	for ok {
		fmt.Println("looping...")
		newIdx = rand.Intn(len(cfg.Prompts)) + 1
		_, ok = cfg.SeenIDs[newIdx]
	}


	// For now just increment by one

	return newIdx
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

	PStr := m.PStr

	// Updates the '|' cursor line to the prompt string
	if m.PIdx+1 < len(m.PStr)+1 {
		PStr = m.PStr[:m.PIdx] + "|" + string(m.PStr[m.PIdx]) + m.PStr[m.PIdx+1:]
	}

	m.State.Accuracy = 100

	if m.State.Hits > 0 {
		m.State.Accuracy = (1.0 - (float32(m.State.Errors) / float32(m.State.Hits))) * 100
	}

	return fmt.Sprintf(
		"%s\n%s\n%s\nScore: %d\nTime elapsed (s): %vs\nAccuracy: %.0f%%\n\n", PStr, m.PUnderlines, m.InStr, m.State.Score, m.State.Time, m.State.Accuracy,
	)
}

func main() {
	// Parse 'library.yaml' for a list of prompts
	cfg, err := parseYAML("")
	if err != nil {
		log.Fatal(err)
	}
	cfg.SeenIDs = make(map[int]int)

	//PStr := "In Golang, string replacement is primarily handled by functions within the strings package. The two main functions for this purpose are"
	pStrsID := getNewPromptId(cfg, -1)
	pStr := cfg.Prompts[pStrsID-1].Text
	pSlice := strings.Split(pStr, " ")
	pUnderlines := UNDERLINE_CHAR
	for range pStr {
		pUnderlines += " "
	}

	p := tea.NewProgram(Model{
		Cfg: cfg,
		PStrsID: pStrsID,
		PStr:        pStr,
		PSlice:      pSlice,
		PUnderlines: pUnderlines,
		State:       State{
			SeenIdxSet: make(map[int]int),
			Hits: 1,
		},
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
