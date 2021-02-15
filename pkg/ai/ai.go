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
	algos           map[Difficulty]algorithms.IAlgorithm
	difficultyNames map[Difficulty]string
}

func NewAI() AI {
	return AI{
		algos: map[Difficulty]algorithms.IAlgorithm{
			Hard: algorithms.NewNegaScout(),
		},
		difficultyNames: map[Difficulty]string{
			Easy:   "Easy",
			Medium: "Medium",
			Hard:   "Hard",
		},
	}
}

func (ai AI) getAlgorithm(dif Difficulty) (algorithms.IAlgorithm, error) {
	algo, ok := ai.algos[dif]
	if !ok {
		difName, ok := ai.difficultyNames[dif]
		if !ok {
			return nil, errors.Errorf("unknown difficulty code %v", dif)
		}
		return nil, errors.Errorf("%v difficulty is not supported", difName)
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
