package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/itimky/go-telegram-bot-tictactoe/pkg/server"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Token is not provided")
	}

	srv, err := server.NewServer(token)
	if err != nil {
		log.Fatal(err)
	}
	srv.Run()
}
