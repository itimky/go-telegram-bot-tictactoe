package game

import (
	"math"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const initialDepth int = 8

type MoveCache struct {
	mx    sync.RWMutex
	cache map[Game]Coordinates
}

func (mc *MoveCache) Store(g Game, c Coordinates) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	mc.cache[g] = c
}

func (mc *MoveCache) Load(g Game) (Coordinates, bool) {
	mc.mx.RLock()
	defer mc.mx.RUnlock()
	c, ok := mc.cache[g]
	return c, ok
}

func NewMoveCache() *MoveCache {
	return &MoveCache{
		cache: make(map[Game]Coordinates),
	}
}

var nextMoveCache = NewMoveCache()

func NegaScout(game Game, depth int, alpha, beta float64) (float64, error) {
	if depth == 0 || game.IsOver() {
		return game.GetScore(), nil
	}

	bestValue := math.Inf(-1)
	for _, move := range game.GetPossibleMoves() {
		possibleGame := game
		if err := possibleGame.MakeMove(move); err != nil {
			return bestValue, errors.Wrap(err, "failed to calc next move")
		}
		possibleGame.SwapPlayers()

		moveAlpha, err := NegaScout(possibleGame, depth-1, -beta, -alpha)
		if err != nil {
			return bestValue, errors.Wrap(err, "failed to calc further steps")
		}
		moveAlpha = -moveAlpha
		bestValue = math.Max(bestValue, moveAlpha)
		if alpha < moveAlpha {
			alpha = moveAlpha
			if alpha >= beta {
				break
			}
		}
	}

	return bestValue, nil
}

func GetAINextMove(game Game) (Coordinates, error) {
	if move, ok := nextMoveCache.Load(game); ok {
		log.Info("Move cache used")
		return move, nil
	}

	possibleMoves := game.GetPossibleMoves()
	resultMove := possibleMoves[0]
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	depth := initialDepth

	for _, move := range possibleMoves {
		possibleGame := game
		if err := possibleGame.MakeMove(move); err != nil {
			return resultMove, errors.Wrap(err, "failed to calc next move")
		}
		possibleGame.SwapPlayers()
		moveAlpha, err := NegaScout(possibleGame, depth-1, -beta, -alpha)
		if err != nil {
			return resultMove, errors.Wrap(err, "failed to calc further steps")
		}
		moveAlpha = -moveAlpha
		if alpha < moveAlpha {
			alpha = moveAlpha
			resultMove = move
		}
	}

	nextMoveCache.Store(game, resultMove)
	return resultMove, nil
}
