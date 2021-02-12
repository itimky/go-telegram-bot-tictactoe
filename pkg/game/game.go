package game

import (
	"math"

	"github.com/pkg/errors"
)

type Coordinates struct {
	I int
	J int
}

type Game struct {
	board    Board
	player   Mark
	opponent Mark
}

func NewGame(player Mark, board Board) Game {
	game := Game{
		player:   player,
		opponent: getOpponent(player),
		board:    board,
	}
	return game
}

func StartNewGame(player Mark) (Game, error) {
	game := Game{
		player:   player,
		opponent: getOpponent(player),
	}
	if game.isAIFirst() {
		game.SwapPlayers()
		if err := game.MakeAIMove(); err != nil {
			return game, errors.Wrap(err, "failed to start game")
		}
		game.SwapPlayers()
	}

	return game, nil
}

func getOpponent(player Mark) Mark {
	return MarkX ^ MarkO ^ player
}

func (g *Game) PlayRound(coords Coordinates) error {
	if err := g.MakeMove(coords); err != nil {
		return errors.Wrap(err, "failed to make move")
	}
	if !g.IsOver() {
		g.SwapPlayers()
		if err := g.MakeAIMove(); err != nil {
			return errors.Wrap(err, "failed to make ai move")
		}
		g.SwapPlayers()
	}
	return nil
}

func (g *Game) GetBoard() Board {
	return g.board
}

func (g *Game) GetPlayer() Mark {
	return g.player
}

func (g *Game) GetPossibleMoves() []Coordinates {
	moves := make([]Coordinates, 0, g.board.Size())
	for i := range g.board {
		for j := range g.board[i] {
			if g.board[i][j] == MarkEmpty {
				moves = append(moves, Coordinates{I: i, J: j})
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
	for _, line := range g.board.GetLines() {
		score += g.getLineScore(line)
	}
	return score
}

func (g *Game) SwapPlayers() {
	g.player, g.opponent = g.opponent, g.player
}

func (g *Game) isAIFirst() bool {
	return g.player == MarkO
}

func (g *Game) MakeAIMove() error {
	move, err := GetAINextMove(*g)
	if err != nil {
		return errors.Wrap(err, "failed to get AI next move")
	}

	if err = g.MakeMove(move); err != nil {
		return errors.Wrap(err, "failed to make move")
	}
	return nil
}

func (g *Game) MakeMove(c Coordinates) error {
	if err := g.board.PlaceMark(c.I, c.J, g.player); err != nil {
		return errors.Wrap(err, "failed to place mark")
	}
	return nil
}

func (g *Game) IsOver() bool {
	moves := g.GetPossibleMoves()
	if len(moves) == 0 {
		return true
	}
	for _, line := range g.board.GetLines() {
		if math.Abs(g.getLineScore(line)) == 100 {
			return true
		}
	}
	return false
}
