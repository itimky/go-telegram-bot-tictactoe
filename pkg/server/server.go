package server

import (
	"fmt"
	"github.com/itimky/go-telegram-bot-tictactoe/pkg/ai"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/game"
)

type Server struct {
	bot *tb.Bot
}

func NewServer(token string) (*Server, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create bot instance: %w", err)
	}
	server := Server{bot: b}
	b.Handle("/start", server.onStart)
	b.Handle(tb.OnCallback, server.onCallback)
	return &server, nil
}

func (s *Server) Run() {
	s.bot.Start()
}

func (s *Server) onStart(m *tb.Message) {
	start := time.Now()
	log.Warn(s)
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
	replyMarkup, err := handleCallback(q.Data)
	log.WithField("elapsed", time.Since(start).Seconds()).Info("callback calculated")
	if err != nil {
		log.Error("failed to handle callback: ", err)
	}
	if _, err := s.bot.EditReplyMarkup(q.Message, replyMarkup); err != nil {
		log.Error("failed to send response: ", err)
	}
	log.WithField("elapsed", time.Since(start).Seconds()).Info("Handling callback end")
}

func getChoosePlayerMarkup() *tb.ReplyMarkup {
	buttons := []tb.InlineButton{
		{
			Text: "Play as X",
			Data: MarkToString(game.MarkX),
		},
		{
			Text: "Play as O",
			Data: MarkToString(game.MarkO),
		},
	}

	return &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{buttons},
	}
}

func handleCallback(data string) (*tb.ReplyMarkup, error) {
	if len(data) == 1 {
		markup, err := startGame(data)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start game")
		}
		return markup, err
	}
	markup, err := playRound(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to play round")
	}
	return markup, nil
}

func startGame(data string) (*tb.ReplyMarkup, error) {
	playerMark, err := GetMarkFromString(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mark from string")
	}
	g := game.NewGame(playerMark)
	if !g.IsPlayerFirst() {
		aiPlayer := ai.NewAI(ai.Hard)
		if err = aiPlayer.PlayOpponent(g); err != nil {
			return nil, errors.Wrap(err, "failed to play AI opponent")
		}
	}
	return getGameMarkup(g), nil
}

func playRound(data string) (*tb.ReplyMarkup, error) {
	player, boardDump, coords, err := ParseCallback(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse callback")
	}
	board, err := GetBoardFromString(boardDump)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get board")
	}
	g := game.ContinueGame(player, board)
	if err := g.MakeMove(coords); err != nil {
		return nil, errors.Wrap(err, "failed to make move")
	}
	aiPlayer := ai.NewAI(ai.Hard)
	if err = aiPlayer.PlayOpponent(g); err != nil {
		return nil, errors.Wrap(err, "failed to play AI opponent")
	}
	return getGameMarkup(g), nil
}

func getGameMarkup(g *game.Game) *tb.ReplyMarkup {
	board := g.GetBoard()
	player := g.GetPlayer()
	boardDump := BoardToString(board)
	buttons := make([][]tb.InlineButton, 0, len(board))
	for i, line := range board {
		rowButtons := make([]tb.InlineButton, 0, len(line))
		for j, mark := range line {
			button := tb.InlineButton{
				Text: GetMarkRepresentation(mark),
				Data: DumpCallback(player, boardDump, i, j),
			}
			rowButtons = append(rowButtons, button)
		}
		buttons = append(buttons, rowButtons)
	}

	return &tb.ReplyMarkup{
		InlineKeyboard: buttons,
	}
}
