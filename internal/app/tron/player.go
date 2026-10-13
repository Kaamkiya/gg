package tron

import "github.com/Kaamkiya/gg/internal/app/tron/gameboard"

type player struct {
	Direction gameboard.Direction
	Color     gameboard.Color
	Crashed   bool
}

func (p *player) changeDirection(direction gameboard.Direction) {
	if direction.IsOpposite(p.Direction) {
		return
	}

	p.Direction = direction
}
