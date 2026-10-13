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
	LightRed
	Blue
	LightBlue
	Green
	Orange
	Pink
	Purple
	LightPurple
	Beige
)

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9f6f2"))

var Colors = map[Color]lipgloss.Style{
	None:        defaultStyle.Background(lipgloss.NoColor{}),
	Blue:        defaultStyle.Background(lipgloss.Color(Blue.Hex())),
	LightBlue:   defaultStyle.Background(lipgloss.Color(LightBlue.Hex())),
	Green:       defaultStyle.Background(lipgloss.Color(Green.Hex())),
	Orange:      defaultStyle.Background(lipgloss.Color(Orange.Hex())),
	Pink:        defaultStyle.Background(lipgloss.Color(Pink.Hex())),
	Purple:      defaultStyle.Background(lipgloss.Color(Purple.Hex())),
	LightPurple: defaultStyle.Background(lipgloss.Color(LightPurple.Hex())),
	Beige:       defaultStyle.Background(lipgloss.Color(Beige.Hex())),
	Red:         defaultStyle.Background(lipgloss.Color(Red.Hex())),
	LightRed:    defaultStyle.Background(lipgloss.Color(LightRed.Hex())),
}

func (c Color) animate() Color {
	switch c {
	case None:
		return None
	case Red:
		return LightRed
	case LightRed:
		return Red
	case Blue:
		return LightBlue
	case LightBlue:
		return Blue
	case Green:
		return Green
	case Orange:
		return Orange
	case Pink:
		return Pink
	case Purple:
		return LightPurple
	case LightPurple:
		return Purple
	case Beige:
		return Beige
	}

	return None
}

func (c Color) String() string {
	switch c {
	case None:
		return "None"
	case Red:
		return "Red"
	case LightRed:
		return "LightRed"
	case Blue:
		return "Blue"
	case LightBlue:
		return "LightBlue"
	case Green:
		return "Green"
	case Orange:
		return "Orange"
	case Pink:
		return "Pink"
	case Purple:
		return "Purple"
	case LightPurple:
		return "LightPurple"
	case Beige:
		return "Beige"
	}

	return "None"
}

func (c Color) Hex() string {
	switch c {
	case None:
		return ""
	case Red:
		return "#F20032"
	case LightRed:
		return "#FFB0B2"
	case Blue:
		return "#2F49FE"
	case LightBlue:
		return "#63D4F2"
	case Green:
		return "#4CA74F"
	case Orange:
		return "#CF6209"
	case Pink:
		return "#D85B85"
	case Purple:
		return "#9047A3"
	case LightPurple:
		return "#DDB2E1"
	case Beige:
		return "#FFFDD0"
	}

	return ""
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
	Grid         [Height][Width]Color
	playerStates map[Color]*playerState
}

type playerState struct {
	position position
	status   status
}

type position struct {
	width  int
	height int
}

type status struct {
	crashed            bool
	animationCountdown int
	removed            bool
}

func (gb *Gameboard) Clear() {
	gb.Grid = [Height][Width]Color{}
}

func NewGameboard(playerColors ...Color) Gameboard {
	gb := Gameboard{
		playerStates: make(map[Color]*playerState),
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
		gb.playerStates[playerColors[i]] = &playerState{
			position: position,
		}
	}

	return gb
}

func (gb *Gameboard) UpdatePlayer(playerColor Color, direction Direction, speed int, gameFinished bool) bool {
	playerState, ok := gb.playerStates[playerColor]

	if !ok {
		panic("Unknown player Color: " + strconv.Itoa(int(playerColor)))
	}

	if playerState.status.removed {
		return true
	}

	if playerState.status.crashed {
		gb.animateCrashedPlayer(playerColor, playerState)
		return true
	} else if !gameFinished {
		return gb.movePlayer(playerState, playerColor, direction, speed)
	}

	return true
}

func (gb *Gameboard) movePlayer(playerState *playerState, playerColor Color, direction Direction, speed int) bool {
	for i := 0; i < speed; i++ {
		var newPosition position
		switch direction {
		case Up:
			newPosition = position{playerState.position.width, playerState.position.height - 1}
		case Down:
			newPosition = position{playerState.position.width, playerState.position.height + 1}
		case Left:
			newPosition = position{playerState.position.width - 1, playerState.position.height}
		case Right:
			newPosition = position{playerState.position.width + 1, playerState.position.height}
		}
		if !gb.isPositionValid(newPosition) {
			markPlayerAsCrashed(playerState)
			return false
		}
		playerState.position = newPosition
		gb.markPositionForPlayer(playerState.position, playerColor)

	}

	return true
}

func (gb *Gameboard) animateCrashedPlayer(playerColor Color, playerState *playerState) {
	var newColor, oldColor Color
	if playerState.status.animationCountdown%2 == 0 {
		newColor = playerColor.animate()
		oldColor = playerColor
	} else {
		oldColor = playerColor.animate()
		newColor = playerColor
	}

	if playerState.status.animationCountdown == 0 {
		newColor = None
		playerState.status.removed = true
	}

	for y, row := range gb.Grid {
		for x := range row {
			if gb.Grid[y][x] == oldColor {
				gb.Grid[y][x] = newColor
			}
		}
	}

	playerState.status.animationCountdown = playerState.status.animationCountdown - 1
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

func markPlayerAsCrashed(playerState *playerState) {
	playerState.status.crashed = true
	playerState.status.animationCountdown = 4
}
