package middleware

import (
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

type ShutdownMiddleware struct {
	shuttingDown atomic.Bool
}

func NewShutdownMiddleware() *ShutdownMiddleware {
	return &ShutdownMiddleware{}
}

func (s *ShutdownMiddleware) BeginShutdown() {
	s.shuttingDown.Store(true)
}

func (s *ShutdownMiddleware) RejectDuringShutdown(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if s.shuttingDown.Load() {
			return NewErrorDetail(
				http.StatusServiceUnavailable,
				errors.New("server is shutting down"),
			)
		}

		return next(ctx)
	}
}
