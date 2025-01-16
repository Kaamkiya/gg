package engine

import (
	"fmt"
	"math"
	"math/rand/v2"
)

const (
	C_VALUE = 1.41
	DEPTH   = 100
)

type AI interface {
	Solve(board *Board) int
}

type mcts struct {
	engine *Engine
}

func NewMCTS(engine *Engine) AI {
	return &mcts{engine}
}

func (m *mcts) Solve(board *Board) int {
	root := newNode(m.engine, board, -1, nil)

	for i := 0; i < DEPTH; i++ {
		node := root
		for node.isExpanded() {
			child, err := node.selectChild()
			if err != nil {
				panic(err)
			}
			node = child
		}

		isOver, value := m.engine.CheckGameOver(node.board, node.move)
		value = m.engine.GetOpponent(value)

		if !isOver {
			child, err := node.expand()
			if err != nil {
			} else {
				value = child.simulate()
				node = child
			}
		}

		node.backpropagate(value)
	}

	visits := make([]float64, board.Size*board.Size)
	dist := make([]float64, board.Size*board.Size)
	sum := 0.0

	for _, child := range root.children {
		visits[child.move] = float64(child.visitCount)
		sum += visits[child.move]
	}

	for i, visit := range visits {
		dist[i] = visit / sum
	}

	bestMove := -1
	bestValue := 0.0

	for i, value := range dist {
		if value > bestValue {
			bestMove = i
			bestValue = value
		}
	}

	return bestMove
}

// Represents a game node in mcts tree
type node struct {
	engine     *Engine
	board      *Board
	move       int
	parent     *node
	children   []*node
	legalMoves map[int]bool
	valueSum   int
	visitCount int
}

// Create a new node
func newNode(engine *Engine, board *Board, move int, parent *node) *node {
	legalMoves := engine.GetLegalMoves(board)
	moves := make(map[int]bool, len(legalMoves))
	for _, m := range legalMoves {
		moves[m] = true
	}

	return &node{
		engine:     engine,
		board:      board,
		move:       move,
		parent:     parent,
		children:   []*node{},
		legalMoves: moves,
		valueSum:   0,
		visitCount: 0,
	}
}

// Simulate all moves until game is over;
// Returns winner
func (n *node) simulate() int {
	isOver, winner := n.engine.CheckGameOver(n.board, n.move)
	if isOver {
		return n.engine.GetOpponent(winner)
	}

	board := n.board.Copy()
	player := P1
	result := 0

	for {
		legalMoves := n.engine.GetLegalMoves(board)
		moves := make(map[int]bool, len(legalMoves))
		for _, m := range legalMoves {
			moves[m] = true
		}
		move, err := popRandomMove(moves)
		if err != nil {
			break
		}

		board.SetCell(move, player)
		isOver, winner = n.engine.CheckGameOver(board, move)
		if isOver {
			result = winner
			break
		}

		player = n.engine.GetOpponent(player)
	}

	return result
}

func (n *node) expand() (*node, error) {
	move, err := popRandomMove(n.legalMoves)
	if err != nil {
		return nil, err
	}

	n.legalMoves[move] = false

	board := n.board.Copy()

	// Every node considers itself as p1
	board.SetCell(move, P1)
	board.ChangePerspective()

	child := newNode(n.engine, board, move, n)
	n.children = append(n.children, child)

	return child, nil
}

func (n *node) backpropagate(value int) {
	n.visitCount++
	n.valueSum += value

	if n.parent != nil {
		n.parent.backpropagate(value * -1)
	}
}

// Get next child with highest UCB
func (n *node) selectChild() (*node, error) {
	if len(n.children) == 0 {
		return nil, fmt.Errorf("No child nodes")
	}

	var selected *node
	var bestValue float64 = math.Inf(-1)

	for _, child := range n.children {
		ucb := n.getUCB(child)
		if selected == nil || ucb > bestValue {
			selected = child
			bestValue = ucb
		}
	}

	return selected, nil
}

func popRandomMove(moves map[int]bool) (int, error) {
	legalMoves := []int{}
	for m, v := range moves {
		if v {
			legalMoves = append(legalMoves, m)
		}
	}

	if len(legalMoves) == 0 {
		return -1, fmt.Errorf("No legal moves")
	}

	index := rand.IntN(len(legalMoves))
	move := legalMoves[index]

	return move, nil
}

func (n *node) isExpanded() bool {
	allVisited := true
	for _, m := range n.legalMoves {
		if m {
			allVisited = false
			break
		}
	}

	return len(n.children) > 0 && allVisited
}

func (n *node) getUCB(child *node) float64 {
	q := 1 - ((float64(child.valueSum)/float64(child.visitCount))+1)/2
	return q + C_VALUE*math.Sqrt(math.Log(float64(n.visitCount))/float64(child.visitCount))
}
