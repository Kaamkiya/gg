package tron

import "github.com/Kaamkiya/gg/internal/app/tron/gameboard"

type player struct {
	Direction            gameboard.Direction
	Color                gameboard.Color
	Crashed              bool
	Removed              bool
	Nitros               int
	nitroRoundsCountdown int
}

func newPlayer(color gameboard.Color, direction gameboard.Direction, nitros int) *player {
	return &player{
		Direction: direction,
		Color:     color,
		Nitros:    nitros,
	}
}

func (p *player) changeDirection(direction gameboard.Direction) {
	if direction.IsOpposite(p.Direction) {
		return
	}

	p.Direction = direction
}

func (p *player) nitroBoost() {
	if p.Nitros <= 0 {
		return
	}

	p.Nitros--
	p.nitroRoundsCountdown += roundsPerNitroUsage
}

func (p *player) speedForNextTick() int {
	if p.nitroRoundsCountdown > 0 {
		p.nitroRoundsCountdown--
		return nitroBoost
	}

	return 1
}
