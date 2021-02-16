package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

const gameRedisStorageKeyPrefix = "ttt_game"
const gameRedisStorageTTL = 2 * (24 * time.Hour) // 2 Days

type Row []byte

func (r Row) MarshalJSON() ([]byte, error) {
	var result string
	if r == nil {
		return nil, errors.Errorf("board row cannot be nil")
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", r)), ",")
	}
	return []byte(result), nil
}

func makeKey(key string) string {
	return gameRedisStorageKeyPrefix + "_" + key
}

type gameContainer struct {
	Board   []Row     `json:"board"`
	N       int       `json:"n"`
	LineLen int       `json:"line_len"`
	Player  game.Mark `json:"player"`
}

type SessionContainer struct {
	ID         int           `json:"id"`
	Difficulty ai.Difficulty `json:"difficulty"`
	Game       gameContainer `json:"game"`
	Type       GameType      `json:"type"`
}

type SessionRedisStorage struct {
	client *redis.Client
}

func (gs *SessionRedisStorage) Load(msgID int) (*session, error) {
	log.Debug("Loading game ", msgID)
	key := makeKey(strconv.Itoa(msgID))
	val, err := gs.client.Get(context.Background(), key).Result()
	switch err {
	case nil:
	case redis.Nil:
		return nil, errors.New("game not found")
	default:
		return nil, errors.Wrap(err, "failed to hget game from redis")
	}

	session, err := unmarshalGameFromRedis([]byte(val))
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal game from redis")
	}

	return session, nil

}

func (gs *SessionRedisStorage) Save(session *session) error {
	log.Debug("Saving session ", session.id)
	key := makeKey(strconv.Itoa(session.id))
	val, err := marshalGameToRedis(session)
	if err != nil {
		return errors.Wrap(err, "failed to marshal game")
	}
	err = gs.client.Set(context.Background(), key, val, gameRedisStorageTTL).Err()
	if err != nil {
		return errors.Wrap(err, "failed to hset game to redis")
	}
	return nil
}

func marshalGameToRedis(session *session) ([]byte, error) {
	n := session.game.Size()
	rows := make([]Row, 0, n)
	for _, gameRow := range session.game.GetBoard() {
		row := make(Row, 0, n)
		for _, mark := range gameRow {
			row = append(row, byte(mark))
		}
		rows = append(rows, row)
	}

	container := SessionContainer{
		ID:         session.id,
		Difficulty: session.difficulty,
		Type:       session.gameType,
		Game: gameContainer{
			Board:   rows,
			N:       n,
			LineLen: session.game.LineLen(),
			Player:  session.game.GetPlayer(),
		},
	}
	val, err := json.Marshal(container)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal game")
	}
	log.Debug("Marshaled game container: ", string(val))
	return val, nil
}

func unmarshalGameFromRedis(data []byte) (*session, error) {
	container := SessionContainer{}
	if err := json.Unmarshal(data, &container); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal redis data")
	}
	log.Debug("Unmarshaled game container: ", container)

	board := make(game.Board, 0, container.Game.N)
	for _, row := range container.Game.Board {
		markRow := make([]game.Mark, 0, container.Game.N)
		for _, m := range row {
			markRow = append(markRow, game.Mark(m))
		}
		board = append(board, markRow)
	}
	session := session{
		id:         container.ID,
		game:       game.ContinueGame(container.Game.Player, board, container.Game.N, container.Game.LineLen),
		difficulty: container.Difficulty,
		gameType:   container.Type,
	}
	return &session, nil
}

func NewSessionRedisStorage(client *redis.Client) *SessionRedisStorage {
	return &SessionRedisStorage{
		client: client,
	}
}
