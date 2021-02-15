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

type Line [3]Mark
type board [10][10]Mark
type Board [][]Mark

type Game struct {
	board    board
	n        int
	player   Mark
	opponent Mark
}

func NewGame(player Mark, size int) Game {
	return Game{
		player:   player,
		n:        size,
		opponent: getOpponent(player),
	}
}

func ContinueGame(player Mark, brd Board, size int) Game {
	b := board{}

	for i, row := range brd {
		for j, mark := range row {
			b[i][j] = mark
		}
	}

	return Game{
		player:   player,
		n:        size,
		opponent: getOpponent(player),
		board:    b,
	}
}

func (g *Game) GetLines() []Line {
	lines := make([]Line, 0, g.LineCount())
	n := g.Size()
	for _, row := range g.board {
		for j := 0; j < n-2; j++ {
			lines = append(lines, Line{row[j], row[j+1], row[j+2]})
		}
	}

	for j := 0; j < n; j++ {
		for i := 0; i < n-2; i++ {
			lines = append(lines, Line{g.board[i][j], g.board[i+1][j], g.board[i+2][j]})
		}
	}

	for i := 0; i < n-2; i++ {
		for j := 0; j < n-2; j++ {
			lines = append(lines, Line{g.board[i][j], g.board[i+1][j+1], g.board[i+2][j+2]})
		}
	}

	for i := 0; i < n-2; i++ {
		for j := 2; j < n; j++ {
			lines = append(lines, Line{g.board[i][j], g.board[i+1][j-1], g.board[i+2][j-2]})
		}
	}
	return lines
}

func (g *Game) LineCount() int {
	/*
		In the future board size could be > 3, but win condition would be the same:  line with 3 same marks in a row.
		Rows, columns and diagonals with len = 3

		• • • •		• • • •		• * • •		• • • •
		* * * •		• * • •		• • * •		• • • *
		• • • •		• * • •		• • • *		• • * •
		• • • •		• * • •		• • • •		• * • •
	*/
	rowsAndCols := 2 * g.n * (g.n - 2)
	diagonals := 2 * (g.n - 2) * (g.n - 2)
	return rowsAndCols + diagonals
}

func (g *Game) SquareCount() int {
	return g.n * g.n
}

func (g *Game) Size() int {
	return g.n
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
	opponentPoints := 0.
	playerPoints := 0.
	var prevValue Mark
	for _, value := range line {
		if value == g.opponent {
			if value == prevValue {
				opponentPoints *= 10
			} else {
				opponentPoints += 1
			}
		} else if value == g.player {
			if value == prevValue {
				playerPoints *= 10
			} else {
				playerPoints += 1
			}
		}
		prevValue = value
	}
	return playerPoints - opponentPoints
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
	if len(moves) == 0 {
		return true
	}
	for _, line := range g.GetLines() {
		if math.Abs(g.getLineScore(line)) == 100 {
			return true
		}
	}
	return false
}
