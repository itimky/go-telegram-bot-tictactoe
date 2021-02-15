package game

import (
	"fmt"
)

type Mark byte

type Line [3]Mark
type Board [3]Line

const (
	MarkEmpty Mark = 0
	MarkX     Mark = 1
	MarkO     Mark = 2
)

func (b *Board) PlaceMark(i, j int, mark Mark) error {
	if b[i][j] != MarkEmpty {
		return fmt.Errorf("square (%v; %v) is not empty", i, j)
	}
	b[i][j] = mark
	return nil
}

func (b *Board) RotateLeft() Board {
	rb := *b
	n := len(rb)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			rb[i][j], rb[j][i] = rb[j][i], rb[i][j]
		}
	}

	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		rb[i], rb[j] = rb[j], rb[i]
	}

	return rb
}

func (b *Board) GetLines() []Line {
	lines := make([]Line, 0, b.LineCount())
	for _, row := range b {
		lines = append(lines, row)
	}

	for j := range b {
		var column Line
		for i, row := range b {
			column[i] = row[j]
		}
		lines = append(lines, column)
	}

	var diagonal Line
	for i := range b {
		diagonal[i] = b[i][i]
	}
	lines = append(lines, diagonal)

	var counterDiagonal Line
	for i := range b {
		counterDiagonal[i] = b[i][len(b)-i-1]
	}
	lines = append(lines, counterDiagonal)
	return lines
}

func (b *Board) LineCount() int {
	n := len(b)
	/*
		In the future board size could be > 3, but win condition would be the same:  line with 3 same marks in a row.
		Rows, columns and diagonals with len = 3

		• • • •		• • • •		• * • •		• • • •
		* * * •		• * • •		• • * •		• • • *
		• • • •		• * • •		• • • *		• • * •
		• • • •		• * • •		• • • •		• * • •
	*/
	rowsAndCols := 2 * n * (n - 2)
	diagonals := 2 * (n - 2) * (n - 2)
	return rowsAndCols + diagonals
}

func (b *Board) Size() int {
	n := len(b)
	return n * n
}
