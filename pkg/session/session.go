package session

import (
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

type GameType string

const (
	GameTypeVersusAI    = "vs_ai"
	GameTypeVersusHuman = "vs_human"
)

type session struct {
	id         int
	game       game.Game
	gameType   GameType
	difficulty ai.Difficulty
}

type Service struct {
	ai                 ai.AI
	gameSizeLineLenMap map[int]int
	storage            ISessionStorage
}

func (s *Service) New(sessionID int, player game.Mark, dif ai.Difficulty, size int) (game.Game, error) {
	session := &session{
		id:         sessionID,
		game:       game.NewGame(player, size, s.gameSizeLineLenMap[size]),
		gameType:   GameTypeVersusAI,
		difficulty: dif,
	}
	if session.gameType == GameTypeVersusAI {
		if session.game.GetPlayer() == game.MarkO {
			g, err := s.ai.MakeAIMove(session.difficulty, session.game)
			if err != nil {
				return g, err
			}
			session.game = g
			session.game.SwapPlayers()
		}
	}

	if err := s.storage.Save(session); err != nil {
		return session.game, err
	}

	return session.game, nil
}

func (s *Service) Play(sessionID int, move game.Move) (game.Game, error) {
	session, err := s.storage.Load(sessionID)
	if err != nil {
		return game.Game{}, err
	}

	if err = session.game.MakeMove(move); err != nil {
		return session.game, err
	}

	if session.gameType == GameTypeVersusAI {
		session.game.SwapPlayers()
		g, err := s.ai.MakeAIMove(session.difficulty, session.game)
		if err != nil {
			return g, err
		}
		session.game = g
		session.game.SwapPlayers()
	}

	if err = s.storage.Save(session); err != nil {
		return session.game, err
	}

	return session.game, nil
}

func NewService(redisClient *redis.Client) *Service {
	var s ISessionStorage
	if redisClient != nil {
		s = NewSessionRedisStorage(redisClient)
		log.Info("using redis session storage")
	} else {
		s = NewInMemStorage()
		log.Info("using in-mem session storage")
	}

	gameSizeLineLenMap := map[int]int{
		3:  3,
		4:  4,
		5:  4,
		6:  4,
		7:  5,
		8:  5,
		9:  5,
		10: 6,
	}

	return &Service{
		ai:                 ai.NewAI(),
		gameSizeLineLenMap: gameSizeLineLenMap,
		storage:            s,
	}
}
