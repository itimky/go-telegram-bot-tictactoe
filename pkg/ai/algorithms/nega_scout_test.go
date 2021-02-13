package algorithms

import (
	"testing"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
	"github.com/stretchr/testify/assert"
)

func Benchmark_NegaScout(b *testing.B) {
	g := game.NewGame(game.MarkX, game.Board{})
	_, err := GetAINextMove(g)
	assert.NoError(b, err)
}
