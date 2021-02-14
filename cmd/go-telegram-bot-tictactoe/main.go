package main

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/server"
)

func initLogger(cfg config) {
	switch cfg.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		fallthrough
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		fallthrough
	default:
		log.SetLevel(log.WarnLevel)
	}
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	initLogger(cfg)

	redisClient, err := getRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewServer(cfg.TgBotToken, redisClient)
	if err != nil {
		log.Fatal(err)
	}
	srv.Run()
}

func getRedisClient(url string) (*redis.Client, error) {
	if url == "" {
		return nil, nil
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "invalid redis url")
	}

	return redis.NewClient(opts), nil
}
