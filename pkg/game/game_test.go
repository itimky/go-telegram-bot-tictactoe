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

func Test_RotateLeft_4x4(t *testing.T) {
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

	resultBoard := board{
		{2, 0, 0, 1},
		{0, 2, 1, 0},
		{2, 1, 2, 2},
		{1, 0, 1, 2},
	}
	game.RotateLeft()

	assert.Equal(t, resultBoard, game.board)
}

func Test_RotateLeft_5x5(t *testing.T) {
	game := Game{
		n:       5,
		lineLen: 4,
		board: board{
			{1, 2, 0, 2, 0},
			{0, 1, 2, 0, 1},
			{1, 2, 1, 0, 1},
			{2, 2, 0, 1, 2},
			{0, 0, 1, 2, 1},
		},
	}

	resultBoard := board{
		{0, 1, 1, 2, 1},
		{2, 0, 0, 1, 2},
		{0, 2, 1, 0, 1},
		{2, 1, 2, 2, 0},
		{1, 0, 1, 2, 0},
	}
	game.RotateLeft()

	log.Debug(game.board)

	assert.Equal(t, resultBoard, game.board)
}

func BenchmarkBoard_RotateLeft(b *testing.B) {
	game := Game{}
	game.RotateLeft()
}
