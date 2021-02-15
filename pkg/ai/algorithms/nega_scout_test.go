package algorithms

import (
	"testing"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
	"github.com/stretchr/testify/assert"
)

func Benchmark_NegaScout(b *testing.B) {
	g := game.NewGame(game.MarkX, 3)
	negaScout := NewNegaScout()
	_, err := negaScout.GetNextMove(g)
	assert.NoError(b, err)
}
