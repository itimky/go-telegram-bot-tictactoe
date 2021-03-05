package algorithms

import (
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Load_Rotated0(t *testing.T) {
	c := NewScoreCache()

	board := game.Board{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 0},
	}

	g := game.ContinueGame(game.MarkX, board, 3, 3)
	c.Store(g, 1000)

	rotatedBoard := game.Board{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 0},
	}

	gRotated := game.ContinueGame(game.MarkX, rotatedBoard, 3, 3)
	score, ok := c.Load(gRotated)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1000., score)
}

func Test_Load_Rotated90(t *testing.T) {
	c := NewScoreCache()

	board := game.Board{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 0},
	}

	g := game.ContinueGame(game.MarkX, board, 3, 3)
	c.Store(g, 1000)

	rotatedBoard := game.Board{
		{0, 0, 0},
		{0, 2, 0},
		{1, 0, 0},
	}

	gRotated := game.ContinueGame(game.MarkX, rotatedBoard, 3, 3)
	score, ok := c.Load(gRotated)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1000., score)
}

func Test_Load_Rotated180(t *testing.T) {
	c := NewScoreCache()

	board := game.Board{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 0},
	}

	g := game.ContinueGame(game.MarkX, board, 3, 3)
	c.Store(g, 1000)

	rotatedBoard := game.Board{
		{0, 0, 0},
		{0, 2, 0},
		{0, 0, 1},
	}

	gRotated := game.ContinueGame(game.MarkX, rotatedBoard, 3, 3)
	score, ok := c.Load(gRotated)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1000., score)
}

func Test_Load_Rotated270(t *testing.T) {
	c := NewScoreCache()

	board := game.Board{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 0},
	}

	g := game.ContinueGame(game.MarkX, board, 3, 3)
	c.Store(g, 1000)

	rotatedBoard := game.Board{
		{0, 0, 1},
		{0, 2, 0},
		{0, 0, 0},
	}

	gRotated := game.ContinueGame(game.MarkX, rotatedBoard, 3, 3)
	score, ok := c.Load(gRotated)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1000., score)
}
