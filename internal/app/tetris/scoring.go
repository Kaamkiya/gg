package tetris

func (gs *gameState) addLineScore(completedLinesNum int) {
	switch completedLinesNum {
	case 4:
		gs.scorePoints(800)
	case 3:
		gs.scorePoints(500)
	case 2:
		gs.scorePoints(300)
	default:
		gs.scorePoints(100)
	}
}

func (gs *gameState) addStillLivingScore() {
	gs.scorePoints(1)
}

func (gs *gameState) addLivingDangerouslyScore() {
	gs.scorePoints(2)
}

func (gs *gameState) scorePoints(points uint) {
	gs.score += uint(float32(points) * gs.currentDifficulty.level)
}
