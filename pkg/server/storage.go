package server

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
)

type GameStorage struct {
	mx    sync.RWMutex
	games map[string]ai.AIGame
}

func (gs *GameStorage) Load(msgID string) (*ai.AIGame, error) {
	gs.mx.RLock()
	defer gs.mx.RUnlock()
	g, ok := gs.games[msgID]
	if !ok {
		return nil, errors.New("game not found")
	}
	return &g, nil
}

func (gs *GameStorage) Save(msgID string, g *ai.AIGame) error {
	gs.mx.Lock()
	defer gs.mx.Unlock()
	gs.games[msgID] = *g
	return nil
}

func NewGameStorage() *GameStorage {
	return &GameStorage{
		games: make(map[string]ai.AIGame),
	}
}
