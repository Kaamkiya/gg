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
	start := Cell{startX, startY}
	curr := start

	walls := maze.GetFrontiers(startX, startY, true)
	visited := make(map[Cell]bool)
	for _, wall := range walls {
		visited[wall] = true
	}

	for len(walls) > 0 {
		// Pop random wall
		randIdx := rand.Intn(len(walls))
		wall := walls[randIdx]
		walls = append(walls[:randIdx], walls[randIdx+1:]...)

		if maze.Get(wall.x, wall.y) == PATH {
			continue
		}

		paths := maze.GetFrontiers(wall.x, wall.y, false)
		if len(paths) == 0 {
			continue
		}
		path := paths[rand.Intn(len(paths))]

		// skip special case: last wall before boundary
		if wall.Diff(path) != 1 {
			// Connect wall and path
			x, y := (wall.x+path.x)/2, (wall.y+path.y)/2
			between := Cell{x, y}
			maze.MakePath(between)
		}

		maze.MakePath(wall)
		// Add walls
		neighbors := maze.GetFrontiers(wall.x, wall.y, true)
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				walls = append(walls, neighbor)
			}
		}

		// find the longest point
		if maze.IsBoundary(wall.x, wall.y) && wall.Diff(start) > curr.Diff(start) {
			curr = wall
		}
	}

	maze.SetEnd(curr.x, curr.y)
}
