package main

import (
	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	TgBotToken string `envconfig:"TG_BOT_TOKEN" desc:"Telegram bot token"`
	RedisURL   string `envconfig:"REDIS_URL" default:"" desc:"Redis url. If not provided, in-mem cache is used"`

	LogLevel  string `envconfig:"LOG_LEVEL" default:"warn" desc:"Level of logging verbosity. Could be debug, info, warn, error"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"text" desc:"Logging format. Could be text, json"`
}

func getConfig() (config, error) {
	cfg := config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to parse config from env")
	}
	return cfg, nil
}
