package mazegenerator

import (
	"math/rand"
)

type MazeGenerator interface {
	Generate(maze *Maze)
}

func NewMazeGenerator(generator string) MazeGenerator {
	switch generator {
	case "prim":
		return &PrimGenerator{}
	default:
		return &PrimGenerator{}
	}
}

type PrimGenerator struct {
}

func (p *PrimGenerator) Generate(maze *Maze) {
	startX, startY := maze.GetStartPos()

	walls := maze.GetNeighbors(startX, startY, true)

	curr := Cell{startX, startY}

	for len(walls) > 0 {
		// Choose a random wall
		randIdx := rand.Intn(len(walls))

		// Pop the wall
		wall := walls[randIdx]
		walls = append(walls[:randIdx], walls[randIdx+1:]...)

		if maze.Get(wall.x, wall.y) == PATH {
			continue
		}

		paths := maze.GetNeighbors(wall.x, wall.y, false)
		if len(paths) >= 2 || len(paths) == 0 {
			continue
		}
		path := paths[rand.Intn(len(paths))]

		// Get opposite cell of the path
		nextX, nextY := wall.x+(wall.x-path.x), wall.y+(wall.y-path.y)
		next := Cell{nextX, nextY}
		maze.MakePath(wall)

		if maze.IsInner(next.x, next.y) {
			maze.MakePath(next)
			// Add walls
			walls = append(walls, maze.GetNeighbors(next.x, next.y, true)...)
		} else {
			walls = append(walls, maze.GetNeighbors(wall.x, wall.y, true)...)
			// find the longest path on the boundary
			if next.Diff(curr) > curr.Diff(curr) {
				curr = next
			}
		}

	}

	maze.SetEnd(curr.x, curr.y)
}
