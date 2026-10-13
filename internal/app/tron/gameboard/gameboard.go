package gameboard

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const Height = 40
const Width = 80

type Color int

const (
	None Color = iota
	Red
	Blue
	Green
	Orange
	Pink
	Purple
	Beige
)

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9f6f2"))

var Colors = map[Color]lipgloss.Style{
	None:   defaultStyle.Background(lipgloss.NoColor{}),
	Blue:   defaultStyle.Background(lipgloss.Color("#063970")),
	Green:  defaultStyle.Background(lipgloss.Color("#4CA74F")),
	Orange: defaultStyle.Background(lipgloss.Color("#CF6209")),
	Pink:   defaultStyle.Background(lipgloss.Color("#D85B85")),
	Purple: defaultStyle.Background(lipgloss.Color("#9047A3")),
	Beige:  defaultStyle.Background(lipgloss.Color("#FFFDD0")),
	Red:    defaultStyle.Background(lipgloss.Color("#F20032")),
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) IsOpposite(other Direction) bool {
	switch d {
	case Up:
		return other == Down
	case Down:
		return other == Up
	case Left:
		return other == Right
	case Right:
		return other == Left
	}

	return false
}

type Gameboard struct {
	Grid            [Height][Width]Color
	playerPositions map[Color]position
}

type position struct {
	width  int
	height int
}

func (gb *Gameboard) Clear() {
	gb.Grid = [Height][Width]Color{}
}

func NewGameboard(playerColors ...Color) Gameboard {
	gb := Gameboard{
		playerPositions: make(map[Color]position),
	}
	initialPositions := make([]position, 0)

	switch len(playerColors) {
	case 2:
		initialHeight := Height / 2
		initialPositions = append(initialPositions, position{Width/8 - 1, initialHeight})
		initialPositions = append(initialPositions, position{Width - Width/8, initialHeight})
	case 3:
		initialPositions = append(initialPositions, position{1, 1})
		initialPositions = append(initialPositions, position{Width - 2, 1})
		initialPositions = append(initialPositions, position{Width / 2, Height - 1})
	default:
		panic("Only 2 or 3 number of players are supported.")
	}

	for i, position := range initialPositions {
		gb.markPositionForPlayer(position, playerColors[i])
		gb.playerPositions[playerColors[i]] = position
	}

	return gb
}

func (gb *Gameboard) MovePlayer(playerColor Color, direction Direction, speed int) bool {
	currentPosition, ok := gb.playerPositions[playerColor]

	if !ok {
		panic("Unknown player Color: " + strconv.Itoa(int(playerColor)))
	}

	for i := 0; i < speed; i++ {
		var newPosition position
		switch direction {
		case Up:
			newPosition = position{currentPosition.width, currentPosition.height - 1}
		case Down:
			newPosition = position{currentPosition.width, currentPosition.height + 1}
		case Left:
			newPosition = position{currentPosition.width - 1, currentPosition.height}
		case Right:
			newPosition = position{currentPosition.width + 1, currentPosition.height}
		}
		if !gb.isPositionValid(newPosition) {
			return false
		}
		currentPosition = newPosition
		gb.markPositionForPlayer(currentPosition, playerColor)

	}
	gb.playerPositions[playerColor] = currentPosition

	return true
}

func (gb *Gameboard) isPositionValid(position position) bool {
	if position.height < 0 || position.width < 0 {
		return false
	}

	if position.height >= Height || position.width >= Width {
		return false
	}

	return gb.Grid[position.height][position.width] == None
}

func (gb *Gameboard) markPositionForPlayer(position position, playerColor Color) {
	gb.Grid[position.height][position.width] = playerColor
}
