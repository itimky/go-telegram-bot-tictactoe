package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RotateLeft(t *testing.T) {
	b := Board{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	rb := Board{
		{3, 6, 9},
		{2, 5, 8},
		{1, 4, 7},
	}

	assert.Equal(t, rb, b.RotateLeft())
}

func BenchmarkBoard_RotateLeft(b *testing.B) {
	board := Board{}
	board.RotateLeft()
}
