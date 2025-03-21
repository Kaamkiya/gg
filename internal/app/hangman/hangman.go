package hangman

import (
	"math/rand/v2"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	word     string
	showWord []rune
	guesses  int
	guessed  []string
	art      []string
}

func initialModel() tea.Model {
	word := wordlist[rand.IntN(len(wordlist))]

	showWord := make([]rune, len(word))
	for i := range word {
		showWord[i] = '_'
	}

	art := []string{
		` +--+
 |  |
    |
    |
    |
    |
=====`,
		` +--+
 |  |
 O  |
    |
    |
    |
=====`,
		` +--+
 |  |
 O  |
 |  |
    |
    |
=====`,
		` +--+
 |  |
 O  |
/|  |
    |
    |
=====`,
		` +--+
 |  |
 O  |
/|\ |
    |
    |
=====`,
		` +--+
 |  |
 O  |
/|\ |
/   |
    |
=====`,
		` +--+
 |  |
 O  |
/|\ |
/ \ |
    |
=====`,
	}

	return model{
		word:     word,
		showWord: showWord,
		guesses:  6,
		guessed:  []string{},
		art:      art,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
			letter := msg.String()
			if slices.Contains(m.guessed, letter) {
				return m, nil
			}

			inWord := false
			for i, char := range m.word {
				if string(char) == letter {
					m.showWord[i] = char
					inWord = true
				}
			}

			if !inWord {
				m.guessed = append(m.guessed, letter)
				m.guesses--
			}
		}
	}

	if m.guesses <= -1 || m.word == string(m.showWord) {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	if m.guesses < 0 {
		s += m.art[len(m.art)-1]
	} else {
		s += m.art[6-m.guesses]
	}

	s += "\n\nGuessed: "
	for _, guessed := range m.guessed {
		s += guessed
	}

	s += "\n\nWord: "
	for _, char := range m.showWord {
		s += string(char)
	}

	s += "\n\n"

	if m.guesses < 0 {
		s += `The word was "` + m.word + "\".\n\n"
	}

	return s
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
