package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port             string `envconfig:"SERVER_PORT" required:"true"`
	Environment      string `envconfig:"ENVIRONMENT" required:"true"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"info"`
	HiddenSuccessLog bool   `envconfig:"HIDDEN_SUCCESS_LOG" default:"true"`
}

func loadConfig() (config, error) {
	var conf config
	if err := envconfig.Process("", &conf); err != nil {
		return conf, newError(fmt.Errorf("failed to load config: %w", err))
	}

	return conf, nil
}
