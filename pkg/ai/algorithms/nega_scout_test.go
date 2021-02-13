package algorithms

import (
	"testing"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
	"github.com/stretchr/testify/assert"
)

func Benchmark_NegaScout(b *testing.B) {
	g := game.NewGame(game.MarkX)
	_, err := GetNextMoveNegaScout(g)
	assert.NoError(b, err)
}
