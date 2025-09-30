package blackjack

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Card struct {
	Suit string
	Rank string
}

type Deck []Card

type model struct {
	deck         Deck
	playerHand   []Card
	dealerHand   []Card
	playerTurn   bool
	gameOver     bool
	message      string
	playerStyle  lipgloss.Style
	dealerStyle  lipgloss.Style
	defaultStyle lipgloss.Style
}

func NewDeck() Deck {
	suits := []string{"♥", "♦", "♣", "♠"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	deck := make(Deck, 0, 52)

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}
	return deck
}

func (d Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func (d *Deck) Draw() Card {
	card := (*d)[0]
	*d = (*d)[1:]
	return card
}

func HandValue(hand []Card) int {
	value := 0
	aces := 0
	for _, card := range hand {
		switch card.Rank {
		case "A":
			aces++
			value += 11
		case "K", "Q", "J":
			value += 10
		default:
			rankValue := 0
			fmt.Sscanf(card.Rank, "%d", &rankValue)
			value += rankValue
		}
	}

	for value > 21 && aces > 0 {
		value -= 10
		aces--
	}
	return value
}

func initialModel() tea.Model {
	deck := NewDeck()
	deck.Shuffle()

	playerHand := []Card{deck.Draw(), deck.Draw()}
	dealerHand := []Card{deck.Draw(), deck.Draw()}

	return model{
		deck:         deck,
		playerHand:   playerHand,
		dealerHand:   dealerHand,
		playerTurn:   true,
		gameOver:     false,
		message:      "Hit (h) or Stand (s)?",
		playerStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("99")),
		dealerStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("208")),
		defaultStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			if m.playerTurn && !m.gameOver {
				m.playerHand = append(m.playerHand, m.deck.Draw())
				if HandValue(m.playerHand) > 21 {
					m.message = "Player busts! Dealer wins."
					m.gameOver = true
				}
			}
		case "s":
			if m.playerTurn && !m.gameOver {
				m.playerTurn = false
				// Dealer's turn
				for HandValue(m.dealerHand) < 17 {
					m.dealerHand = append(m.dealerHand, m.deck.Draw())
				}

				playerValue := HandValue(m.playerHand)
				dealerValue := HandValue(m.dealerHand)

				if dealerValue > 21 || playerValue > dealerValue {
					m.message = "Player wins!"
				} else if dealerValue > playerValue {
					m.message = "Dealer wins!"
				} else {
					m.message = "Push (Tie)!"
				}
				m.gameOver = true
			}
		case "n":
			if m.gameOver {
				return initialModel(), nil
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Blackjack\n\n"

	s += m.dealerStyle.Render("Dealer's Hand:\n")
	if m.playerTurn && !m.gameOver {
		s += fmt.Sprintf("[ %s %s ] [ ? ]\n", m.dealerHand[0].Rank, m.dealerHand[0].Suit)
	} else {
		for _, card := range m.dealerHand {
			s += fmt.Sprintf("[ %s %s ] ", card.Rank, card.Suit)
		}
		s += fmt.Sprintf(" (Value: %d)\n", HandValue(m.dealerHand))
	}

	s += "\n" + m.playerStyle.Render("Player's Hand:\n")
	for _, card := range m.playerHand {
		s += fmt.Sprintf("[ %s %s ] ", card.Rank, card.Suit)
	}
	s += fmt.Sprintf(" (Value: %d)\n", HandValue(m.playerHand))

	s += "\n" + m.defaultStyle.Render(m.message) + "\n"
	if m.gameOver {
		s += "\nPress 'q' to quit or 'n' to start a new game.\n"
	}

	return s
}

func Run() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
