package algorithms

import (
	"math"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

const initialDepth byte = 8

type MoveCache struct {
	mx    sync.RWMutex
	cache map[game.Game]game.Move
}

func (mc *MoveCache) Store(g game.Game, c game.Move) {
	mc.mx.Lock()
	defer mc.mx.Unlock()
	mc.cache[g] = c
}

func (mc *MoveCache) Load(g game.Game) (game.Move, bool) {
	mc.mx.RLock()
	defer mc.mx.RUnlock()
	c, ok := mc.cache[g]
	return c, ok
}

func NewMoveCache() *MoveCache {
	return &MoveCache{
		cache: make(map[game.Game]game.Move),
	}
}

type NegaScout struct {
	initialDepth  byte
	nextMoveCache *MoveCache
}

func NewNegaScout() NegaScout {
	return NegaScout{
		initialDepth:  initialDepth,
		nextMoveCache: NewMoveCache(),
	}
}

func (ns NegaScout) GetNextMove(g game.Game) (game.Move, error) {
	start := time.Now()
	move, err := ns.getNextMove(g)
	log.Debug("NegaScout.GetNextMove elapsed: ", time.Since(start).Seconds())
	return move, err
}

func (ns NegaScout) getNextMove(g game.Game) (game.Move, error) {
	if move, ok := ns.nextMoveCache.Load(g); ok {
		log.Info("Move cache used")
		return move, nil
	}

	possibleMoves := g.GetPossibleMoves()
	resultMove := possibleMoves[0]
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	depth := initialDepth

	for _, move := range possibleMoves {
		possibleGame := g
		if err := possibleGame.MakeMove(move); err != nil {
			return resultMove, errors.Wrap(err, "failed to calc next move")
		}
		possibleGame.SwapPlayers()
		moveAlpha, err := ns.getBestScoreRecursive(possibleGame, depth-1, -beta, -alpha)
		if err != nil {
			return resultMove, errors.Wrap(err, "failed to calc further steps")
		}
		moveAlpha = -moveAlpha
		if alpha < moveAlpha {
			alpha = moveAlpha
			resultMove = move
		}
	}

	ns.nextMoveCache.Store(g, resultMove)
	return resultMove, nil
}

func (ns NegaScout) getBestScoreRecursive(g game.Game, depth byte, alpha, beta float64) (float64, error) {
	if depth == 0 || g.IsOver() {
		return g.GetScore(), nil
	}

	bestValue := math.Inf(-1)
	for _, move := range g.GetPossibleMoves() {
		possibleGame := g
		if err := possibleGame.MakeMove(move); err != nil {
			return bestValue, errors.Wrap(err, "failed to calc next move")
		}
		possibleGame.SwapPlayers()

		moveAlpha, err := ns.getBestScoreRecursive(possibleGame, depth-1, -beta, -alpha)
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
