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
	"github.com/charmbracelet/lipgloss"
)


const (
	RESET       = "\033[0m"
	RED         = "\033[31m"
	RED_HEX = "#FF2900"
	GREEN_HEX = "#00FF00"
	DARK_BLUE_HEX = "#0011FF"
	YELLOW_HEX = "#FFD900"
	ORANGE_HEX = "#FF8C00"
	TEAL_HEX = "#00D5FF"
	BROWN_HEX = "#A34900"
	PURPLE_HEX = "#9500FF"
	WHITE_HEX = "#FFFFFF"

	CURSOR_CHAR = "█"

	RESET_LEN      = len(RESET)
	COLOR_LEN      = len(RED)
	UNDERLINE_CHAR = " "
)

var (
	RED_HEX_LEN = len(lipgloss.NewStyle().Foreground(lipgloss.Color(RED_HEX)).Render("a"))
	GREEN_HEX_LEN = len(lipgloss.NewStyle().Foreground(lipgloss.Color(GREEN_HEX)).Render("a"))
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

	// Correct chars typed/minute
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

	InputStrPlain string

	// Length of user input string
	InputLen int

	CharLenSlice []int

	State *State
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func backspace(m *Model) {
	if m.InputLen > 0 {
		m.InputStr = m.InputStr[:len(m.InputStr)-m.CharLenSlice[len(m.CharLenSlice)-1]] //CHAR_LEN]
		m.InputLen--
		m.CharLenSlice = m.CharLenSlice[:len(m.CharLenSlice)-1] // pop
		m.InputStrPlain = m.InputStrPlain[:len(m.InputStrPlain)-1]

		// Decrement underline pointer if greater than lower limit
		if m.PromptIdx > m.PromptIdxLowerLimit {
			m.PromptIdx--
			m.PromptUnderlines = updateUnderlines(m.PromptIdx, m.PromptIdx+1, m.PromptUnderlines)
		}
	}
}

func typeChar(m *Model, in string) {
	if m.PromptIdx < len(m.PromptStr) {
		m.InputStrPlain += in

		if in == string(m.PromptStr[m.PromptIdx]) {

			if _, ok := m.State.SeenIdxSet[m.PromptIdx]; !ok && in != " " {
				m.State.Hits++
				m.State.SeenIdxSet[m.PromptIdx] = 1
			}

			greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(GREEN_HEX))
			in = greenStyle.Render(in)
			m.CharLenSlice = append(m.CharLenSlice, GREEN_HEX_LEN)

		} else {
			if _, ok := m.State.SeenIdxSet[m.PromptIdx]; !ok && in != " " {
				m.State.Errors++
				m.State.SeenIdxSet[m.PromptIdx] = 1
			}
			redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(RED_HEX))
			in = redStyle.Render(in)
			m.CharLenSlice = append(m.CharLenSlice, RED_HEX_LEN)
		} 

		m.InputLen++

		m.PromptIdx++
		m.PromptUnderlines = updateUnderlines(m.PromptIdx, m.PromptIdx-1, m.PromptUnderlines)

		m.InputStr += in
	}

	// User typed word correctly
	if len(m.InputStrPlain) >= len(m.PromptSlice[m.WordIdx]) && strings.TrimSpace(m.InputStrPlain) == (m.PromptSlice[m.WordIdx]) {
		m.State.WordCompletions++
		m.WordIdx++
		m.PromptIdxLowerLimit = m.PromptIdx
		m.InputStr = ""
		m.InputStrPlain = ""
		for range m.PromptIdxLowerLimit {
			m.InputStr += " "
		}
		m.InputLen = 0
	}
}

func getPrompt(prompts []Prompt, promptID int) Prompt {
	return prompts[promptID-1]
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.PromptStrsID == -2 {
		return m, tea.Quit
	}
	switch msg := msg.(type) {

	case tea.KeyMsg:
		in := msg.String()

		switch in {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", "ctrl+w", "ctrl+h", "ctrl+backspace", "tab", "ctrl+tab":

		case "backspace":
			backspace(&m)
		default:
			if isValidChar(in) {
				typeChar(&m, in)

				if m.WordIdx > len(m.PromptSlice)-1 {
					m.State.PromptCompletions++

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
	cfg.SeenIDs[curID] = 1

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
	if m.PromptIdx+1 < len(m.PromptStr)+1 {
		return m.PromptStr[:m.PromptIdx] + CURSOR_CHAR + string(m.PromptStr[m.PromptIdx]) + m.PromptStr[m.PromptIdx+1:]
	}

	return m.PromptStr
}

func updateWPM(s *State) {
	timeInMinutes := float32(s.Time) / float32(60)

	if timeInMinutes > 0 {
		s.WPM = (float32(s.WordCompletions) / timeInMinutes)
	} else {
		s.WPM = float32(s.WordCompletions)
	}
}

func updateCPM(s *State) {
	timeInMinutes := float32(s.Time) / float32(60)

	if timeInMinutes > 0 {
		s.CPM = (float32(s.Hits) / timeInMinutes)
	} else {
		s.CPM = float32(s.Hits)
	}
}

func (m Model) View() string {
	updateAccuracy(m.State)
	updateWPM(m.State)
	updateCPM(m.State)

	PromptStr := shiftCursor(&m)

	var display string

	display = fmt.Sprintf(
		"%s\n%s\n%s\n%s\nPrompt completions: %d\nWord completions: %d\nTime elapsed (s): %vs\nAccuracy: %.0f%%\nWPM: %.02f\nCPM: %0.02f\n\n", m.Cfg.PromptFormattedPrintString, PromptStr, m.PromptUnderlines, m.InputStr, m.State.PromptCompletions, m.State.WordCompletions, m.State.Time, m.State.Accuracy, m.State.WPM, m.State.CPM)

	// -2 means game should quit
	if m.PromptStrsID != -2 {
		return display
	}

	return display + lipgloss.NewStyle().Foreground(lipgloss.Color(GREEN_HEX)).Render("\nFinished!\n")  
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

	var pType string
	var pTypeColor string
	var hexColor string

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
				huh.NewOption("rust", "rust"),
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
		hexColor = DARK_BLUE_HEX
	case "python":
		hexColor = ORANGE_HEX
	case "golang":
		hexColor = TEAL_HEX
	case "rust":
		hexColor = BROWN_HEX
	case "java":
		hexColor = YELLOW_HEX
	case "any", "standard":
		hexColor = WHITE_HEX
	}

	pTypeColor = lipgloss.NewStyle().Foreground(lipgloss.Color(hexColor)).Render("---------" + pType + "---------")

	cfg, err := parseYAML()
	if err != nil {
		log.Fatal(err)
	}

	cfg.PromptType = pType
	cfg.PromptFormattedPrintString = pTypeColor
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
	charLenSlice := []int{}

	p := tea.NewProgram(Model{
		Cfg:              cfg,
		PromptStrsID:     pStrsID,
		PromptStr:        pStr,
		PromptSlice:      pSlice,
		PromptUnderlines: pUnderlines,
		State:            &state,
		CharLenSlice: charLenSlice,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
