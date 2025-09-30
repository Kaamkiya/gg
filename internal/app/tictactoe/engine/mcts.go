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
	// Returns the best move for the current player
	Solve(board *Board) int
}

type GameEngine interface {
	// Returns gameover (bool) & a value if there's a winner
	CheckGameOver(board *Board, lastMove int) (bool, int)
	// Get all available moves
	GetLegalMoves(board *Board) []int
	// Get the opponent of a player
	GetOpponent(player int) int
	// Play a move on the board
	PlayMove(board *Board, player int, move int) error
}

type mcts struct {
	engine GameEngine
	depth  int
}

func NewMCTS(engine GameEngine, depth int) AI {
	return &mcts{engine, depth}
}

func (m *mcts) Solve(board *Board) int {
	root := newNode(m.engine, board, -1, nil)

	for i := 0; i < m.depth; i++ {
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

type node struct {
	engine     GameEngine
	board      *Board
	move       int
	parent     *node
	children   []*node
	legalMoves []int
	valueSum   int
	visitCount int
}

func newNode(engine GameEngine, board *Board, move int, parent *node) *node {
	legalMoves := engine.GetLegalMoves(board)

	return &node{
		engine:     engine,
		board:      board,
		move:       move,
		parent:     parent,
		children:   []*node{},
		legalMoves: legalMoves,
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
		move, _, err := popRandomMove(n.engine.GetLegalMoves(board))
		if err != nil {
			break
		}

		n.engine.PlayMove(board, player, move)
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
	move, rest, err := popRandomMove(n.legalMoves)
	if err != nil {
		return nil, err
	}

	n.legalMoves = rest

	board := n.board.Copy()
	n.engine.PlayMove(board, P1, move)

	// Every node considers itself as p1
	board.ChangePerspective()
	child := newNode(n.engine, board, move, n)
	n.children = append(n.children, child)

	return child, nil
}

func (n *node) backpropagate(value int) {
	n.visitCount++
	n.valueSum += value

	if n.parent != nil {
		n.parent.backpropagate(n.engine.GetOpponent(value))
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

func popRandomMove(legalMoves []int) (int, []int, error) {
	if len(legalMoves) == 0 {
		return -1, legalMoves, fmt.Errorf("No legal moves")
	}

	index := rand.IntN(len(legalMoves))
	move := legalMoves[index]
	legalMoves = append(legalMoves[:index], legalMoves[index+1:]...)

	return move, legalMoves, nil
}

func (n *node) isExpanded() bool {
	return len(n.children) > 0 && len(n.legalMoves) == 0
}

func (n *node) getUCB(child *node) float64 {
	q := 1 - ((float64(child.valueSum)/float64(child.visitCount))+1)/2
	return q + C_VALUE*math.Sqrt(math.Log(float64(n.visitCount))/float64(child.visitCount))
}
