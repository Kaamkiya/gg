// Package gameboard maintains the game's play area.
//
// It is responsible for drawing the player's paths and detecting
// collisions. The gameboard is not dependent on players and any related
// player state is kept internally in the package. A player is represented
// only by a [Color] and each position in the grid can either have a [Color]
// when it is full or [lipgloss.NoColor] when empty.
package gameboard

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const Height = 40
const Width = 80
const crashAnimationCountdown = 4

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

// animate provides an alternate color to be used in flashy animations.
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

type Board struct {
	Grid         [Height][Width]Color
	playerStates map[Color]*playerState
}

type playerState struct {
	position           position
	animationOngoing   bool
	animationCountdown int
}

type position struct {
	width  int
	height int
}

func (gb *Board) Clear() {
	gb.Grid = [Height][Width]Color{}
}

// NewBoard initializes the play area for two or three players.
func NewBoard(playerColors ...Color) *Board {
	gb := &Board{
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
			position:           position,
			animationCountdown: -1,
		}
	}

	return gb
}

// MovePlayer updates the player position. It should be called at each loop tick for all players while the game is not finished.
//
// Returns whether the player has crashed.
func (gb *Board) MovePlayer(playerColor Color, direction Direction, speed int) (crashed bool) {
	playerState, ok := gb.playerStates[playerColor]

	if !ok {
		panic("Unknown player Color: " + strconv.Itoa(int(playerColor)))
	}

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
			return true
		}
		playerState.position = newPosition
		gb.markPositionForPlayer(playerState.position, playerColor)

	}

	return false
}

// AnimateCrashedPlayer performs a flashy animation when the player crashes. Once the animation
// is completed the player is removed from the player area and his spaces become empty.
func (gb *Board) AnimatePlayer(playerColor Color) {
	playerState, ok := gb.playerStates[playerColor]

	if !ok {
		panic("Unknown player Color: " + strconv.Itoa(int(playerColor)))
	}

	if !playerState.animationOngoing {
		return
	}

	var newColor, oldColor Color
	if playerState.animationCountdown%2 == 0 {
		newColor = playerColor.animate()
		oldColor = playerColor
	} else {
		oldColor = playerColor.animate()
		newColor = playerColor
	}

	if playerState.animationCountdown == 0 {
		newColor = None
		playerState.animationOngoing = false
	}

	for y, row := range gb.Grid {
		for x := range row {
			if gb.Grid[y][x] == oldColor {
				gb.Grid[y][x] = newColor
			}
		}
	}

	playerState.animationCountdown -= 1
}

func (gb *Board) isPositionValid(position position) bool {
	if position.height < 0 || position.width < 0 {
		return false
	}

	if position.height >= Height || position.width >= Width {
		return false
	}

	return gb.Grid[position.height][position.width] == None
}

func (gb *Board) markPositionForPlayer(position position, playerColor Color) {
	gb.Grid[position.height][position.width] = playerColor
}

func markPlayerAsCrashed(playerState *playerState) {
	playerState.animationOngoing = true
	playerState.animationCountdown = crashAnimationCountdown
}
