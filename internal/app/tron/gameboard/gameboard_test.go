package gameboard

import (
	"testing"
)

func TestWrongNumberOfPlayersPanics(t *testing.T) {
	testCases := []struct {
		name string
		in   []Color
	}{
		{"no players", []Color{}},
		{"more than three playrs", []Color{
			Red,
			Green,
			Blue,
			None,
		},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("should panic when a wrong number of players is used")
				}
			}()

			NewGameboard(tt.in...)
		})
	}

}

func TestTwoPlayerInitilization(t *testing.T) {
	gameboard := NewGameboard(Red, Blue)

	if gameboard.grid[20][9] != Red {
		t.Error("expected that first player position was 9,20")
	}

	if gameboard.grid[20][70] != Blue {
		t.Error("expected that second player position was 71,20")
	}
}

func TestPlayersCanCrashOneOnAnother(t *testing.T) {
	gameboard := NewGameboard(Red, Blue)

	if !gameboard.MovePlayer(Red, Right, 10) {
		t.Error("first player should be able to move 10 positions to the right")
	}

	if !gameboard.MovePlayer(Red, Right, 50) {
		t.Error("first player should be able to move just left of the other player")
	}

	if gameboard.MovePlayer(Red, Right, 1) {
		t.Error("first player should not be able to move on the second player")
	}
}

func TestPlayerShouldCrashWhenGoingOutOfBounds(t *testing.T) {
	testCases := []struct {
		name      string
		color     Color
		direction Direction
		speed     int
		want      bool
	}{
		{
			"should be able to go before the top bound",
			Red,
			Up,
			20,
			true,
		},
		{
			"should crash when going over the top bound",
			Red,
			Up,
			21,
			false,
		},
		{
			"should be able to go before the bottom bound",
			Red,
			Down,
			19,
			true,
		},
		{
			"should crash when going over the bottom bound",
			Red,
			Down,
			20,
			false,
		},
		{
			"should be able to go before the left bound",
			Red,
			Left,
			9,
			true,
		},
		{
			"should crash when going over the left bound",
			Red,
			Left,
			10,
			false,
		},
		{
			"should be able to go before the right bound",
			Blue,
			Right,
			9,
			true,
		},
		{
			"should crash when going over the right bound",
			Blue,
			Right,
			10,
			false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gameboard := NewGameboard(Red, Blue)

			if gameboard.MovePlayer(tt.color, tt.direction, tt.speed) != tt.want {
				t.Errorf("Expected MovePlayer to return %t", tt.want)
			}
		})
	}
}
