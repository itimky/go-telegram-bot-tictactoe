package server

import (
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BoardToString(t *testing.T) {
	var b game.Board
	assert.Equal(t, "000000000", BoardToString(b))
}
