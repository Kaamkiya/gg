package solitaire

import (
	"fmt"
	"strings"

	"math/rand/v2"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type suit int

const (
	spades suit = iota
	hearts
	diamonds
	clubs
)

type card struct {
	rank   int  // 1..13 => A..K
	suit   suit
	faceUp bool
}

func (c card) colorStyle(red, black lipgloss.Style) lipgloss.Style {
	if c.suit == hearts || c.suit == diamonds {
		return red
	}
	return black
}

func (c card) suitRune() string {
	switch c.suit {
	case spades:
		return "♠"
	case hearts:
		return "♥"
	case diamonds:
		return "♦"
	default:
		return "♣"
	}
}

func (c card) rankString() string {
	switch c.rank {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		return fmt.Sprintf("%d", c.rank)
	}
}

func (c card) String(red, black lipgloss.Style) string {
	if !c.faceUp {
		return "XX"
	}
	val := c.rankString() + c.suitRune()
	return c.colorStyle(red, black).Render(val)
}

func isRed(s suit) bool { return s == hearts || s == diamonds }

// Game model

type sourceKind int

const (
	srcNone sourceKind = iota
	srcWaste
	srcTableau
)

type source struct {
	kind sourceKind
	idx  int // for tableau index when kind == srcTableau
}

type model struct {
	stock       []card
	waste       []card
	foundations [4][]card
	tableau     [7][]card

	selected source
	message  string

	redStyle   lipgloss.Style
	blackStyle lipgloss.Style
}

func initialModel() tea.Model {
	// Build and shuffle deck
	deck := make([]card, 0, 52)
	for s := 0; s < 4; s++ {
		for r := 1; r <= 13; r++ {
			deck = append(deck, card{rank: r, suit: suit(s), faceUp: false})
		}
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	m := model{
		stock:      []card{},
		waste:      []card{},
		selected:   source{kind: srcNone},
		redStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")),
		blackStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#dddddd")),
	}

	// Deal tableau: columns 0..6 get 1..7 cards, last is face up
	pos := 0
	for col := 0; col < 7; col++ {
		for i := 0; i <= col; i++ {
			c := deck[pos]
			pos++
			if i == col {
				c.faceUp = true
			}
			m.tableau[col] = append(m.tableau[col], c)
		}
	}

	// Remaining go to stock (face down)
	m.stock = append(m.stock, deck[pos:]...)
	return m
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) draw() {
	if len(m.stock) == 0 {
		// Recycle waste back into stock (face down, reverse order)
		if len(m.waste) == 0 {
			m.message = "nothing to draw"
			return
		}
		for i := len(m.waste) - 1; i >= 0; i-- {
			c := m.waste[i]
			c.faceUp = false
			m.stock = append(m.stock, c)
		}
		m.waste = m.waste[:0]
		m.message = "recycled waste back to stock"
		return
	}
	// Draw one card to waste
	c := m.stock[len(m.stock)-1]
	m.stock = m.stock[:len(m.stock)-1]
	c.faceUp = true
	m.waste = append(m.waste, c)
	m.message = "drew a card"
}

func (m *model) topWaste() (card, bool) {
	if len(m.waste) == 0 {
		return card{}, false
	}
	return m.waste[len(m.waste)-1], true
}

func (m *model) topTableau(i int) (card, bool) {
	col := m.tableau[i]
	if len(col) == 0 {
		return card{}, false
	}
	return col[len(col)-1], true
}

func (m *model) flipIfNeeded(i int) {
	col := m.tableau[i]
	if len(col) == 0 {
		return
	}
	last := &col[len(col)-1]
	if !last.faceUp {
		last.faceUp = true
	}
	m.tableau[i] = col
}

func canPlaceOnTableau(c card, dest []card) bool {
	if len(dest) == 0 {
		return c.rank == 13 // kings on empty
	}
	top := dest[len(dest)-1]
	if !top.faceUp {
		return false
	}
	// alternating colors and descending rank
	if isRed(c.suit) == isRed(top.suit) {
		return false
	}
	return c.rank == top.rank-1
}

func canMoveToFoundation(c card, dest []card) bool {
	if !c.faceUp {
		return false
	}
	if len(dest) == 0 {
		return c.rank == 1
	}
	top := dest[len(dest)-1]
	return c.suit == top.suit && c.rank == top.rank+1
}

func (m *model) tryAutoFoundation() bool {
	// Try waste first
	if wc, ok := m.topWaste(); ok {
		fi := findFoundationIndexFor(m, wc)
		if fi >= 0 && canMoveToFoundation(wc, m.foundations[fi]) {
			m.foundations[fi] = append(m.foundations[fi], wc)
			m.waste = m.waste[:len(m.waste)-1]
			m.message = "moved waste -> foundation"
			return true
		}
	}
	// Then each tableau top
	for i := 0; i < 7; i++ {
		if tc, ok := m.topTableau(i); ok && tc.faceUp {
			fi := findFoundationIndexFor(m, tc)
			if fi >= 0 && canMoveToFoundation(tc, m.foundations[fi]) {
				m.foundations[fi] = append(m.foundations[fi], tc)
				m.tableau[i] = m.tableau[i][:len(m.tableau[i])-1]
				m.flipIfNeeded(i)
				m.message = fmt.Sprintf("moved T%d -> foundation", i+1)
				return true
			}
		}
	}
	return false
}

func findFoundationIndexFor(m *model, c card) int {
	// Prefer exact suit pile if it started, else any empty
	for i := 0; i < 4; i++ {
		pile := m.foundations[i]
		if len(pile) == 0 {
			continue
		}
		if pile[len(pile)-1].suit == c.suit {
			return i
		}
	}
	// Not found: choose empty foundation
	for i := 0; i < 4; i++ {
		if len(m.foundations[i]) == 0 {
			return i
		}
	}
	return -1
}

func (m *model) moveWasteToTableau(dst int) bool {
	wc, ok := m.topWaste()
	if !ok {
		m.message = "waste empty"
		return false
	}
	if !canPlaceOnTableau(wc, m.tableau[dst]) {
		m.message = "cannot place waste there"
		return false
	}
	m.tableau[dst] = append(m.tableau[dst], wc)
	m.waste = m.waste[:len(m.waste)-1]
	m.message = fmt.Sprintf("waste -> T%d", dst+1)
	return true
}

func (m *model) moveTableauToTableau(src, dst int) bool {
	srcCol := m.tableau[src]
	if len(srcCol) == 0 {
		m.message = "source empty"
		return false
	}
	// Find first face-up index
	firstFace := -1
	for i := range srcCol {
		if srcCol[i].faceUp {
			firstFace = i
			break
		}
	}
	if firstFace == -1 {
		m.message = "no face-up run"
		return false
	}
	// Find the minimal index within the face-up run that can be placed on dst
	for i := firstFace; i < len(srcCol); i++ {
		lead := srcCol[i]
		if canPlaceOnTableau(lead, m.tableau[dst]) && isValidDescendingAlt(srcCol[i:]) {
			// Move the run
			m.tableau[dst] = append(m.tableau[dst], srcCol[i:]...)
			m.tableau[src] = srcCol[:i]
			m.flipIfNeeded(src)
			m.message = fmt.Sprintf("T%d -> T%d", src+1, dst+1)
			return true
		}
	}
	m.message = "no valid move"
	return false
}

func isValidDescendingAlt(run []card) bool {
	if len(run) == 0 {
		return false
	}
	for i := 0; i < len(run)-1; i++ {
		if !run[i].faceUp || !run[i+1].faceUp {
			return false
		}
		if isRed(run[i].suit) == isRed(run[i+1].suit) {
			return false
		}
		if run[i].rank != run[i+1].rank+1 {
			return false
		}
	}
	return true
}

func (m *model) hasWon() bool {
	total := 0
	for i := 0; i < 4; i++ {
		total += len(m.foundations[i])
	}
	return total == 52
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "d", " ":
			m.draw()
		case "0":
			m.selected = source{kind: srcWaste}
			m.message = "selected: waste"
		case "f":
			if !m.tryAutoFoundation() {
				m.message = "no foundation moves"
			}
		case "esc":
			m.selected = source{kind: srcNone}
			m.message = "cleared selection"
		case "1", "2", "3", "4", "5", "6", "7":
			idx := int(key[0]-'1')
			if m.selected.kind == srcNone {
				m.selected = source{kind: srcTableau, idx: idx}
				m.message = fmt.Sprintf("selected: T%d", idx+1)
			} else if m.selected.kind == srcWaste {
				m.moveWasteToTableau(idx)
				m.selected = source{kind: srcNone}
			} else if m.selected.kind == srcTableau {
				if idx == m.selected.idx {
					// same column pressed again -> try top to foundation
					if c, ok := m.topTableau(idx); ok {
						fi := findFoundationIndexFor(&m, c)
						if fi >= 0 && canMoveToFoundation(c, m.foundations[fi]) {
							m.foundations[fi] = append(m.foundations[fi], c)
							m.tableau[idx] = m.tableau[idx][:len(m.tableau[idx])-1]
							m.flipIfNeeded(idx)
							m.message = "moved tableau -> foundation"
						} else {
							m.message = "cannot move to foundation"
						}
					} else {
						m.message = "empty column"
					}
					m.selected = source{kind: srcNone}
				} else {
					m.moveTableauToTableau(m.selected.idx, idx)
					m.selected = source{kind: srcNone}
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString("Solitaire (Klondike, 1-card draw)\n")
	b.WriteString("q: quit  d/space: draw  0: select waste  1-7: select/move tableau  f: auto-foundation  esc: clear\n\n")

	// Header: stock, waste top, foundations top
	stockCount := len(m.stock)
	wasteStr := "--"
	if c, ok := m.topWaste(); ok {
		wasteStr = c.String(m.redStyle, m.blackStyle)
	}
	b.WriteString(fmt.Sprintf("Stock[%d]  Waste[%s]    Foundations: ", stockCount, wasteStr))
	for i := 0; i < 4; i++ {
		pile := m.foundations[i]
		if len(pile) == 0 {
			b.WriteString("[__] ")
			continue
		}
		b.WriteString("[" + pile[len(pile)-1].String(m.redStyle, m.blackStyle) + "] ")
	}
	b.WriteString("\n")
	if m.selected.kind == srcWaste {
		b.WriteString("Selected: waste\n")
	} else if m.selected.kind == srcTableau {
		b.WriteString(fmt.Sprintf("Selected: T%d\n", m.selected.idx+1))
	} else {
		b.WriteString("Selected: none\n")
	}
	if m.message != "" {
		b.WriteString(m.message + "\n")
	}
	b.WriteString("\n")

	// Compute max height across tableau (including face-down)
	maxH := 0
	for i := 0; i < 7; i++ {
		if len(m.tableau[i]) > maxH {
			maxH = len(m.tableau[i])
		}
	}
	// Column headers
	b.WriteString("   ")
	for i := 0; i < 7; i++ {
		label := fmt.Sprintf(" T%d ", i+1)
		b.WriteString(label)
		b.WriteString("   ")
	}
	b.WriteString("\n")
	// Rows
	for row := 0; row < maxH; row++ {
		b.WriteString("   ")
		for col := 0; col < 7; col++ {
			pile := m.tableau[col]
			if row < len(pile) {
				b.WriteString("[" + pile[row].String(m.redStyle, m.blackStyle) + "]")
			} else {
				b.WriteString("[  ]")
			}
			b.WriteString("   ")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func Run() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
