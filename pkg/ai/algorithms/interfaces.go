package algorithms

import "github.com/itimky/go-telegram-bot-tictactoe/pkg/game"

type IAlgorithm interface {
	GetNextMove(game.Game) (game.Move, error)
}
