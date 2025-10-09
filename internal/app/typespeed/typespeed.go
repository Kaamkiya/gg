package typespeed

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	ORANGE = "\033[33m" // closest standard color to orange
	TEAL   = "\033[36m" // cyan/teal
	BROWN  = "\033[33m" // same as orange in basic ANSI
	PURPLE = "\033[35m" // magenta/purple

	RESET_LEN      = len(RESET)
	COLOR_LEN      = len(RED)
	UNDERLINE_CHAR = "^"
	CHAR_LEN       = len(RED) + 1 + len(RESET)
)

type TickMsg time.Time

type State struct {
  // Number of prompts completed
	PromptCompletions int

  // Number of words completed
  WordCompletions int

	// As decimal
	Accuracy float32

  // Correct words typed/minute
  WPM float32

  //
  CPM float32

	// Time elapsed in seconds
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
	PromptStrsID int

	PromptStr string

	// Current character the user is trying to solve for
	// based on PromptStr
	PromptIdx int

	PromptIdxLowerLimit int

	PromptUnderlines string

	PromptSlice []string

	// Current word being attempted, initially 0
	WordIdx int

	InputStr string

	// Length of user input string
	InputLen int

	State *State
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
	if m.InputLen > 0 {
		m.InputStr = m.InputStr[:len(m.InputStr)-CHAR_LEN] //CHAR_LEN]
		m.InputLen--

		// Decrement underline pointer if greater than lower limit
		if m.PromptIdx > m.PromptIdxLowerLimit {
			m.PromptIdx--
			m.PromptUnderlines = updateUnderlines(m.PromptIdx, m.PromptIdx+1, m.PromptUnderlines)
		}
	}
}

// Types a character from the user input, and updates the underline string
func typeChar(m *Model, in string) {
//	lastWordIncr := 1
//	space := " "
//	if m.WordIdx == len(m.PromptSlice)-1 {
//		lastWordIncr = 2
//		space = ""
//	}

	//if m.InputLen < len(m.PromptSlice[m.WordIdx])+lastWordIncr {
	// Update underline pointer and add colors to characters for output
	if m.PromptIdx < len(m.PromptStr) {

		if in == string(m.PromptStr[m.PromptIdx]) {
			if _, ok := m.State.SeenIdxSet[m.PromptIdx]; !ok && in != " " {
				m.State.Hits++
				m.State.SeenIdxSet[m.PromptIdx] = 1
			}

			in = GREEN + in + RESET

		} else {
			if _, ok := m.State.SeenIdxSet[m.PromptIdx]; !ok && in != " " {
				m.State.Errors++
				m.State.SeenIdxSet[m.PromptIdx] = 1
			}
			in = RED + in + RESET
		}
		m.InputLen++

		m.PromptIdx++
		m.PromptUnderlines = updateUnderlines(m.PromptIdx, m.PromptIdx-1, m.PromptUnderlines)

		m.InputStr += in
	}

	// User typed word correctly
	inputStrPlain := removeColors(m.InputStr[m.PromptIdxLowerLimit:])
	if len(inputStrPlain) >= len(m.PromptSlice[m.WordIdx]) && strings.TrimSpace(removeColors(inputStrPlain)) == (m.PromptSlice[m.WordIdx]) {
    m.State.WordCompletions++
		m.WordIdx++
		m.PromptIdxLowerLimit = m.PromptIdx  
		m.InputStr = ""
		for range m.PromptIdxLowerLimit {
			m.InputStr += " "
		}
		m.InputLen = 0
	}
}

// Takes an ID and returns Prompt struct
func getPrompt(prompts []Prompt, promptID int) Prompt {
	return prompts[promptID-1]

}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.PromptStrsID == -2 {
		return m, tea.Quit
	}
	switch msg := msg.(type) {

	case tea.KeyMsg:
		in := msg.String()

		switch in {
		case "ctrl+c":
			// Quit on q or ctrl+c
			return m, tea.Quit
		case "enter", "ctrl+w", "ctrl+h", "ctrl+backspace", "tab", "ctrl+tab":

		case "backspace":
			backspace(&m)
		default:
			if isValidChar(in) {
				typeChar(&m, in)

				// Exit if finished
				if m.WordIdx > len(m.PromptSlice)-1 {
					m.State.PromptCompletions++

					// Select next prompt
					m.PromptStrsID = getNewPromptId(m.Cfg, m.PromptStrsID, m.State)
					// Exit the game
					if m.PromptStrsID == -2 {
						fmt.Println("Finished!")
						return m, nil
					}

					// Reinitialize variables
					m.PromptStr = getPrompt(m.Cfg.Prompts, m.PromptStrsID).Text
					m.PromptSlice = strings.Split(m.PromptStr, " ")
					m.PromptUnderlines = UNDERLINE_CHAR
					for range m.PromptStr {
						m.PromptUnderlines += " "
					}
					m.InputStr = ""
					m.PromptIdx = 0
					m.PromptIdxLowerLimit = 0
					m.WordIdx = 0
					m.InputLen = 0

					// Wipe the seen set
					m.State.SeenIdxSet = make(map[int]int)
				}
			}
		}
	case TickMsg:
		m.State.Time++
		return m, doTick()
	}

	return m, nil
}

func isValidChar(in string) bool {
	r := rune(in[0])
	return len(in) == 1 && r >= 32 && r <= 126
}

