package storage

import "github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"

type IGameStorage interface {
	Load(msgID int) (ai.AIGame, error)
	Save(msgID int, g ai.AIGame) error
}
