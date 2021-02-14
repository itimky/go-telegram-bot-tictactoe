package storage

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

const gameRedisStorageKeyPrefix = "ttt_game"
const gameRedisStorageTTL = 2 * (24 * time.Hour) // 2 Days

func makeKey(key string) string {
	return gameRedisStorageKeyPrefix + "_" + key
}

type gameContainer struct {
	Board  game.Board `json:"board"`
	Player game.Mark  `json:"player"`
}

type aiGameContainer struct {
	Difficulty ai.Difficulty `json:"difficulty"`
	Game       gameContainer `json:"game"`
}

type GameRedisStorage struct {
	client *redis.Client
}

func (gs *GameRedisStorage) Load(msgID int) (ai.AIGame, error) {
	log.Debug("Loading game ", msgID)
	key := makeKey(strconv.Itoa(msgID))
	val, err := gs.client.Get(context.Background(), key).Result()
	g := ai.AIGame{}
	switch err {
	case nil:
	case redis.Nil:
		return g, errors.New("game not found")
	default:
		return g, errors.Wrap(err, "failed to hget game from redis")
	}

	g, err = unmarshalGameFromRedis([]byte(val))
	if err != nil {
		return g, errors.Wrap(err, "failed to unmarshal game from redis")
	}

	return g, nil

}

func (gs *GameRedisStorage) Save(msgID int, g ai.AIGame) error {
	log.Debug("Saving game ", msgID)
	key := makeKey(strconv.Itoa(msgID))
	val, err := marshalGameToRedis(g)
	if err != nil {
		return errors.Wrap(err, "failed to marshal game")
	}
	err = gs.client.Set(context.Background(), key, val, gameRedisStorageTTL).Err()
	if err != nil {
		return errors.Wrap(err, "failed to hset game to redis")
	}
	return nil
}

func marshalGameToRedis(g ai.AIGame) ([]byte, error) {
	container := aiGameContainer{
		Difficulty: g.Difficulty,
		Game: gameContainer{
			Board:  g.Game.GetBoard(),
			Player: g.Game.GetPlayer(),
		},
	}
	val, err := json.Marshal(container)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal game")
	}
	log.Debug("Marshaled game container: ", string(val))
	return val, nil
}

func unmarshalGameFromRedis(data []byte) (ai.AIGame, error) {
	container := aiGameContainer{}
	aiGame := ai.AIGame{}

	if err := json.Unmarshal(data, &container); err != nil {
		return aiGame, errors.Wrap(err, "failed to unmarshal redis data")
	}
	log.Debug("Unmarshaled game container: ", container)

	aiGame.Game = game.ContinueGame(container.Game.Player, container.Game.Board)
	aiGame.Difficulty = container.Difficulty
	return aiGame, nil
}

func NewGameRedisStorage(client *redis.Client) *GameRedisStorage {
	return &GameRedisStorage{
		client: client,
	}
}
