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

type AIGame struct {
	Game       game.Game
	Difficulty Difficulty
}

func (aig *AIGame) MakeAIMove() error {
	if aig.Game.IsOver() {
		return nil
	}

	var move game.Move
	var err error
	switch aig.Difficulty {
	case Easy:
		return errors.Errorf("Easy difficulty is not supported")
	case Medium:
		return errors.Errorf("Medium difficulty is not supported")
	case Hard:
		move, err = algorithms.GetNextMoveNegaScout(aig.Game)
	}

	if err != nil {
		return err
	}

	aig.Game.SwapPlayers()
	if err = aig.Game.MakeMove(move); err != nil {
		return err
	}
	aig.Game.SwapPlayers()
	return nil
}

func StartAIGame(dif Difficulty, player game.Mark) (*AIGame, error) {
	aiGame := AIGame{
		Game:       game.NewGame(player),
		Difficulty: dif,
	}
	if !aiGame.Game.IsPlayerFirst() {
		if err := aiGame.MakeAIMove(); err != nil {
			return nil, errors.Wrap(err, "failed to play AI opponent")
		}
	}
	return &aiGame, nil
}
