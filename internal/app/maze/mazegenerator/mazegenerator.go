package mazegenerator

func GenerateMaze(width, height int, algorithm string) *Maze {
	maze := NewMaze(width, height)
	generator := NewMazeGenerator(algorithm)
	generator.Generate(maze)

	return maze
}
