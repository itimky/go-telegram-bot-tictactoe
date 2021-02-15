package session

import (
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

type GameType byte

const (
	GameTypeVersusAI = iota
	GameTypeVersusHuman
)

type session struct {
	id         int
	game       game.Game
	gameType   GameType
	difficulty ai.Difficulty
}

type Service struct {
	ai      ai.AI
	storage ISessionStorage
}

func (s *Service) New(sessionID int, player game.Mark, dif ai.Difficulty) (game.Game, error) {
	session := &session{
		id:         sessionID,
		game:       game.NewGame(player),
		gameType:   GameTypeVersusAI,
		difficulty: dif,
	}
	if session.gameType == GameTypeVersusAI {
		if session.game.GetPlayer() == game.MarkO {
			g, err := s.ai.MakeAIMove(session.difficulty, session.game)
			if err != nil {
				return g, errors.Wrap(err, "failed to make ai move")
			}
			session.game = g
			session.game.SwapPlayers()
		}
	}

	if err := s.storage.Save(session); err != nil {
		return session.game, errors.Wrap(err, "failed to save session")
	}

	return session.game, nil
}

func (s *Service) Play(sessionID int, move game.Move) (game.Game, error) {
	session, err := s.storage.Load(sessionID)
	if err != nil {
		return game.Game{}, errors.Wrap(err, "failed to load session")
	}

	if err = session.game.MakeMove(move); err != nil {
		return session.game, errors.Wrap(err, "failed to make move")
	}

	if session.gameType == GameTypeVersusAI {
		session.game.SwapPlayers()
		g, err := s.ai.MakeAIMove(session.difficulty, session.game)
		if err != nil {
			return g, errors.Wrap(err, "failed to make ai move")
		}
		session.game = g
		session.game.SwapPlayers()
	}

	if err = s.storage.Save(session); err != nil {
		return session.game, errors.Wrap(err, "failed to save session")
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
	return &Service{
		ai:      ai.NewAI(),
		storage: s,
	}
}
