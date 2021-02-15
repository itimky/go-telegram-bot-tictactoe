package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetLines(t *testing.T) {
	game := Game{
		n: 3,
		board: board{
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

func Test_GetPossibleMoves(t *testing.T) {
	game := Game{
		n: 1,
		board: board{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	}
	assert.Equal(t, []Move{{0, 0}}, game.GetPossibleMoves())

}
