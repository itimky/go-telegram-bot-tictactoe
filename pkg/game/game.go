package game

import (
	"math"

	"github.com/pkg/errors"
)

type Move struct {
	I int
	J int
}

type Mark byte

const (
	MarkEmpty Mark = 0
	MarkX     Mark = 1
	MarkO     Mark = 2
)

type Line []Mark
type board [10][10]Mark

type Board [][]Mark

type Game struct {
	board    board
	n        int
	lineLen  int
	player   Mark
	opponent Mark
}

func NewGame(player Mark, size int, lineLen int) Game {
	return Game{
		player:   player,
		n:        size,
		lineLen:  lineLen,
		opponent: getOpponent(player),
	}
}

func ContinueGame(player Mark, brd Board, size int, lineLen int) Game {
	b := board{}

	for i, row := range brd {
		for j, mark := range row {
			b[i][j] = mark
		}
	}

	return Game{
		player:   player,
		n:        size,
		lineLen:  lineLen,
		opponent: getOpponent(player),
		board:    b,
	}
}

func (g *Game) GetLines() []Line {
	// TODO: optimize
	lines := make([]Line, 0, g.LineCount())
	n := g.Size()
	lineLen := g.LineLen()
	linesInRow := n - lineLen + 1
	for i := 0; i < n; i++ {
		for j := 0; j < linesInRow; j++ {
			line := make(Line, 0, lineLen)
			for k := 0; k < lineLen; k++ {
				line = append(line, g.board[i][j+k])
			}
			lines = append(lines, line)
		}
	}

	for j := 0; j < n; j++ {
		for i := 0; i < linesInRow; i++ {
			line := make(Line, 0, lineLen)
			for k := 0; k < lineLen; k++ {
				line = append(line, g.board[i+k][j])
			}
			lines = append(lines, line)
		}
	}

	for i := 0; i < linesInRow; i++ {
		for j := 0; j < linesInRow; j++ {
			line := make(Line, 0, lineLen)
			for k := 0; k < lineLen; k++ {
				line = append(line, g.board[i+k][j+k])
			}
			lines = append(lines, line)
		}
	}

	for i := 0; i < linesInRow; i++ {
		for j := lineLen - 1; j < n; j++ {
			line := make(Line, 0, lineLen)
			for k := 0; k < lineLen; k++ {
				line = append(line, g.board[i+k][j-k])
			}
			lines = append(lines, line)
		}
	}
	return lines
}

func (g *Game) LineCount() int {
	n := g.Size()
	linesInRow := n - g.LineLen() + 1
	rowsAndCols := 2 * n * linesInRow
	diagonals := 2 * linesInRow * linesInRow
	return rowsAndCols + diagonals
}

func (g *Game) SquareCount() int {
	return g.n * g.n
}

func (g *Game) Size() int {
	return g.n
}

func (g *Game) LineLen() int {
	return g.lineLen
}

func getOpponent(player Mark) Mark {
	return MarkX ^ MarkO ^ player
}

func (g *Game) GetBoard() Board {
	n := g.Size()
	brd := make([][]Mark, 0, n)
	for i := 0; i < n; i++ {
		row := make([]Mark, 0, n)
		for j := 0; j < n; j++ {
			row = append(row, g.board[i][j])
		}
		brd = append(brd, row)
	}
	return brd
}

func (g *Game) GetPlayer() Mark {
	return g.player
}

func (g *Game) GetPossibleMoves() []Move {
	moves := make([]Move, 0, g.SquareCount())
	n := g.Size()
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if g.board[i][j] == MarkEmpty {
				moves = append(moves, Move{I: i, J: j})
			}
		}
	}
	return moves
}

func (g *Game) getLineScore(line Line) float64 {
	var opponentPoints float64
	var playerPoints float64
	var prevValue Mark
	for _, value := range line {
		if value == g.player {
			if value == prevValue {
				playerPoints *= 10
			} else {
				playerPoints += 1
			}
		} else {
			if value == prevValue {
				opponentPoints *= 10
			} else {
				opponentPoints += 1
			}
		}
		prevValue = value
	}
	return playerPoints - opponentPoints
}

func (g *Game) WinScore() float64 {
	return math.Pow(10, float64(g.lineLen-1))
}

func (g *Game) GetScore() float64 {
	var score float64
	for _, line := range g.GetLines() {
		score += g.getLineScore(line)
	}
	return score
}

func (g *Game) SwapPlayers() {
	g.player, g.opponent = g.opponent, g.player
}

func (g *Game) MakeMove(move Move) error {
	if g.board[move.I][move.J] != MarkEmpty {
		return errors.Errorf("square (%v; %v) is not empty", move.I, move.J)
	}
	g.board[move.I][move.J] = g.player
	return nil
}

func (g *Game) IsOver() bool {
	moves := g.GetPossibleMoves()
	winScore := g.WinScore()
	if len(moves) == 0 {
		return true
	}
	for _, line := range g.GetLines() {
		if math.Abs(g.getLineScore(line)) == winScore {
			return true
		}
	}
	return false
}
