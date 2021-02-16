package algorithms

import (
	"math/rand"
	"time"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

type Random struct {
	random *rand.Rand
}

func NewRandom() Random {
	seed := rand.NewSource(time.Now().Unix())
	return Random{random: rand.New(seed)}
}

func (r Random) GetNextMove(g game.Game) (game.Move, error) {
	moves := g.GetPossibleMoves()
	randomMove := moves[r.random.Intn(len(moves))]
	return randomMove, nil
}
