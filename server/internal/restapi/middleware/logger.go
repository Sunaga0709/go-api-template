package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/Sunaga0709/go-api-template/internal/logger"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

type LoggerMiddleware struct {
	logger        logger.Logger
	hiddenSuccess bool
}

func NewLoggerMiddleware(logger logger.Logger, hiddenSuccess bool) *LoggerMiddleware {
	return &LoggerMiddleware{logger: logger, hiddenSuccess: hiddenSuccess}
}

// AccessLogging アクセスログを出力するミドルウェア
func (l *LoggerMiddleware) AccessLogging(f gen.StrictHandlerFunc, operationID string) gen.StrictHandlerFunc {
	return func(ctx echo.Context, request any) (any, error) {
		receivedAt := time.Now()
		res, err := f(ctx, request)
		elapsed := time.Since(receivedAt)

		status := http.StatusOK
		if err != nil {
			status = http.StatusInternalServerError
			if detail, ok := errors.AsType[ErrorDetail](err); ok {
				status = detail.Status
			}
		}

		if l.hiddenSuccess && status == http.StatusOK {
			return res, err
		}

		event := l.logger.Inner.WithLevel(logLevel(status)).
			Str("operation_id", operationID).
			Str("method", ctx.Request().Method).
			Str("path", ctx.Request().URL.Path).
			Int("status", status).
			Str("timestamp", receivedAt.UTC().Format(time.RFC3339)).
			Str("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds()))

		if err != nil {
			event = event.Err(err)
		}

		event.Send()

		return res, err
	}
}

// logLevel ステータスコードに応じたログレベルを返す。
func logLevel(status int) zerolog.Level {
	switch {
	case status >= 500:
		return zerolog.ErrorLevel
	case status >= 400:
		return zerolog.WarnLevel
	default:
		return zerolog.InfoLevel
	}
}
