package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetLines(t *testing.T) {
	game := Game{
		n: 3,
		board: Board{
			{1, 2, 0},
			{0, 1, 2},
			{1, 2, 1},
		},
	}

	lines := game.GetLines()
	expected := []Line{
		{1, 2, 0},
		{0, 1, 2},
		{1, 2, 1},

		{1, 0, 1},
		{2, 1, 2},
		{0, 2, 1},

		{1, 1, 1},
		{0, 1, 1},
	}
	assert.Equal(t, expected, lines)
}
