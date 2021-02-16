package game

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func Test_GetLines(t *testing.T) {
	game := Game{
		n:       4,
		lineLen: 4,
		board: board{
			{1, 2, 0, 2},
			{0, 1, 2, 0},
			{1, 2, 1, 0},
			{2, 2, 0, 1},
		},
	}

	lines := game.GetLines()
	expected := []Line{
		{1, 2, 0, 2},
		{0, 1, 2, 0},
		{1, 2, 1, 0},
		{2, 2, 0, 1},

		{1, 0, 1, 2},
		{2, 1, 2, 2},
		{0, 2, 1, 0},
		{2, 0, 0, 1},

		{1, 1, 1, 1},
		{2, 2, 2, 2},
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

func Test_IsOver(t *testing.T) {
	game := Game{
		n:        4,
		lineLen:  4,
		player:   MarkX,
		opponent: MarkO,
		board: board{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{1, 1, 1, 1},
		},
	}
	assert.Equal(t, true, game.IsOver())
}
