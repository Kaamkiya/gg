package tetris

import "time"

const (
	// initialDifficulyCountDown is the number of pieces that trigger a difficulty increase
	initialDifficulyCountDown = 10
	// initialDifficulyLevel is the factor that increases scoring and decreases the game tick. Increased by 0.1 on difficulty increase.
	initialDifficulyLevel = 1.0
)

type difficulty struct {
	countdown             int
	level                 float32
	gameProgressTickDelay time.Duration
}

func (gs *gameState) adjustDifficulty() {
	if gs.currentDifficulty.countdown <= 1 {
		gs.currentDifficulty.countdown = initialDifficulyCountDown
		gs.currentDifficulty.level += 0.1
		gs.currentDifficulty.gameProgressTickDelay = time.Duration(float32(initialGameProgressTickDelay) / gs.currentDifficulty.level)

		return
	}

	gs.currentDifficulty.countdown--
}
