package ai

import (
	"github.com/pkg/errors"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai/algorithms"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

type Difficulty byte

const (
	Easy   = 0
	Medium = 1
	Hard   = 2
)

type AI struct {
	difficulty Difficulty
}

func (a AI) PlayOpponent(g *game.Game) error {
	if g.IsOver() {
		return nil
	}

	var move game.Coordinates
	var err error
	switch a.difficulty {
	case Easy:
		return errors.Errorf("Easy difficulty is not supported")
	case Medium:
		return errors.Errorf("Medium difficulty is not supported")
	case Hard:
		move, err = algorithms.GetNextMoveNegaScout(*g)
	}

	if err != nil {
		return err
	}

	g.SwapPlayers()
	if err = g.MakeMove(move); err != nil {
		return err
	}
	g.SwapPlayers()
	return nil
}

func NewAI(dif Difficulty) AI {
	return AI{
		difficulty: dif,
	}
}
