package algorithms

import (
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
	"sync"
)

type MoveCache struct {
	mx    sync.RWMutex
	cache map[game.Game]game.Move
}

func (mc *MoveCache) Store(g game.Game, move game.Move) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	mc.cache[g] = move
}

func (mc *MoveCache) Load(g game.Game) (game.Move, bool) {
	mc.mx.RLock()
	defer mc.mx.RUnlock()
	move, ok := mc.cache[g]
	return move, ok
}

func NewMoveCache() *MoveCache {
	return &MoveCache{
		cache: make(map[game.Game]game.Move),
	}
}

type ScoreCache struct {
	mx    sync.RWMutex
	cache map[game.Game]float64
}

func NewScoreCache() *ScoreCache {
	scoreCache := &ScoreCache{
		cache: make(map[game.Game]float64),
	}
	return scoreCache
}

func (sc *ScoreCache) Store(g game.Game, score float64) {
	if _, ok := sc.Load(g); ok {
		return
	}
	sc.mx.Lock()
	defer sc.mx.Unlock()
	sc.cache[g] = score
}

func (sc *ScoreCache) Load(g game.Game) (float64, bool) {
	sc.mx.RLock()
	defer sc.mx.RUnlock()
	if score, ok := sc.cache[g]; ok {
		return score, ok
	}

	maxRotations := 3
	for rotationCount := 0; rotationCount < maxRotations; rotationCount++ {
		g.RotateLeft()
		if score, ok := sc.cache[g]; ok {
			return score, ok
		}
	}
	return 0, false
}
