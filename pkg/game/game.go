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
type Board [3]Line

type Game struct {
	board    Board
	n        int
	player   Mark
	opponent Mark
}

func NewGame(player Mark) Game {
	return Game{
		player:   player,
		opponent: getOpponent(player),
	}
}

func ContinueGame(player Mark, board Board) Game {
	return Game{
		player:   player,
		opponent: getOpponent(player),
		board:    board,
	}
}

func (g *Game) GetLines() []Line {
	lines := make([]Line, 0, g.LineCount())
	for _, row := range g.board {
		lines = append(lines, row)
	}

	for j := range g.board {
		var column Line
		for i, row := range g.board {
			column[i] = row[j]
		}
		lines = append(lines, column)
	}

	var diagonal Line
	for i := range g.GetBoard() {
		diagonal[i] = g.board[i][i]
	}
	lines = append(lines, diagonal)

	var counterDiagonal Line
	for i := range g.board {
		counterDiagonal[i] = g.board[i][len(g.board)-i-1]
	}
	lines = append(lines, counterDiagonal)
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

func (g *Game) Size() int {
	return g.n * g.n
}

func getOpponent(player Mark) Mark {
	return MarkX ^ MarkO ^ player
}

func (g *Game) GetBoard() Board {
	return g.board
}

func (g *Game) GetPlayer() Mark {
	return g.player
}

func (g *Game) GetPossibleMoves() []Move {
	moves := make([]Move, 0, g.Size())
	for i := range g.board {
		for j := range g.board[i] {
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
