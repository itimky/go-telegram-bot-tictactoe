package ai

import (
	"github.com/pkg/errors"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai/algorithms"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

type Difficulty string

const (
	DifficultyNovice Difficulty = "novice"
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
	DifficultyUnfair Difficulty = "unfair"
)

type AI struct {
	algos map[Difficulty]algorithms.IAlgorithm
}

func maxDepthPercent(percent float32) byte {
	return byte(float32(algorithms.NegaScoutMaxDepth) * percent)
}

func NewAI() AI {
	return AI{
		algos: map[Difficulty]algorithms.IAlgorithm{
			DifficultyNovice: algorithms.NewRandom(),
			DifficultyEasy:   algorithms.NewNegaScout(maxDepthPercent(0.2)),
			DifficultyMedium: algorithms.NewNegaScout(maxDepthPercent(0.3)),
			DifficultyHard:   algorithms.NewNegaScout(maxDepthPercent(0.5)),
			DifficultyUnfair: algorithms.NewNegaScout(algorithms.NegaScoutMaxDepth),
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
		return g, err
	}
	move, err := algo.GetNextMove(g)
	if err != nil {
		return g, err
	}
	if err = g.MakeMove(move); err != nil {
		return g, err
	}
	return g, nil
}
