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

			NewBoard(tt.in...)
		})
	}

}

func TestTwoPlayerInitilization(t *testing.T) {
	gameboard := NewBoard(Red, Blue)

	if gameboard.Grid[20][9] != Red {
		t.Error("expected that first player position was 9,20")
	}

	if gameboard.Grid[20][70] != Blue {
		t.Error("expected that second player position was 71,20")
	}
}

func TestPlayersCanCrashOneOnAnother(t *testing.T) {
	gameboard := NewBoard(Red, Blue)

	if gameboard.UpdatePlayer(Red, Right, 10, false) {
		t.Error("first player should be able to move 10 positions to the right")
	}

	if gameboard.UpdatePlayer(Red, Right, 50, false) {
		t.Error("first player should be able to move just left of the other player")
	}

	if !gameboard.UpdatePlayer(Red, Right, 1, false) {
		t.Error("first player should not be able to move on the second player")
	}
}

func TestPlayerShouldBeMarkedAsCrashedWhenCrashes(t *testing.T) {
	gameboard := NewBoard(Red, Blue)

	if !gameboard.UpdatePlayer(Red, Right, 1000, false) {
		t.Error("first player should crash when moving 1000 positions to the right")
	}

	if !gameboard.playerStates[Red].status.crashed {
		t.Error("first player should be marked as crashed when crashing")
	}
}

func TestPlayerShouldHaveCountdownReducedAfterCrashing(t *testing.T) {
	gameboard := NewBoard(Red, Blue)

	if !gameboard.UpdatePlayer(Red, Right, 1000, false) {
		t.Error("first player should crash when moving 1000 positions to the right")
	}

	if gameboard.playerStates[Red].status.animationCountdown != 4 {
		t.Error("first player should have 4 as the initial countdown")
	}

	gameboard.UpdatePlayer(Red, Right, 1, false)
	gameboard.UpdatePlayer(Red, Right, 1, false)
	gameboard.UpdatePlayer(Red, Right, 1, false)
	gameboard.UpdatePlayer(Red, Right, 1, false)
	gameboard.UpdatePlayer(Red, Right, 1, false)

	countdown := gameboard.playerStates[Red].status.animationCountdown
	removed := gameboard.playerStates[Red].status.removed
	if countdown != -1 || !removed {
		t.Errorf("first player should be removed and have the countdown to -1 after animation stops. Countdown %d status %t", countdown, removed)
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
			false,
		},
		{
			"should crash when going over the top bound",
			Red,
			Up,
			21,
			true,
		},
		{
			"should be able to go before the bottom bound",
			Red,
			Down,
			19,
			false,
		},
		{
			"should crash when going over the bottom bound",
			Red,
			Down,
			20,
			true,
		},
		{
			"should be able to go before the left bound",
			Red,
			Left,
			9,
			false,
		},
		{
			"should crash when going over the left bound",
			Red,
			Left,
			10,
			true,
		},
		{
			"should be able to go before the right bound",
			Blue,
			Right,
			9,
			false,
		},
		{
			"should crash when going over the right bound",
			Blue,
			Right,
			10,
			true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gameboard := NewBoard(Red, Blue)

			if gameboard.UpdatePlayer(tt.color, tt.direction, tt.speed, false) != tt.want {
				t.Errorf("Expected UpdatePlayer to return %t", tt.want)
			}
		})
	}
}
