package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const defaultLevel = zerolog.InfoLevel

type Logger struct {
	Inner zerolog.Logger
}

// NewLogger ロガーを生成する
//
// input:
//   - level (trace | debug | info | warn | error | faatl | panic)
//
// output:
//   - zerolog.Logger
//   - error
func NewLogger(level string) (*Logger, error) {
	parsedLevel := defaultLevel

	level = strings.TrimSpace(level)
	if level != "" {
		l, err := zerolog.ParseLevel(strings.ToLower(level))
		if err != nil {
			return nil, newError(fmt.Errorf("invalid log level: %w", err))
		}

		parsedLevel = l
	}

	return &Logger{Inner: zerolog.New(os.Stdout).Level(parsedLevel)}, nil
}
