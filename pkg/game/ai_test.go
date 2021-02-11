package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Benchmark_NegaScout(b *testing.B) {
	game := Game{player: MarkX, opponent: MarkO}
	_, err := GetAINextMove(game)
	assert.NoError(b, err)
}
