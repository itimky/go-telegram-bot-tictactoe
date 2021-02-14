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

var difficultyNames = map[Difficulty]string{
	Easy:   "Easy",
	Medium: "Medium",
	Hard:   "Hard",
}

var algoMap = map[Difficulty]algorithms.IAlgorithm{
	Hard: algorithms.NewNegaScout(),
}

func getAlgorithm(dif Difficulty) (algorithms.IAlgorithm, error) {
	algo, ok := algoMap[dif]
	if !ok {
		difName, ok := difficultyNames[dif]
		if !ok {
			return nil, errors.Errorf("unknown difficulty code %v", dif)
		}
		return nil, errors.Errorf("%v difficulty is not supported", difName)
	}
	return algo, nil
}

type AIGame struct {
	Game       game.Game
	Difficulty Difficulty
}

func (aig *AIGame) MakeAIMove() error {
	if aig.Game.IsOver() {
		return nil
	}

	algo, err := getAlgorithm(aig.Difficulty)
	if err != nil {
		return errors.Wrap(err, "failed to get algorithm")
	}
	move, err := algo.GetNextMove(aig.Game)
	if err != nil {
		return errors.Wrap(err, "failed to get next move")
	}

	aig.Game.SwapPlayers()
	if err = aig.Game.MakeMove(move); err != nil {
		return errors.Wrap(err, "failed to make move")
	}
	aig.Game.SwapPlayers()
	return nil
}

func StartAIGame(dif Difficulty, player game.Mark) (AIGame, error) {
	aiGame := AIGame{
		Game:       game.NewGame(player),
		Difficulty: dif,
	}
	if !aiGame.Game.IsPlayerFirst() {
		if err := aiGame.MakeAIMove(); err != nil {
			return aiGame, errors.Wrap(err, "failed to play AI opponent")
		}
	}
	return aiGame, nil
}
