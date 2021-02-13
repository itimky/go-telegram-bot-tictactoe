package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

func GetMoveFromString(data string) (game.Move, error) {
	var coords game.Move
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

func DumpMoveToString(i, j int) string {
	return fmt.Sprintf("%v,%v", i, j)
}

func GetMarkFromString(data string) (game.Mark, error) {
	markInt, err := strconv.Atoi(data)
	if err != nil {
		return game.MarkEmpty, errors.Wrap(err, "failed to parse int from data")
	}

	return game.Mark(markInt), nil
}

func DumpMarkToString(mark game.Mark) string {
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
