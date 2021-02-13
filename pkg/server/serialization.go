package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

func GetCoordinatesFromString(data string) (game.Coordinates, error) {
	var coords game.Coordinates
	splCoords := strings.Split(data, ",")
	if len(splCoords) != 2 {
		return coords, errors.Errorf("error number of splitted coords: %v", len(splCoords))
	}
	var err error
	coords.I, err = strconv.Atoi(splCoords[0])
	if err != nil {
		return coords, errors.Wrap(err, "failed to parse i coordinate")
	}
	coords.J, err = strconv.Atoi(splCoords[1])
	if err != nil {
		return coords, errors.Wrap(err, "failed to parse j coordinate")
	}
	return coords, err
}

func DumpCoordinates(coords game.Coordinates) string {
	return fmt.Sprintf("%v,%v", coords.I, coords.J)
}

func GetBoardFromString(data string) (game.Board, error) {
	var board game.Board
	if len(data) != board.Size() {
		return board, errors.Errorf("len of dump (%v) not equal to size of board (%v)", len(data), board.Size())
	}
	for n, r := range data {
		mark := game.Mark(r - '0')
		i := n / len(board)
		j := n - i*len(board)
		board[i][j] = mark
	}
	return board, nil
}

func BoardToString(board game.Board) string {
	var sb strings.Builder
	sb.Grow(9)
	for i := range board {
		for j := range board[i] {
			sb.WriteString(strconv.Itoa(int(board[i][j])))
		}
	}

	return sb.String()
}

func GetMarkFromString(data string) (game.Mark, error) {
	markInt, err := strconv.Atoi(data)
	if err != nil {
		return game.MarkEmpty, errors.Wrap(err, "failed to parse int from data")
	}

	return game.Mark(markInt), nil
}

func MarkToString(mark game.Mark) string {
	return strconv.Itoa(int(mark))
}

func GetMarkRepresentation(mark game.Mark) string {
	switch mark {
	case game.MarkX:
		return "X"
	case game.MarkO:
		return "0"
	default:
		return " "
	}
}

func DumpCallback(player game.Mark, boardDump string, i, j int) string {
	return fmt.Sprintf("%v:%v:%v,%v", player, boardDump, i, j)
}

func ParseCallback(data string) (game.Mark, string, game.Coordinates, error) {
	var player game.Mark
	var boardDump string
	var coords game.Coordinates

	splitted := strings.Split(data, ":")
	if len(splitted) != 3 {
		return player, boardDump, coords, errors.Errorf("error number of splitted strigns: %v", len(splitted))
	}

	var err error
	player, err = GetMarkFromString(splitted[0])
	if err != nil {
		return player, boardDump, coords, errors.Wrap(err, "failed to parse player")
	}
	boardDump = splitted[1]
	strCoords := splitted[2]
	coords, err = GetCoordinatesFromString(strCoords)
	if err != nil {
		return player, boardDump, coords, errors.Wrap(err, "failed to get coords")
	}

	return player, boardDump, coords, nil
}