func getNewPromptId(cfg *Config, curID int, state *State) int {
	// Add the ID just completed to the seen set
	cfg.SeenIDs[curID] = 1

	// Exit game case
	if cfg.ActivePromptsLen == state.PromptCompletions {
		return -2
	}

	newID := rand.Intn(len(cfg.Prompts)) + 1 // 1 indexed

	// Keep looping if already used newID prompt
	_, ok := cfg.SeenIDs[newID]
	for ok || cfg.PromptType != "any" && getPrompt(cfg.Prompts, newID).Type != cfg.PromptType {
		newID = rand.Intn(len(cfg.Prompts)) + 1
		_, ok = cfg.SeenIDs[newID]
	}

	return newID
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

func getActivePromptsLen(pType string, prompts []Prompt) int {
	var res int

	if pType == "any" {
		return len(prompts)
	}

	for _, prompt := range prompts {
		if prompt.Type == pType {
			res++
		}
	}

	return res
}

func shiftCursor(m *Model) string {
	// Updates the '|' cursor line to the prompt string
	if m.PromptIdx+1 < len(m.PromptStr)+1 {
		return m.PromptStr[:m.PromptIdx] + "|" + string(m.PromptStr[m.PromptIdx]) + m.PromptStr[m.PromptIdx+1:]
	}

  return m.PromptStr
}

func updateWPM(s *State) {
  timeInMinutes := float32(s.Time)/float32(60)

  if timeInMinutes > 0 {
    s.WPM = (float32(s.WordCompletions)/timeInMinutes)
  } else {
    s.WPM = float32(s.WordCompletions)
  }
} 

func updateCPM(s *State) {
  timeInMinutes := float32(s.Time)/float32(60)

  if timeInMinutes > 0 {
    s.CPM = (float32(s.Hits)/timeInMinutes)
  } else {
    s.CPM = float32(s.Hits)
  }
} 

// View renders the UI
func (m Model) View() string {
  updateAccuracy(m.State)
  updateWPM(m.State)
  updateCPM(m.State)

	PromptStr := shiftCursor(&m)


	pType := m.Cfg.PromptTypeColor + "--------- " + m.Cfg.PromptType + " ---------" + RESET

  var display string

  display = fmt.Sprintf(
    "%s\n%s\n%s\n%s\nPrompt completions: %d\nWord completions: %d\nTime elapsed (s): %vs\nAccuracy: %.0f%%\nWPM: %.02f\nCPM: %0.02f\n\n", pType, PromptStr, m.PromptUnderlines, m.InputStr, m.State.PromptCompletions, m.State.WordCompletions, m.State.Time, m.State.Accuracy, m.State.WPM, m.State.CPM,)
// -2 means game should quit
	if m.PromptStrsID != -2 {
    return display
	}

	return display + GREEN+"\nFinished!\n"+RESET
}

func updateAccuracy(s *State) {
	if s.Hits > 0 {
		s.Accuracy = (1.0 - (float32(s.Errors) / float32(s.Hits))) * 100
	}
}

func Run() {
	var gameMode string
	err := huh.NewSelect[string]().
		Title("Select a mode").
		Options(
			huh.NewOption("standard (standard prompts for typing speed)", "standard"),
			huh.NewOption("coding (common coding motifs from different languages)", "coding"),
			huh.NewOption("any", "any"),
		).Value(&gameMode).Run()

	if err != nil {
		fmt.Println("Error: failed to run selected game mode")
		panic(err)
	}

	// Prompt type (standard, c++, Python etc...)
	var pType string
	var pTypeColor string

	switch gameMode {
	case "standard":
		pType = "standard"

	case "coding":
		err := huh.NewSelect[string]().
			Title("Select a coding language").
			Options(
				huh.NewOption("c++", "c++"),
				huh.NewOption("golang", "golang"),
				huh.NewOption("python", "python"),
				huh.NewOption("java", "java"),
			).Value(&pType).Run()
		if err != nil {
			fmt.Println("Error: failed to run selected game mode")
			panic(err)
		}

	case "any":
		pType = "any"

	default:
		panic("The game mode " + gameMode + " is not currently supported")
	}

	switch pType {
	case "c++":
		pTypeColor = BLUE
	case "python":
		pTypeColor = ORANGE
	case "golang":
		pTypeColor = TEAL
	case "rust":
		pTypeColor = BROWN
	case "java":
		pTypeColor = YELLOW
	case "any", "standard":
		pTypeColor = YELLOW
	}

	// Parse 'library.yaml' for a list of prompts
	cfg, err := parseYAML()
	if err != nil {
		log.Fatal(err)
	}

	cfg.PromptType = pType
	cfg.PromptTypeColor = pTypeColor
	cfg.ActivePromptsLen = getActivePromptsLen(pType, cfg.Prompts)

	state := State{
		SeenIdxSet: make(map[int]int),
	}

	cfg.SeenIDs = make(map[int]int)
	pStrsID := getNewPromptId(cfg, -1, &state)
	pStr := getPrompt(cfg.Prompts, pStrsID).Text
	pSlice := strings.Split(pStr, " ")
	pUnderlines := UNDERLINE_CHAR
	for range pStr {
		pUnderlines += " "
	}

	p := tea.NewProgram(Model{
		Cfg:         cfg,
		PromptStrsID:     pStrsID,
		PromptStr:        pStr,
		PromptSlice:      pSlice,
		PromptUnderlines: pUnderlines,
		State:       &state,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
