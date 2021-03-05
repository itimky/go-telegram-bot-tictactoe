package algorithms

import (
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

const NegaScoutMaxDepth byte = 10

type NegaScout struct {
	initialDepth  byte
	nextMoveCache *MoveCache
	scoreCache    *ScoreCache
}

func NewNegaScout(depth byte) NegaScout {
	return NegaScout{
		initialDepth:  depth,
		nextMoveCache: NewMoveCache(),
		scoreCache:    NewScoreCache(),
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
		log.Debug("Move cache used")
		return move, nil
	} else {
		log.Debug("Move cache miss")
	}

	possibleMoves := g.GetPossibleMoves()
	resultMove := possibleMoves[0]
	alpha := math.Inf(-1)
	beta := math.Inf(1)

	for _, move := range possibleMoves {
		possibleGame := g
		if err := possibleGame.MakeMove(move); err != nil {
			return resultMove, err
		}
		possibleGame.SwapPlayers()
		moveAlpha, err := ns.getBestScoreRecursive(possibleGame, ns.initialDepth-1, -beta, -alpha)
		if err != nil {
			return resultMove, err
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
	if score, ok := ns.scoreCache.Load(g); ok {
		return score, nil
	}

	if depth == 0 || g.IsOver() {
		score := g.GetScore()
		ns.scoreCache.Store(g, score)
		return score, nil
	}

	bestValue := math.Inf(-1)
	for _, move := range g.GetPossibleMoves() {
		possibleGame := g
		if err := possibleGame.MakeMove(move); err != nil {
			return bestValue, err
		}
		possibleGame.SwapPlayers()

		moveAlpha, err := ns.getBestScoreRecursive(possibleGame, depth-1, -beta, -alpha)
		if err != nil {
			return bestValue, err
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

	ns.scoreCache.Store(g, bestValue)
	return bestValue, nil
}
