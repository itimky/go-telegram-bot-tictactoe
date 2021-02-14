package storage

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
)

type GameAppStorage struct {
	mx    sync.RWMutex
	games map[int]ai.AIGame
}

func (gs *GameAppStorage) Load(msgID int) (ai.AIGame, error) {
	gs.mx.RLock()
	defer gs.mx.RUnlock()
	g, ok := gs.games[msgID]
	if !ok {
		return g, errors.New("game not found")
	}
	return g, nil
}

func (gs *GameAppStorage) Save(msgID int, g ai.AIGame) error {
	gs.mx.Lock()
	defer gs.mx.Unlock()
	gs.games[msgID] = g
	return nil
}

func NewGameAppStorage() *GameAppStorage {
	return &GameAppStorage{
		games: make(map[int]ai.AIGame),
	}
}
