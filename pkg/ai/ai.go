package ai

import (
	"github.com/pkg/errors"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai/algorithms"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type AI struct {
	algos map[Difficulty]algorithms.IAlgorithm
}

func NewAI() AI {
	return AI{
		algos: map[Difficulty]algorithms.IAlgorithm{
			DifficultyHard: algorithms.NewNegaScout(),
		},
	}
}

func (ai AI) getAlgorithm(dif Difficulty) (algorithms.IAlgorithm, error) {
	algo, ok := ai.algos[dif]
	if !ok {
		return nil, errors.Errorf("%v difficulty is not supported", dif)
	}
	return algo, nil
}

func (ai AI) MakeAIMove(dif Difficulty, g game.Game) (game.Game, error) {
	if g.IsOver() {
		return g, nil
	}

	algo, err := ai.getAlgorithm(dif)
	if err != nil {
		return g, errors.Wrap(err, "failed to get algorithm")
	}
	move, err := algo.GetNextMove(g)
	if err != nil {
		return g, errors.Wrap(err, "failed to get next move")
	}
	if err = g.MakeMove(move); err != nil {
		return g, errors.Wrap(err, "failed to make move")
	}
	return g, nil
}
