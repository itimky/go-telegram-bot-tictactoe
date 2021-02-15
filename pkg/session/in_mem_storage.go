package session

import (
	"sync"

	"github.com/pkg/errors"
)

type InMemStorage struct {
	mx    sync.RWMutex
	games map[int]session
}

func (gs *InMemStorage) Load(id int) (*session, error) {
	gs.mx.RLock()
	defer gs.mx.RUnlock()
	g, ok := gs.games[id]
	if !ok {
		return nil, errors.New("game not found")
	}
	return &g, nil
}

func (gs *InMemStorage) Save(session *session) error {
	gs.mx.Lock()
	defer gs.mx.Unlock()
	gs.games[session.id] = *session
	return nil
}

func NewInMemStorage() *InMemStorage {
	return &InMemStorage{
		games: make(map[int]session),
	}
}
