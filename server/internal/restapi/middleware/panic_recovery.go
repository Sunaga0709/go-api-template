package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"

	"github.com/Sunaga0709/go-api-template/internal/logger"
)

type PanicRecoveryMiddleware struct {
	logger logger.Logger
}

func NewPanicRecoveryMiddleware(logger logger.Logger) *PanicRecoveryMiddleware {
	return &PanicRecoveryMiddleware{logger: logger}
}

func (p *PanicRecoveryMiddleware) Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				detail := NewErrorDetail(
					http.StatusInternalServerError,
					fmt.Errorf("panic recovered: %v", recovered),
				)
				p.logger.Inner.Error().
					Err(detail).
					Str("event", "panic_recovered").
					Str("method", ctx.Request().Method).
					Str("path", ctx.Request().URL.Path).
					Int("status", http.StatusInternalServerError).
					Str("panic", fmt.Sprint(recovered)).
					Str("stack", string(debug.Stack())).
					Send()
				err = detail
			}
		}()

		return next(ctx)
	}
}
