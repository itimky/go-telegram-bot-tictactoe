package server

import (
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/server/storage"
)

type Server struct {
	bot         *tb.Bot
	gameStorage storage.IGameStorage
}

func NewServer(token string, redisClient *redis.Client) (*Server, error) {
	var gameStorage storage.IGameStorage
	if redisClient != nil {
		gameStorage = storage.NewGameRedisStorage(redisClient)
		log.Info("using redis game storage")
	} else {
		gameStorage = storage.NewGameAppStorage()
		log.Info("using app game storage")
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create bot instance")
	}

	server := Server{
		bot:         b,
		gameStorage: gameStorage,
	}
	b.Handle("/start", server.onStart)
	b.Handle(tb.OnCallback, server.onCallback)
	return &server, nil
}

func (s *Server) Run() {
	s.bot.Start()
}

func (s *Server) onStart(m *tb.Message) {
	start := time.Now()
	log.Info("Handling /start being")
	_, err := s.bot.Send(m.Sender, "Choose player", getChoosePlayerMarkup())
	if err != nil {
		log.Error("failed to send response: ", err)
	}
	log.WithField("elapsed", time.Since(start).Seconds()).Info("Handling /start end")
}

func (s *Server) onCallback(q *tb.Callback) {
	start := time.Now()
	log.Info("Handling callback begin")
	log.Info("Callback: ", *q)
	replyMarkup, err := s.handleCallback(q.Message.ID, q.Data)
	log.WithField("elapsed", time.Since(start).Seconds()).Info("callback calculated")
	if err != nil {
		log.Error("failed to handle callback: ", err)
	}
	if _, err := s.bot.EditReplyMarkup(q.Message, replyMarkup); err != nil {
		log.Error("failed to send response: ", err)
	}
	log.WithField("elapsed", time.Since(start).Seconds()).Info("Handling callback end")
}

func (s *Server) handleCallback(msgID int, data string) (*tb.ReplyMarkup, error) {
	if len(data) == 1 {
		markup, err := s.startGame(msgID, data)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start game")
		}
		return markup, err
	}
	markup, err := s.playRound(msgID, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to play round")
	}
	return markup, nil
}

func (s *Server) startGame(msgID int, data string) (*tb.ReplyMarkup, error) {
	playerMark, err := GetMarkFromString(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mark from string")
	}
	aiGame, err := ai.StartAIGame(ai.Hard, playerMark)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start AI game")
	}
	if err = s.gameStorage.Save(msgID, aiGame); err != nil {
		return nil, errors.Wrap(err, "failed to save game")
	}

	return getGameMarkup(aiGame.Game), nil
}

func (s *Server) playRound(msgID int, data string) (*tb.ReplyMarkup, error) {
	aiGame, err := s.gameStorage.Load(msgID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load game")
	}
	move, err := GetMoveFromString(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse move")
	}
	if err := aiGame.Game.MakeMove(move); err != nil {
		return nil, errors.Wrap(err, "failed to make move")
	}
	if err = aiGame.MakeAIMove(); err != nil {
		return nil, errors.Wrap(err, "failed to play AI opponent")
	}
	if err = s.gameStorage.Save(msgID, aiGame); err != nil {
		return nil, errors.Wrap(err, "failed to save game")
	}
	return getGameMarkup(aiGame.Game), nil
}

func getChoosePlayerMarkup() *tb.ReplyMarkup {
	buttons := []tb.InlineButton{
		{
			Text: "Play as X",
			Data: DumpMarkToString(game.MarkX),
		},
		{
			Text: "Play as O",
			Data: DumpMarkToString(game.MarkO),
		},
	}

	return &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{buttons},
	}
}

func getGameMarkup(g game.Game) *tb.ReplyMarkup {
	board := g.GetBoard()
	buttons := make([][]tb.InlineButton, 0, len(board))
	for i, line := range board {
		rowButtons := make([]tb.InlineButton, 0, len(line))
		for j, mark := range line {
			button := tb.InlineButton{
				Text: GetMarkRepresentation(mark),
				Data: DumpMoveToString(i, j),
			}
			rowButtons = append(rowButtons, button)
		}
		buttons = append(buttons, rowButtons)
	}

	return &tb.ReplyMarkup{
		InlineKeyboard: buttons,
	}
}
