package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type secret struct {
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func loadSecret() (secret, error) {
	var secr secret
	if err := envconfig.Process("", &secr); err != nil {
		return secr, newError(fmt.Errorf("failed to load secret: %w", err))
	}

	return secr, nil
}
